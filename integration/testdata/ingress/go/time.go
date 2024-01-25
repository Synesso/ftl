//ftl:module time
package _go

import (
	"context"
	"time"
)

type TimeRequest struct {
	Message string `alias:"m"`
}
type TimeResponse struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

// Time returns the current time.
//
//ftl:verb
//ftl:ingress GET /time
func Time(ctx context.Context, req TimeRequest) (TimeResponse, error) {
	return TimeResponse{
		Message: req.Message,
		Time:    time.Now(),
	}, nil
}
