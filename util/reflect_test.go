package util

import (
	//"fmt"
	"reflect"
	"testing"

	ut "github.com/ben-han-cn/cement/unittest"
)

func TestValueKind(t *testing.T) {
	var v interface{}
	v = uint8(10)
	ut.Equal(t, Uint, Inspect(reflect.TypeOf(v)))

	v = int32(10)
	ut.Equal(t, Int, Inspect(reflect.TypeOf(v)))

	type MyFlag string
	v = MyFlag("x")
	ut.Equal(t, String, Inspect(reflect.TypeOf(v)))

	type MyStruct struct {
	}
	v = MyStruct{}
	ut.Equal(t, Struct, Inspect(reflect.TypeOf(v)))
	v = &MyStruct{}
	ut.Equal(t, StructPtr, Inspect(reflect.TypeOf(v)))

	v = []int8{}
	ut.Equal(t, IntSlice, Inspect(reflect.TypeOf(v)))
	v = []string{}
	ut.Equal(t, StringSlice, Inspect(reflect.TypeOf(v)))
	v = []MyStruct{}
	ut.Equal(t, StructSlice, Inspect(reflect.TypeOf(v)))
	v = []*MyStruct{}
	ut.Equal(t, StructPtrSlice, Inspect(reflect.TypeOf(v)))
	v = []MyFlag{}
	ut.Equal(t, StringSlice, Inspect(reflect.TypeOf(v)))

	v = map[string]string{}
	ut.Equal(t, StringStringMap, Inspect(reflect.TypeOf(v)))
	v = map[string]MyStruct{}
	ut.Equal(t, StringStructMap, Inspect(reflect.TypeOf(v)))
	v = map[string]*MyStruct{}
	ut.Equal(t, StringStructPtrMap, Inspect(reflect.TypeOf(v)))
	v = map[string]MyFlag{}
	ut.Equal(t, StringStringMap, Inspect(reflect.TypeOf(v)))
}
