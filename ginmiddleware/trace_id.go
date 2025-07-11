package ginmiddleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mysteriumnetwork/msk/ctxl"
)

const XHeaderTraceID = "X-Trace-Id"

func setTraceID(ctx context.Context, traceId string) context.Context {
	if traceId == "" {
		return ctxl.SetField(ctx, XHeaderTraceID, uuid.New().String())
	}
	return ctxl.SetField(ctx, XHeaderTraceID, traceId)
}

// InboundTraceIDToContext - searches for XHeaderTraceID in request header and if found set's it to c.Request.Context
// if no XHeaderTraceID header is presents generates a new trace id (uuid) and sets it to c.Request.Context
// should be used as the first middleware in the chain
func InboundTraceIDToContext(c *gin.Context) {
	c.Request = c.Request.WithContext(setTraceID(c.Request.Context(), c.GetHeader(XHeaderTraceID)))
	c.Next()
}
