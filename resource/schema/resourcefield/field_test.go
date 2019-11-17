package resourcefield

import (
	"encoding/json"
	ut "github.com/ben-han-cn/cement/unittest"
	"reflect"
	"testing"
)

func TestFieldBuild(t *testing.T) {
	builder := NewBuilder()
	sf, err := builder.Build(reflect.TypeOf(TestStruct{}))
	ut.Assert(t, err == nil, "")

	ut.Equal(t, len(sf.fields), 14)
	fieldNames := []string{
		"Id",
		"Age",
		"Name",
		"StringWithOption",
		"StringWithLenLimit",
		"IntWithRange",
		"StringIntMap",
		"IntSlice",
		"StringSliceWithOption",

		"SliceComposition",
		"StringMapCompostion",

		"PtrComposition",
		"SlicePtrComposition",
		"StringPtrMapCompostion",
	}
	for _, name := range fieldNames {
		_, ok := sf.fields[name]
		ut.Assert(t, ok, "")
	}
}

func TestInvalidField(t *testing.T) {
	type S1 struct {
		StringWithLenLimit string `json:"stringWithLenLimit" rest:"minLen=20,maxLen=10"`
	}
	builder := NewBuilder()
	_, err := builder.Build(reflect.TypeOf(S1{}))
	ut.Assert(t, err != nil, "")

	type S2 struct {
		IntWithRange uint32 `json:"intWithRange" rest:"min=100,max=10"`
	}
	builder = NewBuilder()
	_, err = builder.Build(reflect.TypeOf(S2{}))
	ut.Assert(t, err != nil, "")
}

func TestCheckRequired(t *testing.T) {
	builder := NewBuilder()
	sf, _ := builder.Build(reflect.TypeOf(TestStruct{}))
	ts := TestStruct{
		Name:                  "dd",
		StringWithOption:      "ceph",
		StringWithLenLimit:    "aaa",
		IntWithRange:          100,
		StringIntMap:          map[string]int32{"name": 20},
		IntSlice:              []uint32{1},
		StringSliceWithOption: []MyOption{MyOption("lvm")},
		SliceComposition: []IncludeStruct{
			IncludeStruct{
				Int8WithRange: 5,
			},
		},
		StringMapCompostion: map[string]IncludeStruct{
			"a": IncludeStruct{
				Int8WithRange: 6,
			},
		},
		PtrComposition: &IncludeStruct{
			Int8WithRange: 7,
		},
		SlicePtrComposition: []*IncludeStruct{
			&IncludeStruct{
				Int8WithRange: 5,
			},
		},
		StringPtrMapCompostion: map[string]*IncludeStruct{
			"a": &IncludeStruct{
				Int8WithRange: 5,
			},
		},
	}

	raw := make(map[string]interface{})
	rawByte, _ := json.Marshal(ts)
	json.Unmarshal(rawByte, &raw)
	ut.Assert(t, sf.CheckRequired(raw) == nil, "")

	for _, name := range []string{"name", "stringWithOption", "stringMapComposition", "ptrComposition", "slicePtrComposition", "stringPtrMapComposition", "stringIntMap", "intSlice", "stringSliceWithOption"} {
		json.Unmarshal(rawByte, &raw)
		delete(raw, name)
		ut.Assert(t, sf.CheckRequired(raw) != nil, "")
	}

	ts.Name = ""
	rawByte, _ = json.Marshal(ts)
	json.Unmarshal(rawByte, &raw)
	ut.Assert(t, sf.CheckRequired(raw) != nil, "")
}

func TestCheckNestRequired(t *testing.T) {
	type inner struct {
		Age  int
		Name string `json:"name" rest:"required=true"`
	}

	type wrapper struct {
		Inner inner `json:"name" rest:"required=true"`
	}
	builder := NewBuilder()
	sf, _ := builder.Build(reflect.TypeOf(wrapper{}))
	w := wrapper{
		Inner: inner{
			Age:  10,
			Name: "aaa",
		},
	}
	raw := make(map[string]interface{})
	rawByte, _ := json.Marshal(w)
	json.Unmarshal(rawByte, &raw)
	ut.Assert(t, sf.CheckRequired(raw) == nil, "")

	w.Inner.Name = ""
	raw = make(map[string]interface{})
	rawByte, _ = json.Marshal(w)
	json.Unmarshal(rawByte, &raw)
	ut.Assert(t, sf.CheckRequired(raw) != nil, "")
}

func TestValidate(t *testing.T) {
	builder := NewBuilder()
	sf, _ := builder.Build(reflect.TypeOf(TestStruct{}))

	ts := TestStruct{
		Name:               "dd",
		StringWithOption:   "ceph",
		StringWithLenLimit: "aaa",
		IntWithRange:       100,
		SliceComposition: []IncludeStruct{
			IncludeStruct{
				Int8WithRange: 5,
			},
		},
		StringMapCompostion: map[string]IncludeStruct{
			"a": IncludeStruct{
				Int8WithRange: 6,
			},
		},
		PtrComposition: &IncludeStruct{
			Int8WithRange: 7,
		},
		SlicePtrComposition: []*IncludeStruct{
			nil,
			&IncludeStruct{
				Int8WithRange: 5,
			},
		},
		StringPtrMapCompostion: map[string]*IncludeStruct{
			"a": &IncludeStruct{
				Int8WithRange: 5,
			},
		},
	}

	ut.Assert(t, sf.Validate(ts) == nil, "")

	rawByte, _ := json.Marshal(ts)

	ts.StringWithOption = "oo"
	ut.Assert(t, sf.Validate(ts) != nil, "")

	ts = TestStruct{}
	json.Unmarshal(rawByte, &ts)
	ts.IntWithRange = 10000
	ut.Assert(t, sf.Validate(ts) != nil, "")

	ts = TestStruct{}
	json.Unmarshal(rawByte, &ts)
	ss := ts.StringMapCompostion["a"]
	ss.Int8WithRange = 22
	ts.StringMapCompostion["a"] = ss
	ut.Assert(t, sf.Validate(ts) != nil, "")

	ts = TestStruct{}
	json.Unmarshal(rawByte, &ts)
	ts.PtrComposition.Int8WithRange = 22
	ut.Assert(t, sf.Validate(ts) != nil, "")

	ts = TestStruct{}
	json.Unmarshal(rawByte, &ts)
	ts.SliceComposition[0].Int8WithRange = 22
	ut.Assert(t, sf.Validate(ts) != nil, "")

	ts = TestStruct{}
	json.Unmarshal(rawByte, &ts)
	ts.StringPtrMapCompostion["a"].Int8WithRange = 22
	ut.Assert(t, sf.Validate(ts) != nil, "")
}
