package resourcefield

import (
	"encoding/json"
	ut "github.com/zdnscloud/cement/unittest"
	"reflect"
	"strings"
	"testing"
)

func TestFieldBuild(t *testing.T) {
	type Embed struct {
		Id  string `rest:"required=true"`
		Age int64  `rest:"min=10,max=20"`
	}

	type IncludeStruct struct {
		Int8WithRange     int8   `json:"int8WithRange" rest:"min=1,max=20"`
		Uint16WithDefault uint16 `json:"uint16WithDefault,omitempty"`
	}

	type MyOption string

	type TestStruct struct {
		Embed            `json:",inline"`
		Name             string   `json:"name" rest:"required=true"`
		StringWithOption MyOption `json:"stringWithOption,omitempty" rest:"required=true,options=lvm|ceph"`
		IntWithDefault   int      `json:"intWithDefault,omitempty"`
		IntWithRange     uint32   `json:"intWithRange" rest:"min=1,max=1000"`
		BoolWithDefault  bool     `json:"boolWithDefault,omitempty"`

		IntSlice              []uint32         `json:"intSlice,omitempty" rest:"required=true"`
		StringSliceWithOption []MyOption       `json:"stringSliceWithOption,omitempty" rest:"required=true,options=lvm|ceph"`
		SliceStruct           []IncludeStruct  `json:"sliceStruct" rest:"required=true"`
		SlicePtrStruct        []*IncludeStruct `json:"slicePtrStruct" rest:"required=true"`

		StringMapStruct    map[string]IncludeStruct  `json:"stringMapStruct" rest:"required=true"`
		StringPtrMapStruct map[string]*IncludeStruct `json:"stringPtrMapStruct" rest:"required=true"`
		StringIntMap       map[string]int32          `json:"stringIntMap,omitempty" rest:"required=true"`

		PtrStruct *IncludeStruct `json:"ptrStruct" rest:"required=true"`
	}

	builder := NewBuilder()
	sf, err := builder.Build(reflect.TypeOf(TestStruct{}))
	ut.Assert(t, err == nil, "")

	fieldNames := []string{
		"Id",
		"Age",
		"Name",
		"StringWithOption",
		"IntWithRange",

		"IntSlice",
		"StringSliceWithOption",
		"SliceStruct",
		"SlicePtrStruct",

		"StringMapStruct",
		"StringPtrMapStruct",
		"StringIntMap",

		"PtrStruct",
	}
	ut.Equal(t, len(sf.fields), len(fieldNames))
	for _, name := range fieldNames {
		_, ok := sf.fields[name]
		ut.Assert(t, ok, "%s has no field", name)
	}
}

//for struct without any rest contraint generate nil field
func TestFieldBuildForNoneRestStruct(t *testing.T) {
	type IncludeStruct struct {
		Int8WithRange     int8   `json:"int8WithRange"`
		Uint16WithDefault uint16 `json:"uint16WithDefault,omitempty"`
	}
	builder := NewBuilder()
	sf, _ := builder.Build(reflect.TypeOf(IncludeStruct{}))
	ut.Assert(t, sf == nil, "")

	type TestStruct struct {
		SliceStruct      []IncludeStruct
		StringMapStruct  map[string]IncludeStruct
		IncludeStructPtr *IncludeStruct
	}
	builder = NewBuilder()
	sf, err := builder.Build(reflect.TypeOf(TestStruct{}))
	ut.Assert(t, sf == nil, "")
	ut.Assert(t, err == nil, "")
}

func TestCheckRequired(t *testing.T) {
	type IncludeStruct struct {
		Int8WithRange int8           `json:"int8WithRange,omitempty" rest:"min=1,max=20"`
		StringSlice   []string       `rest:"required=true"`
		StringIntMap  map[string]int `rest:"required=true"`
	}

	type TestStruct struct {
		Name           string `json:"name" rest:"required=true"`
		IntWithDefault int    `json:"int" rest:"required=true"`

		IntSlice       []uint32         `json:"intSlice,omitempty"`
		SliceStruct    []IncludeStruct  `json:"sliceStruct" rest:"required=true"`
		SlicePtrStruct []*IncludeStruct `json:"slicePtrStruct" rest:"required=true"`

		StringMapStruct    map[string]IncludeStruct  `json:"stringMapStruct" rest:"required=true"`
		StringPtrMapStruct map[string]*IncludeStruct `json:"stringPtrMapStruct"`
		StringIntMap       map[string]int32          `json:"stringIntMap,omitempty" rest:"required=true"`

		PtrStruct *IncludeStruct `json:"ptrStruct" rest:"required=true"`
	}

	builder := NewBuilder()
	sf, _ := builder.Build(reflect.TypeOf(TestStruct{}))
	ut.Equal(t, len(sf.fields), 8)
	ts := TestStruct{
		Name:           "dd",
		IntWithDefault: 100,
		IntSlice:       []uint32{1},
		SliceStruct: []IncludeStruct{
			IncludeStruct{
				Int8WithRange: 1,
				StringSlice:   []string{"a"},
				StringIntMap: map[string]int{
					"b": 2,
				},
			},
		},
		SlicePtrStruct: []*IncludeStruct{
			&IncludeStruct{
				Int8WithRange: 3,
				StringSlice:   []string{"d"},
				StringIntMap: map[string]int{
					"e": 5,
				},
			},
		},

		StringMapStruct: map[string]IncludeStruct{
			"f": IncludeStruct{
				Int8WithRange: 6,
				StringSlice:   []string{"g"},
				StringIntMap: map[string]int{
					"h": 7,
				},
			},
		},
		StringPtrMapStruct: map[string]*IncludeStruct{
			"i": &IncludeStruct{
				Int8WithRange: 8,
				StringSlice:   []string{"j"},
				StringIntMap: map[string]int{
					"k": 9,
				},
			},
		},
		StringIntMap: map[string]int32{
			"l": 10,
		},

		PtrStruct: &IncludeStruct{
			Int8WithRange: 11,
			StringSlice:   []string{"j"},
			StringIntMap: map[string]int{
				"m": 12,
			},
		},
	}

	raw := make(map[string]interface{})
	rawByte, _ := json.Marshal(ts)
	json.Unmarshal(rawByte, &raw)
	err := sf.Validate(ts, raw)
	ut.Assert(t, err == nil, "required check should pass but get %v", err)

	//delete required in map will cause error
	for _, name := range []string{"name", "int", "sliceStruct", "slicePtrStruct", "stringMapStruct", "stringIntMap", "ptrStruct"} {
		raw = make(map[string]interface{})
		json.Unmarshal(rawByte, &raw)
		delete(raw, name)
		ut.Assert(t, sf.Validate(ts, raw) != nil, "delete %s should failed", name)
	}

	//empty array will cause error
	for _, name := range []string{"sliceStruct", "slicePtrStruct", "stringMapStruct", "stringIntMap"} {
		raw = make(map[string]interface{})
		json.Unmarshal(rawByte, &raw)
		raw[name] = nil
		ut.Assert(t, sf.Validate(ts, raw) != nil, "set %s to nil should failed", name)
	}

	//empty array will non-required field should be ok
	for _, name := range []string{"stringPtrMapStruct", "intSlice"} {
		raw = make(map[string]interface{})
		json.Unmarshal(rawByte, &raw)
		raw[name] = nil
		ut.Assert(t, sf.Validate(ts, raw) == nil, "unrequired field %s shouldn't fail ", name)
	}

	//check nest required
	var tmp TestStruct
	json.Unmarshal(rawByte, &tmp)
	tmp.SliceStruct = []IncludeStruct{
		IncludeStruct{
			StringSlice:  []string{"c", "d"},
			StringIntMap: make(map[string]int),
		},
	}
	makeSureValidateFailedWithInfo(t, sf, tmp, "StringIntMap")

	tmp = TestStruct{}
	json.Unmarshal(rawByte, &tmp)
	tmp.StringMapStruct = map[string]IncludeStruct{
		"a": IncludeStruct{
			StringIntMap: map[string]int{"a": 20},
		},
	}
	makeSureValidateFailedWithInfo(t, sf, tmp, "StringSlice")
}

func TestValidate(t *testing.T) {
	type IncludeStruct struct {
		Int8WithRange     int8   `json:"int8WithRange" rest:"min=1,max=20"`
		Uint16WithDefault uint16 `json:"uint16WithDefault,omitempty"`
	}

	type MyOption string

	type TestStruct struct {
		IntSlice              []uint32   `json:"intSlice,omitempty" rest:"required=true,min=20,max=30"`
		StringSliceWithOption []MyOption `json:"stringSliceWithOption,omitempty" rest:"required=true,options=lvm|ceph"`

		SliceStruct    []IncludeStruct  `json:"sliceStruct" rest:"required=true"`
		SlicePtrStruct []*IncludeStruct `json:"slicePtrStruct" rest:"required=true"`

		StringMapStruct map[string]IncludeStruct `json:"stringMapStruct" rest:"required=true"`
		StringStringMap map[string]string        `json:"stringStringMap,omitempty" rest:"isDomain=true"`

		PtrStruct *IncludeStruct `json:"ptrStruct" rest:"required=true"`
	}

	builder := NewBuilder()
	sf, err := builder.Build(reflect.TypeOf(TestStruct{}))
	ut.Assert(t, err == nil, "")

	ts := TestStruct{
		IntSlice:              []uint32{21, 22},
		StringSliceWithOption: []MyOption{"lvm", "ceph"},
		SliceStruct: []IncludeStruct{
			IncludeStruct{
				Int8WithRange: 5,
			},
		},
		SlicePtrStruct: []*IncludeStruct{
			&IncludeStruct{
				Int8WithRange: 6,
			},
		},
		StringMapStruct: map[string]IncludeStruct{
			"a": IncludeStruct{
				Int8WithRange: 7,
			},
		},
		StringStringMap: map[string]string{
			"b": "good",
		},
		PtrStruct: &IncludeStruct{
			Int8WithRange: 5,
		},
	}
	rawByte, _ := json.Marshal(ts)
	raw := make(map[string]interface{})
	json.Unmarshal(rawByte, &raw)
	ut.Assert(t, sf.Validate(ts, raw) == nil, "")

	//intslice with min-max
	var tmp TestStruct
	json.Unmarshal(rawByte, &tmp)
	tmp.IntSlice[0] = 32
	makeSureValidateFailedWithInfo(t, sf, tmp, "exceed the range limit")

	//stringslice with option
	tmp = TestStruct{}
	json.Unmarshal(rawByte, &tmp)
	tmp.StringSliceWithOption[0] = "gogod"
	makeSureValidateFailedWithInfo(t, sf, tmp, "included in options")

	//stringslice with min-max
	tmp = TestStruct{}
	json.Unmarshal(rawByte, &tmp)
	tmp.SlicePtrStruct[0].Int8WithRange = 0
	makeSureValidateFailedWithInfo(t, sf, tmp, "exceed the range limit")

	//stringmap with isDomain
	tmp = TestStruct{}
	json.Unmarshal(rawByte, &tmp)
	tmp.StringStringMap["c"] = "Xgoo"
	makeSureValidateFailedWithInfo(t, sf, tmp, "subdomain must consist")
}

func TestValidateWithOptionalField(t *testing.T) {
	type IncludeStruct struct {
		Int8WithRange     int8   `json:"int8WithRange" rest:"min=1,max=20"`
		Uint16WithDefault uint16 `json:"uint16WithDefault,omitempty"`
	}

	type TestStruct struct {
		Name        string          `json:"name,omitempty" rest:"isDomain=true"`
		SliceStruct []IncludeStruct `json:"sliceStruct" rest:"required=true"`
	}

	builder := NewBuilder()
	sf, err := builder.Build(reflect.TypeOf(TestStruct{}))
	ut.Assert(t, err == nil, "")

	ts := TestStruct{
		SliceStruct: []IncludeStruct{
			IncludeStruct{
				Int8WithRange: 5,
			},
		},
	}
	rawByte, _ := json.Marshal(ts)
	raw := make(map[string]interface{})
	json.Unmarshal(rawByte, &raw)
	err = sf.Validate(ts, raw)
	ut.Assert(t, err == nil, "shouldn't get err %v", err)
}

func makeSureValidateFailedWithInfo(t *testing.T, sf Field, structVal interface{}, errorInfo string) {
	rawByte, _ := json.Marshal(structVal)
	raw := make(map[string]interface{})
	json.Unmarshal(rawByte, &raw)
	err := sf.Validate(structVal, raw)
	ut.Assert(t, err != nil, "want err %s", errorInfo)
	ut.Assert(t, strings.Contains(err.Error(), errorInfo), "")
}
