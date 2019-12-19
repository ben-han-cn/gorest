package resource

import (
	"encoding/json"
	"testing"

	ut "github.com/zdnscloud/cement/unittest"
)

func TestCollectionToJson(t *testing.T) {
	r := &dumbResource{
		Number: 10,
	}

	rs, err := NewResourceCollection(r, nil)
	ut.Assert(t, err == nil, "")
	ut.Assert(t, rs.Resources != nil, "")
	d, _ := json.Marshal(rs)
	ut.Equal(t, string(d), `{"type":"collection","data":[]}`)

	rs2, err := NewResourceCollection(r, []*dumbResource{})
	ut.Assert(t, err == nil, "")
	ut.Assert(t, rs2.Resources != nil, "")
	ut.Equal(t, len(rs2.Resources), 0)
	d2, _ := json.Marshal(rs2)
	ut.Equal(t, string(d), string(d2))
}
