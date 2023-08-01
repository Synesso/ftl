// This is the echo module.
//
//ftl:module echo
package echo

import (
	"context"
	"fmt"
	"time"

	timemodule "github.com/TBD54566975/ftl/examples/time"
	"github.com/TBD54566975/ftl/internal/log"

	ftl "github.com/TBD54566975/ftl/go-runtime/sdk"
)

// An echo request.
type EchoRequest struct {
	Name string `json:"name"`
}

type EchoResponse struct {
	Message string `json:"message"`
}

// Echo returns a greeting with the current time.
//
//ftl:verb
//ftl:ingress GET /echo
func Echo(ctx context.Context, req EchoRequest) (EchoResponse, error) {
	logger := log.FromContext(ctx)

	logger.Errorf(nil, "Received a request!")

	tresp, err := ftl.Call(ctx, timemodule.Time, timemodule.TimeRequest{})
	if err != nil {
		return EchoResponse{}, err
	}
	t := time.Unix(int64(tresp.Time), 0)
	return EchoResponse{Message: fmt.Sprintf("Hello, %s!!! It is %s!", req.Name, t)}, nil
}
