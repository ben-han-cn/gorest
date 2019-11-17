package resourcefield

type Embed struct {
	Id  string `rest:"default=xxxx"`
	Age int64  `rest:"default=20"`
}

type IncludeStruct struct {
	Int8WithRange     int8   `json:"int8WithRange" rest:"min=1,max=20"`
	Uint16WithDefault uint16 `json:"uint16WithDefault,omitempty"`
}

type MyOption string

type TestStruct struct {
	Embed `json:",inline"`

	Name                  string           `json:"name" rest:"required=true"`
	StringWithOption      MyOption         `json:"stringWithOption,omitempty" rest:"required=true,options=lvm|ceph"`
	StringWithDefault     string           `json:"stringWithDefault,omitempty"`
	StringWithLenLimit    string           `json:"stringWithLenLimit" rest:"minLen=2,maxLen=10"`
	IntWithDefault        int              `json:"intWithDefault,omitempty"`
	IntWithRange          uint32           `json:"intWithRange" rest:"min=1,max=1000"`
	BoolWithDefault       bool             `json:"boolWithDefault,omitempty"`
	StringIntMap          map[string]int32 `json:"stringIntMap,omitempty" rest:"required=true"`
	IntSlice              []uint32         `json:"intSlice,omitempty" rest:"required=true"`
	StringSliceWithOption []MyOption       `json:"stringSliceWithOption,omitempty" rest:"required=true,options=lvm|ceph"`

	SliceComposition    []IncludeStruct          `json:"sliceComposition" rest:"required=true"`
	StringMapCompostion map[string]IncludeStruct `json:"stringMapComposition" rest:"required=true"`

	PtrComposition         *IncludeStruct            `json:"ptrComposition" rest:"required=true"`
	SlicePtrComposition    []*IncludeStruct          `json:"slicePtrComposition" rest:"required=true"`
	StringPtrMapCompostion map[string]*IncludeStruct `json:"stringPtrMapComposition" rest:"required=true"`
}
