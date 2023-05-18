package rpc

import (
	"context"

	"github.com/bufbuild/connect-go"

	"github.com/TBD54566975/ftl/common/log"
)

const ftlDirectRoutingHeader = "FTL-Direct"

type ftlDirectRoutingKey struct{}

// WithDirectRouting ensures any hops in Verb routing do not redirect.
//
// This is used so that eg. calls from Drives do not create recursive loops
// when calling back to the Agent.
func WithDirectRouting(ctx context.Context) context.Context {
	return context.WithValue(ctx, ftlDirectRoutingKey{}, "1")
}

// IsDirectRouted returns true if the incoming request should be directly
// routed and never redirected.
func IsDirectRouted(ctx context.Context) bool {
	return ctx.Value(ftlDirectRoutingKey{}) != nil
}

func DefaultClientOptions() []connect.ClientOption {
	return []connect.ClientOption{
		connect.WithInterceptors(MetadataInterceptor()),
	}
}

func DefaultHandlerOptions() []connect.HandlerOption {
	return []connect.HandlerOption{
		connect.WithInterceptors(MetadataInterceptor()),
	}
}

// MetadataInterceptor propagates FTL metadata through servers and clients.
func MetadataInterceptor() connect.Interceptor {
	return &metadataInterceptor{}
}

type metadataInterceptor struct{}

func (*metadataInterceptor) WrapStreamingClient(req connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, s connect.Spec) connect.StreamingClientConn {
		log.FromContext(ctx).Tracef("%s (streaming client)", s.Procedure)
		return req(ctx, s)
	}
}

func (*metadataInterceptor) WrapStreamingHandler(req connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, s connect.StreamingHandlerConn) error {
		log.FromContext(ctx).Tracef("%s (streaming handler)", s.Spec().Procedure)
		if s.Spec().IsClient {
			if IsDirectRouted(ctx) {
				s.RequestHeader().Set(ftlDirectRoutingHeader, "1")
			}
		} else if s.RequestHeader().Get(ftlDirectRoutingHeader) != "" {
			ctx = WithDirectRouting(ctx)
		}
		return req(ctx, s)
	}
}

func (*metadataInterceptor) WrapUnary(uf connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		log.FromContext(ctx).Tracef("%s (unary)", req.Spec().Procedure)
		if req.Spec().IsClient {
			if IsDirectRouted(ctx) {
				req.Header().Set(ftlDirectRoutingHeader, "1")
			}
		} else if req.Header().Get(ftlDirectRoutingHeader) != "" {
			ctx = WithDirectRouting(ctx)
		}
		return uf(ctx, req)
	})
}