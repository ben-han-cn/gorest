package gorest

import (
	"net/http"

	goresterr "github.com/ben-han-cn/gorest/error"
	"github.com/ben-han-cn/gorest/resource"
)

type HandlerFunc func(*resource.Context) *goresterr.APIError
type HandlersChain []HandlerFunc

type Server struct {
	Schemas  resource.SchemaManager
	handlers HandlersChain
}

func NewAPIServer(schemas resource.SchemaManager) *Server {
	return &Server{
		Schemas: schemas,
	}
}

func (s *Server) Use(h HandlerFunc) {
	s.handlers = append(s.handlers, h)
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx, err := resource.NewContext(rw, req, s.Schemas)
	if err != nil {
		WriteResponse(rw, err.Status, err)
		return
	}

	for _, h := range s.handlers {
		if err := h(ctx); err != nil {
			WriteResponse(rw, err.Status, err)
			return
		}
	}

	if err := restHandler(ctx); err != nil {
		WriteResponse(rw, err.Status, err)
	}
}
