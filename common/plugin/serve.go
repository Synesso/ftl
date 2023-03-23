package plugin

import (
	"context"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/alecthomas/errors"
	"github.com/alecthomas/kong"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/TBD54566975/ftl/common/log"
	"github.com/TBD54566975/ftl/common/socket"
)

type serveCli struct {
	LogConfig log.Config    `embed:"" group:"Logging:"`
	Socket    socket.Socket `help:"Socket to listen on." env:"FTL_PLUGIN_ENDPOINT" required:""`
	kong.Plugins
}

type startOptions[Impl any] struct {
	register []func(grpc.ServiceRegistrar, Impl)
}

type StartOption[Impl any] func(*startOptions[Impl])

// RegisterAdditionalServer allows a plugin to serve additional gRPC services.
//
// "Impl" must be an implementation of "Iface.
func RegisterAdditionalServer[Impl any, Iface any](register func(grpc.ServiceRegistrar, Iface)) StartOption[Impl] {
	return func(so *startOptions[Impl]) {
		so.register = append(so.register, func(sr grpc.ServiceRegistrar, i Impl) {
			register(sr, any(i).(Iface)) //nolint:forcetypeassert
		})
	}
}

// Start a gRPC server plugin listening on the socket specified by the
// environment variable FTL_PLUGIN_ENDPOINT.
//
// This function does not return.
//
// "Config" is Kong configuration to pass to "create".
// "create" is called to create the implementation of the service.
// "register" is called to register the service with the gRPC server and is typically a generated function.
func Start[Impl any, Iface any, Config any](
	ctx context.Context,
	name string,
	create func(context.Context, Config) (Impl, error),
	register func(grpc.ServiceRegistrar, Iface),
	options ...StartOption[Impl],
) {
	var config Config
	cli := serveCli{Plugins: kong.Plugins{&config}}
	kctx := kong.Parse(&cli, kong.Description(`FTL - Towards a 𝝺-calculus for large-scale systems`))

	so := &startOptions[Impl]{}
	for _, option := range options {
		option(so)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	logger := log.Configure(os.Stderr, cli.LogConfig).Sub(name, log.Default)
	ctx = log.ContextWithLogger(ctx, logger)

	logger.Debugf("Starting on %s", cli.Socket)

	// Signal handling.
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigch
		logger.Infof("Terminated by signal %s", sig)
		cancel()
		_ = syscall.Kill(-syscall.Getpid(), sig.(syscall.Signal)) //nolint:forcetypeassert
		os.Exit(0)
	}()

	svc, err := create(ctx, config)
	kctx.FatalIfErrorf(err)

	l, err := socket.Listen(cli.Socket)
	kctx.FatalIfErrorf(err)
	gs := socket.NewGRPCServer(ctx)
	reflection.Register(gs)
	register(gs, any(svc).(Iface)) //nolint:forcetypeassert
	for _, register := range so.register {
		register(gs, svc)
	}
	err = gs.Serve(l)
	kctx.FatalIfErrorf(err)
	kctx.Exit(0)
}

func allocatePort() (*net.TCPAddr, error) {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to allocate port")
	}
	_ = l.Close()
	return l.Addr().(*net.TCPAddr), nil //nolint:forcetypeassert
}

func cleanup(logger *log.Logger, pidFile string) error {
	pidb, err := ioutil.ReadFile(pidFile)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return errors.WithStack(err)
	}
	pid, err := strconv.Atoi(string(pidb))
	if err != nil && !os.IsNotExist(err) {
		return errors.WithStack(err)
	}
	err = syscall.Kill(pid, syscall.SIGKILL)
	if err != nil && !errors.Is(err, syscall.ESRCH) {
		logger.Warnf("Failed to reap old plugin with pid %d: %s", pid, err)
	}
	return nil
}
