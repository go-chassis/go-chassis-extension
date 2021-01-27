package gin

import (
	"net/http"

	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/server"
)

// ResourceHandler wraps go-chassis restful function
type ResourceHandler struct {
	handleFunc func(ctx *Context)
	rc         *Context
	opts       server.Options
}

// Handle is to handle the router related things
func (h *ResourceHandler) Handle(chain *handler.Chain, inv *invocation.Invocation, cb invocation.ResponseCallBack) {
	Invocation2HTTPRequest(inv, h.rc.C)

	// check body size
	if h.opts.BodyLimit > 0 {
		h.rc.C.Request.Body = http.MaxBytesReader(h.rc.C.Writer, h.rc.C.Request.Body, h.opts.BodyLimit)
	}

	h.rc.Ctx = inv.Ctx
	// call real route func
	h.handleFunc(h.rc)
	ir := &invocation.Response{}
	ir.Status = h.rc.C.Writer.Status()
	ir.Result = h.rc.C.Writer
	//call next chain
	cb(ir)
}

func newHandler(f func(ctx *Context), rc *Context, opts server.Options) handler.Handler {
	return &ResourceHandler{
		handleFunc: f,
		rc:         rc,
		opts:       opts,
	}
}

// Name returns the name string
func (h *ResourceHandler) Name() string {
	return "gin"
}
