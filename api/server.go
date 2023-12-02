package api

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	testv1 "github.com/koo04/gateway-test/internal/gen/proto/go/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type server struct{}

var srv *server = &server{}

func Start() error {
	slog.Info("Starting Service")

	// start the REST proxy endpoints
	go func() {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		mux := runtime.NewServeMux(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
				Marshaler: &runtime.JSONPb{
					MarshalOptions: protojson.MarshalOptions{
						UseProtoNames:   true,
						EmitUnpopulated: true,
					},
					UnmarshalOptions: protojson.UnmarshalOptions{
						DiscardUnknown: true,
					},
				},
			}),
		)
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		// register contact types
		if err := testv1.RegisterTestAPIServiceHandlerFromEndpoint(ctx, mux, ":9000", opts); err != nil {
			slog.Error("failed to register handlers", "error", err)
			return
		}

		httpmux := http.NewServeMux()

		httpmux.Handle("/", mux)

		s := &http.Server{
			Addr:    ":8000",
			Handler: testHandler(httpmux),
		}

		slog.Info("REST server listening on port :8000")
		if err := s.ListenAndServe(); err != nil {
			slog.Error("failed to serve proxy", "error", err)
			return
		}
	}()

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		slog.Error("failed to listen", "error", err)
		return err
	}

	s := grpc.NewServer()
	testv1.RegisterTestAPIServiceServer(s, srv)

	slog.Info("GRPC server listening on port :9000")
	return s.Serve(lis)
}

type ContextTestString struct{}

// middleware function to inject something into the context for later use
func testHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ContextTestString{}, "test string")

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
