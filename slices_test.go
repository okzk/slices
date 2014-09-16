package slices

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func twice(i int) int {
	return i * 2
}

func TestMapInt(t *testing.T) {
	src := []int{2, 4}
	expected := []int{8, 16}
	actual := NewOp(src).Map(twice).Map(twice).Slice.([]int)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Map faild, expected: %v, actual: %v", expected, actual)
	}
}

func TestMapString(t *testing.T) {
	src := []int{100, 400, 1234}
	actual := NewOp(src).Map(func(i int) string {
		return fmt.Sprintf("%d", i)
	}).Slice.([]string)
	expected := []string{"100", "400", "1234"}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Map faild, expected: %v, actual: %v", expected, actual)
	}
}

func isEven(i int) bool {
	return i%2 == 0
}

func TestSelect(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := []int{2, 4, 6, 8, 10}
	actual := NewOp(src).Select(isEven).Slice.([]int)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Select faild, expected: %v, actual: %v", expected, actual)
	}
}

func TestAny(t *testing.T) {
	srcTrue := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if !NewOp(srcTrue).Any(isEven) {
		t.Errorf("Any of %v is not even?", srcTrue)
	}

	srcFalse := []int{1, 3, 5, 7, 9}
	if NewOp(srcFalse).Any(isEven) {
		t.Error("Any of %v is even?", srcFalse)
	}

}

func TestAll(t *testing.T) {
	srcTrue := []int{2, 4, 6, 8, 10}
	if !NewOp(srcTrue).All(isEven) {
		t.Error("All of %v is not even?", srcTrue)
	}

	srcFalse := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if NewOp(srcFalse).All(isEven) {
		t.Error("All of %v is even?", srcFalse)
	}
}

func sum(a, b int) int {
	return a + b
}

func TestInject(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	actual := NewOp(src).Inject(0, sum).(int)
	expected := 55

	if actual != expected {
		t.Errorf("Inject faild, expected: %d, actual: %d\n", expected, actual)
	}
}

func TestGroupBy(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	actual := NewOp(src).GroupBy(isEven).(map[bool][]int)
	expected := map[bool][]int{true: []int{2, 4, 6, 8, 10}, false: []int{1, 3, 5, 7, 9}}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("GroupBy faild, expected: %v, actual: %v\n", expected, actual)
	}
}

func TestCopy(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	dst := NewOp(src).Copy().Slice.([]int)

	if !reflect.DeepEqual(src, dst) {
		t.Errorf("Copy faild, expected: %v, actual: %v\n", src, dst)
	}

	src[0] = 1000
	if src[0] == dst[0] {
		t.Errorf("Copy faild, expected: %d, actual: %d\n", src[0], dst[0])
	}
}

func TestLen(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	actual := NewOp(src).len()
	expected := 10

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Len faild, expected: %v, actual: %v\n", expected, actual)
	}
}

func TestSwap(t *testing.T) {
	src := []int{1, 2, 3}
	NewOp(src).swap(0, 1)
	expected := []int{2, 1, 3}

	if !reflect.DeepEqual(src, expected) {
		t.Errorf("Copy faild, expected: %v, actual: %v\n", expected, src)
	}
}

func TestSort(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 10, 9, 8, 7, 6}
	actual := NewOp(src).Copy().Sort(func(i, j int) bool {
		return i < j
	}).Slice.([]int)
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Sort faild, expected: %v, actual: %v\n", expected, actual)
	}
}

func TestShuffle(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	r := rand.New(rand.NewSource(0))
	actual := NewOp(src).Copy().Shuffle(r).Slice.([]int)
	expected := []int{7, 9, 3, 4, 8, 6, 10, 2, 1, 5}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Sort faild, expected: %v, actual: %v\n", expected, actual)
	}
}

func primitiveInject(t []int, o int, f func(int, int) int) int {
	for _, v := range t {
		o = f(o, v)
	}
	return o
}

func BenchmarkInject(b *testing.B) {
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	op := NewOp(src)
	for i := 0; i < b.N; i++ {
		op.Inject(0, sum)
	}
}

func BenchmarkPrimitiveInject(b *testing.B) {
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < b.N; i++ {
		primitiveInject(src, 0, sum)
	}
}
