package gin

import (
	"context"

	"github.com/gin-gonic/gin"
)

//Context is a struct which has both request and response objects
// and request context
type Context struct {
	Ctx context.Context
	C   *gin.Context
}

//NewBaseServer is a function which return context
func NewBaseServer(ctx context.Context) *Context {
	return &Context{
		Ctx: ctx,
	}
}
