package resourcefield

import (
	"reflect"
	"strings"
	"testing"

	ut "github.com/ben-han-cn/cement/unittest"
)

func TestBuildValidatorWithValidTag(t *testing.T) {
	sv := reflect.ValueOf(TestStruct{
		StringWithOption:      "lvm",
		StringWithLenLimit:    "good",
		IntWithRange:          10,
		StringSliceWithOption: []MyOption{"ceph"},
	})

	sv2 := reflect.ValueOf(TestStruct{
		StringWithOption:      "lvms",
		StringWithLenLimit:    "g",
		IntWithRange:          10000,
		StringSliceWithOption: []MyOption{"ceph", "xxx"},
	})
	testFields := []string{
		"StringWithOption",
		"StringWithLenLimit",
		"IntWithRange",
		"StringSliceWithOption",
	}
	st := sv.Type()
	for _, fn := range testFields {
		f, ok := st.FieldByName(fn)
		ut.Assert(t, ok, "field %s doesn't exist", fn)
		tags := strings.Split(f.Tag.Get("rest"), ",")
		if len(tags) > 0 {
			validator, err := buildValidator(f.Type, tags)
			ut.Assert(t, err == nil, "get err %v", err)
			if validator != nil {
				ut.Assert(t, validator.Validate(sv.FieldByName(fn).Interface()) == nil, "")
				ut.Assert(t, validator.Validate(sv2.FieldByName(fn).Interface()) != nil, "")
			}
		}
	}
}

func TestBuildValidatorWithInValidTag(t *testing.T) {
	type testStruct struct {
		IntWithOption      int    `rest:"required=true,options=lvm|ceph"`
		IntWithLenLimit    int    `rest:"minLen=10,maxLen=11"`
		StringWithLenLimit string `rest:"minLen=12,maxLen=12"`
		ShortOfMax         uint32 `rest:"min=1"`
		ShortOfMin         int8   `rest:"max=1"`
	}

	sv := reflect.ValueOf(testStruct{})
	st := sv.Type()
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		tags := strings.Split(f.Tag.Get("rest"), ",")
		if len(tags) > 0 {
			validator, _ := buildValidator(f.Type, tags)
			ut.Assert(t, validator == nil, "")
		}
	}
}
