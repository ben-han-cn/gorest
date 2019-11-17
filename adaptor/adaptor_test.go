package adaptor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	ut "github.com/ben-han-cn/cement/unittest"
	"github.com/ben-han-cn/gorest/resource"
)

type testStruct struct {
	T *testing.T
}

func (t *testStruct) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ut.Equal(t.T, req.Method, "POST")
	ut.Equal(t.T, req.URL.Path, "/path")
	w.WriteHeader(201)
	fmt.Fprint(w, "hello")
}

func (t *testStruct) route() resource.ResourceRoute {
	route := resource.NewResourceRoute()
	route.AddPathForMethod(http.MethodPost, "/path")
	return route
}

func doRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestRegisterHandler(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	ts := &testStruct{}
	RegisterHandler(router, ts, ts.route())

	w := doRequest(router, "POST", "/path")
	ut.Equal(t, w.Code, 201)
	ut.Equal(t, w.Body.String(), "hello")
}
