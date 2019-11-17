package gorest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ut "github.com/ben-han-cn/cement/unittest"
	goresterr "github.com/ben-han-cn/gorest/error"
	"github.com/ben-han-cn/gorest/resource"
	"github.com/ben-han-cn/gorest/resource/schema"
)

var (
	schemas = schema.NewSchemaManager()
	version = resource.APIVersion{
		Group:   "testing",
		Version: "v1",
	}
)

type dumbHandler struct{}

func (h *dumbHandler) Create(ctx *resource.Context) (resource.Resource, *goresterr.APIError) {
	return nil, nil
}

func (h *dumbHandler) List(ctx *resource.Context) interface{} {
	return nil
}

type Foo struct {
	resource.ResourceBase
}

var gnum int

var dumbHandler1 = func(ctx *resource.Context) *goresterr.APIError {
	ctx.Set("key", &gnum)
	return nil
}

var dumbHandler2 = func(ctx *resource.Context) *goresterr.APIError {
	val_, _ := ctx.Get("key")
	*(val_.(*int)) = 100
	return nil
}

func TestContextPassChain(t *testing.T) {
	schemas.Import(&version, Foo{}, &dumbHandler{})
	req, _ := http.NewRequest("GET", "/apis/testing/v1/foos", nil)
	req.Host = "127.0.0.1:1234"
	w := httptest.NewRecorder()
	s := NewAPIServer(schemas)
	s.Use(dumbHandler1)
	s.Use(dumbHandler2)

	ut.Equal(t, gnum, 0)
	s.ServeHTTP(w, req)
	ut.Equal(t, gnum, 100)

	s.ServeHTTP(w, req)
}
