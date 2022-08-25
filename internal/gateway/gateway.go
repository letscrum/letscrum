package gateway

import (
	"context"
	"github.com/letscrum/letscrum/pkg/logging"
	"mime"
	"net/http"

	"github.com/daocloud/skoala/api/hive/v1alpha1"
	_ "github.com/daocloud/skoala/statik"
	"github.com/rakyll/statik/fs"
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Endpoint struct {
	Network, Addr string
}

type Options struct {
	Addr string

	GRPCServer Endpoint

	OpenAPIDir string

	Mux []gwruntime.ServeMuxOption
}

func newGateway(ctx context.Context, conn *grpc.ClientConn, opts []gwruntime.ServeMuxOption) (http.Handler, error) {

	mux := gwruntime.NewServeMux(opts...)

	for _, f := range []func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error{
		v1alpha1.RegisterHiveHandler,
		v1alpha1.RegisterRegistrationHandler,
		v1alpha1.RegisterBookHandler,
	} {
		if err := f(ctx, mux, conn); err != nil {
			return nil, err
		}
	}
	return mux, nil
}

func Run(ctx context.Context, opts Options) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := dialTCP(ctx, opts.GRPCServer.Addr)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			logging.L(ctx).Errorf("Failed to close a client connection to the gRPC server: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/openapiv2/", openAPIServer(opts.OpenAPIDir))
	mux.HandleFunc("/grpcHealthz", grpcHealthzServer(conn))
	mux.Handle("/sys/", runHealthCheck())
	mime.AddExtensionType(".svg", "image/svg+xml")
	statikFS, err := fs.New()
	if err != nil {
		return err
	}
	fileServer := http.FileServer(statikFS)
	mux.Handle("/openapi-ui/", http.StripPrefix("/openapi-ui/", fileServer))

	gw, err := newGateway(ctx, conn, opts.Mux)
	if err != nil {
		return err
	}
	mux.Handle("/", gw)

	s := &http.Server{
		Addr:    opts.Addr,
		Handler: allowCORS(mux),
	}
	go func() {
		<-ctx.Done()
		logging.L(ctx).Infof("Shutting down the http server")
		if err := s.Shutdown(context.Background()); err != nil {
			logging.L(ctx).Errorf("Failed to shutdown http server: %v", err)
		}
	}()

	logging.L(ctx).Infof("Starting listening at %s", opts.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logging.L(ctx).Errorf("Failed to listen and serve: %v", err)
		return err
	}
	return nil
}

func dialTCP(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
