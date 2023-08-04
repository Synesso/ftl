package rpc

import (
	"context"
	"net/http"

	"github.com/alecthomas/errors"
	"github.com/bufbuild/connect-go"
	otelconnect "github.com/bufbuild/connect-opentelemetry-go"

	"github.com/TBD54566975/ftl/backend/common/log"
	"github.com/TBD54566975/ftl/backend/common/model"
	"github.com/TBD54566975/ftl/backend/common/rpc/headers"
	"github.com/TBD54566975/ftl/backend/schema"
)

type ftlDirectRoutingKey struct{}
type ftlVerbKey struct{}
type requestIDKey struct{}

// WithDirectRouting ensures any hops in Verb routing do not redirect.
//
// This is used so that eg. calls from Drives do not create recursive loops
// when calling back to the Agent.
func WithDirectRouting(ctx context.Context) context.Context {
	return context.WithValue(ctx, ftlDirectRoutingKey{}, "1")
}

// WithVerbs adds the module.verb chain from the current request to the context.
func WithVerbs(ctx context.Context, verbs []*schema.VerbRef) context.Context {
	return context.WithValue(ctx, ftlVerbKey{}, verbs)
}

// VerbFromContext returns the current module.verb of the current request.
func VerbFromContext(ctx context.Context) (*schema.VerbRef, bool) {
	value := ctx.Value(ftlVerbKey{})
	verbs, ok := value.([]*schema.VerbRef)
	if len(verbs) == 0 {
		return nil, false
	}
	return verbs[len(verbs)-1], ok
}

// VerbsFromContext returns the module.verb chain of the current request.
func VerbsFromContext(ctx context.Context) ([]*schema.VerbRef, bool) {
	value := ctx.Value(ftlVerbKey{})
	verbs, ok := value.([]*schema.VerbRef)
	return verbs, ok
}

// IsDirectRouted returns true if the incoming request should be directly
// routed and never redirected.
func IsDirectRouted(ctx context.Context) bool {
	return ctx.Value(ftlDirectRoutingKey{}) != nil
}

// RequestKeyFromContext returns the request Key from the context, if any.
func RequestKeyFromContext(ctx context.Context) (model.IngressRequestKey, bool, error) {
	value := ctx.Value(requestIDKey{})
	keyStr, ok := value.(string)
	if !ok {
		return model.IngressRequestKey{}, false, nil
	}
	key, err := model.ParseIngressRequestKey(keyStr)
	if err != nil {
		return model.IngressRequestKey{}, false, errors.Wrap(err, "invalid request Key")
	}
	return key, true, nil
}

// WithRequestKey adds the request Key to the context.
func WithRequestKey(ctx context.Context, key model.IngressRequestKey) context.Context {
	return context.WithValue(ctx, requestIDKey{}, key.String())
}

func DefaultClientOptions(level log.Level) []connect.ClientOption {
	return []connect.ClientOption{
		connect.WithGRPC(), // Use gRPC because some servers will not be using Connect.
		connect.WithInterceptors(MetadataInterceptor(level)),
	}
}

func DefaultHandlerOptions() []connect.HandlerOption {
	return []connect.HandlerOption{
		connect.WithInterceptors(MetadataInterceptor(log.Error)),
		connect.WithInterceptors(otelconnect.NewInterceptor()),
	}
}

// MetadataInterceptor propagates FTL metadata through servers and clients.
func MetadataInterceptor(level log.Level) connect.Interceptor {
	return &metadataInterceptor{
		errorLevel: level,
	}
}

type metadataInterceptor struct {
	errorLevel log.Level
}

func (*metadataInterceptor) WrapStreamingClient(req connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, s connect.Spec) connect.StreamingClientConn {
		// TODO(aat): I can't figure out how to get the client headers here.
		logger := log.FromContext(ctx)
		logger.Tracef("%s (streaming client)", s.Procedure)
		return req(ctx, s)
	}
}

func (m *metadataInterceptor) WrapStreamingHandler(req connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, s connect.StreamingHandlerConn) error {
		logger := log.FromContext(ctx)
		logger.Tracef("%s (streaming handler)", s.Spec().Procedure)
		ctx, err := propagateHeaders(ctx, s.Spec().IsClient, s.RequestHeader())
		if err != nil {
			return err
		}
		err = errors.WithStack(req(ctx, s))
		if err != nil {
			logger.Logf(m.errorLevel, "Streaming RPC failed: %s: %s", err, s.Spec().Procedure)
			return err
		}
		return nil
	}
}

func (m *metadataInterceptor) WrapUnary(uf connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		logger := log.FromContext(ctx)
		logger.Tracef("%s (unary)", req.Spec().Procedure)
		ctx, err := propagateHeaders(ctx, req.Spec().IsClient, req.Header())
		if err != nil {
			return nil, err
		}
		resp, err := uf(ctx, req)
		if err != nil {
			err = errors.WithStack(err)
			logger.Logf(m.errorLevel, "Unary RPC failed: %s: %s", err, req.Spec().Procedure)
			return nil, err
		}
		return resp, nil
	}
}

type clientKey[Client Pingable] struct{}

// ContextWithClient returns a context with an RPC client attached.
func ContextWithClient[Client Pingable](ctx context.Context, client Client) context.Context {
	return context.WithValue(ctx, clientKey[Client]{}, client)
}

// ClientFromContext returns the given RPC client from the context, or panics.
func ClientFromContext[Client Pingable](ctx context.Context) Client {
	value := ctx.Value(clientKey[Client]{})
	if value == nil {
		panic("no RPC client in context")
	}
	return value.(Client) //nolint:forcetypeassert
}

func propagateHeaders(ctx context.Context, isClient bool, header http.Header) (context.Context, error) {
	if isClient {
		if IsDirectRouted(ctx) {
			headers.SetDirectRouted(header)
		}
		if verbs, ok := VerbsFromContext(ctx); ok {
			headers.SetCallers(header, verbs)
		}
		if key, ok, err := RequestKeyFromContext(ctx); ok {
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if ok {
				headers.SetRequestKey(header, key)
			}
		}
	} else {
		if headers.IsDirectRouted(header) {
			ctx = WithDirectRouting(ctx)
		}
		if verbs, err := headers.GetCallers(header); err != nil {
			return nil, errors.WithStack(err)
		} else { //nolint:revive
			ctx = WithVerbs(ctx, verbs)
		}
		if key, ok, err := headers.GetRequestKey(header); err != nil {
			return nil, errors.WithStack(err)
		} else if ok {
			ctx = WithRequestKey(ctx, key)
		}
	}
	return ctx, nil
}