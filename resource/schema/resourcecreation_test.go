package schema

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	ut "github.com/zdnscloud/cement/unittest"
)

type podGenJson struct {
	Name                  string
	Count                 uint32            `json:"Count,omitempty"`
	Annotations           map[string]string `json:"Annotations,omitempty"`
	OtherInfoSlice        []OtherPodInfo    `json:"OtherInfoSlice,omitempty"`
	OtherInfoPointer      *OtherPodInfo     `json:"OtherInfoPointer,omitempty"`
	OtherInfoPointerSlice []*OtherPodInfo   `json:"OtherInfoPointerSlice,omitempty"`

	//simulate not pass OtherInfo
	OtherInfo *OtherPodInfo `json:"OtherInfo,omitempty"`
}

func TestFillDefaultValue(t *testing.T) {
	mgr := createSchemaManager()
	testcases := []struct {
		jsonPod   podGenJson
		expectPod Pod
	}{
		{
			podGenJson{
				Name: "p1",
			},
			Pod{
				Name:  "p1",
				Count: 20,
				OtherInfo: OtherPodInfo{
					Name:    "other",
					Numbers: []uint32{1, 2, 3},
				},
				OtherInfoSlice: []OtherPodInfo{
					OtherPodInfo{
						Name:    "other",
						Numbers: []uint32{1, 2, 3},
					},
				},
				OtherInfoPointer: &OtherPodInfo{
					Name:    "other",
					Numbers: []uint32{1, 2, 3},
				},
				OtherInfoPointerSlice: []*OtherPodInfo{
					&OtherPodInfo{
						Name:    "other",
						Numbers: []uint32{1, 2, 3},
					},
				},
			},
		},
		{
			podGenJson{
				Name:  "p2",
				Count: 30,
				OtherInfo: &OtherPodInfo{
					Name: "other1",
				},
			},
			Pod{
				Name:  "p2",
				Count: 30,
				OtherInfo: OtherPodInfo{
					Name: "other1",
				},
				OtherInfoSlice: []OtherPodInfo{
					OtherPodInfo{
						Name:    "other",
						Numbers: []uint32{1, 2, 3},
					},
				},
				OtherInfoPointer: &OtherPodInfo{
					Name:    "other",
					Numbers: []uint32{1, 2, 3},
				},
				OtherInfoPointerSlice: []*OtherPodInfo{
					&OtherPodInfo{
						Name:    "other",
						Numbers: []uint32{1, 2, 3},
					},
				},
			},
		},
		{
			podGenJson{
				Name:  "p3",
				Count: 30,
				OtherInfoPointer: &OtherPodInfo{
					Name: "other1",
				},
			},
			Pod{
				Name:  "p3",
				Count: 30,
				OtherInfo: OtherPodInfo{
					Name:    "other",
					Numbers: []uint32{1, 2, 3},
				},
				OtherInfoSlice: []OtherPodInfo{
					OtherPodInfo{
						Name:    "other",
						Numbers: []uint32{1, 2, 3},
					},
				},
				OtherInfoPointer: &OtherPodInfo{
					Name: "other1",
				},
				OtherInfoPointerSlice: []*OtherPodInfo{
					&OtherPodInfo{
						Name:    "other",
						Numbers: []uint32{1, 2, 3},
					},
				},
			},
		},
		{
			podGenJson{
				Name:  "p4",
				Count: 30,
				OtherInfoSlice: []OtherPodInfo{
					OtherPodInfo{
						Name:    "other1",
						Numbers: []uint32{1, 3},
					},
				},
			},
			Pod{
				Name:  "p4",
				Count: 30,
				OtherInfo: OtherPodInfo{
					Name:    "other",
					Numbers: []uint32{1, 2, 3},
				},
				OtherInfoSlice: []OtherPodInfo{
					OtherPodInfo{
						Name:    "other1",
						Numbers: []uint32{1, 3},
					},
				},
				OtherInfoPointer: &OtherPodInfo{
					Name:    "other",
					Numbers: []uint32{1, 2, 3},
				},
				OtherInfoPointerSlice: []*OtherPodInfo{
					&OtherPodInfo{
						Name:    "other",
						Numbers: []uint32{1, 2, 3},
					},
				},
			},
		},
		{
			podGenJson{
				Name:  "p5",
				Count: 30,
				OtherInfoPointerSlice: []*OtherPodInfo{
					&OtherPodInfo{
						Name:    "other1",
						Numbers: []uint32{1, 3},
					},
				},
			},
			Pod{
				Name:  "p5",
				Count: 30,
				OtherInfo: OtherPodInfo{
					Name:    "other",
					Numbers: []uint32{1, 2, 3},
				},
				OtherInfoSlice: []OtherPodInfo{
					OtherPodInfo{
						Name:    "other",
						Numbers: []uint32{1, 2, 3},
					},
				},
				OtherInfoPointer: &OtherPodInfo{
					Name:    "other",
					Numbers: []uint32{1, 2, 3},
				},
				OtherInfoPointerSlice: []*OtherPodInfo{
					&OtherPodInfo{
						Name:    "other1",
						Numbers: []uint32{1, 3},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		url := "/apis/testing/v1/clusters/c1/namespaces/n1/deployments/d1/pods/"
		reqBody, _ := json.Marshal(tc.jsonPod)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(string(reqBody)))
		r, err := mgr.CreateResourceFromRequest(req)
		ut.Assert(t, err == nil, "")
		ut.Equal(t, r.GetType(), "pod")
		pod := r.(*Pod)
		PodEqual(t, pod, &tc.expectPod)
	}
}

func PodEqual(t *testing.T, c, p *Pod) {
	ut.Equal(t, c.Name, p.Name)
	ut.Equal(t, c.Count, p.Count)
	ut.Equal(t, c.Annotations, p.Annotations)
	ut.Equal(t, c.OtherInfo, p.OtherInfo)
	ut.Equal(t, c.OtherInfoSlice, p.OtherInfoSlice)
	ut.Equal(t, c.OtherInfoPointer, p.OtherInfoPointer)
	ut.Equal(t, c.OtherInfoPointerSlice, p.OtherInfoPointerSlice)
}

func TestAction(t *testing.T) {
	mgr := createSchemaManager()
	url := "/apis/testing/v1/clusters/c1/namespaces/n1/deployments/d1/pods/p1?action=move"
	reqBody, _ := json.Marshal(Location{NodeName: "n1"})
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(string(reqBody)))
	r, err := mgr.CreateResourceFromRequest(req)
	ut.Assert(t, err == nil, "")
	ut.Equal(t, r.GetType(), "pod")
	action := r.GetAction()
	ut.Equal(t, action.Name, "move")
	ut.Equal(t, action.Input.(*Location).NodeName, "n1")

	url = "/apis/testing/v1/clusters/c1/namespaces/n1/deployments/d1/pods/p1?action=me"
	req, _ = http.NewRequest(http.MethodPost, url, nil)
	_, err = mgr.CreateResourceFromRequest(req)
	ut.Assert(t, err != nil, "")
}
