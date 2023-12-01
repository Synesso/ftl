//ftl:module main
package main

import (
	"context"
	"time"

	"github.com/TBD54566975/ftl/go-runtime/sdk"
)

type Nested struct {
}

type Req struct {
	Int      int
	Int64    int64
	Float    float64
	String   string
	Slice    []string
	Map      map[string]string
	Nested   Nested
	Optional sdk.Option[Nested]
	Time     time.Time
}
type Resp struct{}

//ftl:verb
func Verb(ctx context.Context, req Req) (Resp, error) {
	return Resp{}, nil
}