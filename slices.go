// Package slices implements simple functions to manipulate slices.
package slices

import (
	"math/rand"
	"reflect"
	"sort"
)

// Operator of slice.
type Op struct {
	Slice interface{}
}

// Creates operator of "[]T".
// Type of s must be "T[]".
func NewOp(s interface{}) *Op {
	return &Op{s}
}

// If op is operator of "[]T", type of f must be "func(T) T2".
// Returns Operator of "[]T2".
func (op *Op) Map(f interface{}) *Op {
	vs := reflect.ValueOf(op.Slice)
	vf := reflect.ValueOf(f)

	tf := vf.Type()
	if tf.NumIn() != 1 {
		panic("Number of Argument must be 1")
	}
	if tf.NumOut() != 1 {
		panic("Number of return value must be 1")
	}

	tif := tf.In(0)
	tof := tf.Out(0)
	if tif != vs.Type().Elem() {
		panic("Mismatch function type")
	}

	len := vs.Len()
	vos := reflect.MakeSlice(reflect.SliceOf(tof), 0, len)
	for i := 0; i < len; i++ {
		vos = reflect.Append(vos, vf.Call([]reflect.Value{vs.Index(i)})[0])
	}

	return &Op{vos.Interface()}
}

// If op is operator of "[]T" and type of o is "T2", type of f must be "func(T2, T) T2".
// Returns value of "T2"
func (op *Op) Inject(o, f interface{}) interface{} {
	vs := reflect.ValueOf(op.Slice)
	vo := reflect.ValueOf(o)
	vf := reflect.ValueOf(f)

	tf := vf.Type()
	if tf.NumIn() != 2 {
		panic("Number of Argument must be 2")
	}
	if tf.NumOut() != 1 {
		panic("Number of return value must be 1")
	}

	tfi0 := tf.In(0)
	tfi1 := tf.In(1)
	tfo := tf.Out(0)
	if vo.Type() != tfi0 || tfo != tfi0 || tfi1 != vs.Type().Elem() {
		panic("Mismatch function type")
	}

	len := vs.Len()
	for i := 0; i < len; i++ {
		vo = vf.Call([]reflect.Value{vo, vs.Index(i)})[0]
	}

	return vo.Interface()
}

func validateBooleanFunc(vs, vf *reflect.Value) {
	tf := vf.Type()
	if tf.NumIn() != 1 {
		panic("Number of Argument must be 1")
	}
	if tf.NumOut() != 1 {
		panic("Number of return value must be 1")
	}

	if tf.In(0) != vs.Type().Elem() {
		panic("Function argument Type is invalid")
	}
	if tf.Out(0) != reflect.TypeOf(true) {
		panic("Function return Type is invalid")
	}

}

// If op is operator of "[]T", type of f must be "func(T) bool".
// Returns Operator of []T
func (op *Op) Select(f interface{}) *Op {
	vs := reflect.ValueOf(op.Slice)
	vf := reflect.ValueOf(f)

	validateBooleanFunc(&vs, &vf)

	len := vs.Len()
	vos := reflect.MakeSlice(vs.Type(), 0, len)

	for i := 0; i < len; i++ {
		v := vs.Index(i)
		if vf.Call([]reflect.Value{v})[0].Bool() {
			vos = reflect.Append(vos, v)
		}
	}

	return &Op{vos.Interface()}
}

// If op is operator of "[]T", type of f must be "func(T) bool".
func (op *Op) All(f interface{}) bool {
	vs := reflect.ValueOf(op.Slice)
	vf := reflect.ValueOf(f)

	validateBooleanFunc(&vs, &vf)

	len := vs.Len()
	for i := 0; i < len; i++ {
		if !vf.Call([]reflect.Value{vs.Index(i)})[0].Bool() {
			return false
		}
	}

	return true
}

// If op is operator of "[]T", type of f must be "func(T) bool".
func (op *Op) Any(f interface{}) bool {
	vs := reflect.ValueOf(op.Slice)
	vf := reflect.ValueOf(f)

	validateBooleanFunc(&vs, &vf)

	len := vs.Len()
	for i := 0; i < len; i++ {
		if vf.Call([]reflect.Value{vs.Index(i)})[0].Bool() {
			return true
		}
	}

	return false
}

// If op is Operator of "[]T", type of f must be "func(T) T2".
// And GroupBy(f) returns value of "map[T2] []T"
func (op *Op) GroupBy(f interface{}) interface{} {
	vs := reflect.ValueOf(op.Slice)
	vf := reflect.ValueOf(f)

	tf := vf.Type()
	if tf.NumIn() != 1 {
		panic("Number of Argument must be 1")
	}
	if tf.NumOut() != 1 {
		panic("Number of return value must be 1")
	}

	tif := tf.In(0)
	tof := tf.Out(0)
	if tif != vs.Type().Elem() {
		panic("Mismatch function type")
	}

	len := vs.Len()
	vom := reflect.MakeMap(reflect.MapOf(tof, vs.Type()))
	for i := 0; i < len; i++ {
		v := vs.Index(i)
		vk := vf.Call([]reflect.Value{v})[0]
		vi := vom.MapIndex(vk)
		if vi.IsValid() {
			vom.SetMapIndex(vk, reflect.Append(vi, v))
		} else {
			vom.SetMapIndex(vk, reflect.Append(reflect.MakeSlice(vs.Type(), 0, len), v))
		}
	}

	return vom.Interface()
}

// Returns operator of copied slice.
func (op *Op) Copy() *Op {
	vs := reflect.ValueOf(op.Slice)
	len := vs.Len()
	vos := reflect.MakeSlice(vs.Type(), 0, len)

	for i := 0; i < len; i++ {
		vos = reflect.Append(vos, vs.Index(i))
	}

	return &Op{vos.Interface()}
}

type wrappedOp struct {
	*Op
	less interface{}
}

func (op *Op) len() int {
	return reflect.ValueOf(op.Slice).Len()
}

func (wop *wrappedOp) Len() int {
	return wop.Op.len()
}

func (op *Op) swap(i, j int) {
	vs := reflect.ValueOf(op.Slice)
	tmp := reflect.Indirect(reflect.New(vs.Type().Elem()))
	tmp.Set(vs.Index(i))
	vs.Index(i).Set(vs.Index(j))
	vs.Index(j).Set(tmp)
}

func (wop *wrappedOp) Swap(i, j int) {
	wop.Op.swap(i, j)
}

func (wop *wrappedOp) Less(i, j int) bool {
	vs := reflect.ValueOf(wop.Op.Slice)
	vf := reflect.ValueOf(wop.less)
	return vf.Call([]reflect.Value{vs.Index(i), vs.Index(j)})[0].Bool()
}

func validateSortFunc(vs, vf *reflect.Value) {
	tf := vf.Type()
	if tf.NumIn() != 2 {
		panic("Number of Argument must be 2")
	}
	if tf.NumOut() != 1 {
		panic("Number of return value must be 1")
	}

	tfi0 := tf.In(0)
	tfi1 := tf.In(1)
	tfo := tf.Out(0)
	if tfi0 != tfi1 || tfi0 != vs.Type().Elem() || tfo != reflect.TypeOf(true) {
		panic("Mismatch function type")
	}
}

// If op is Operator of "[]T", type of f must be "func(T, T) bool".
// Sorts self in place, and returns self operator.
//
// This method performs unstable sort.
func (op *Op) Sort(f interface{}) *Op {
	vs := reflect.ValueOf(op.Slice)
	vf := reflect.ValueOf(f)

	validateSortFunc(&vs, &vf)

	w := wrappedOp{op, f}
	sort.Sort(&w)
	return op
}

// If op is Operator of "[]T", type of f must be "func(T, T) bool".
// Sorts self in place, and returns self operator.
//
// This method performs stable sort.
func (op *Op) StableSort(f interface{}) *Op {
	vs := reflect.ValueOf(op.Slice)
	vf := reflect.ValueOf(f)

	validateSortFunc(&vs, &vf)

	w := wrappedOp{op, f}
	sort.Stable(&w)
	return op
}

// Shuffles self in place, and returns self operator.
func (op *Op) Shuffle(r *rand.Rand) *Op {
	len := op.len()
	for i := len - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		if i != j {
			op.swap(i, j)
		}
	}

	return op
}
