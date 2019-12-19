package resource

import (
	"testing"

	ut "github.com/zdnscloud/cement/unittest"
)

type Deployment struct {
	ResourceBase
}

func (d Deployment) Default() Resource {
	return &Deployment{}
}

func TestKindAndResourceName(t *testing.T) {
	ut.Equal(t, DefaultKindName(Deployment{}), "deployment")
	ut.Equal(t, DefaultResourceName(Deployment{}), "deployments")
}
