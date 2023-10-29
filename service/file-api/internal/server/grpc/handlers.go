package grpc

import (
	"context"
	pb "github.com/GCFactory/dbo-system/service/file-api/proto/api/v1"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"strings"
)

type ServiceServer struct {
	pb.UnimplementedFileApiServiceServer
}

func (s *Server) MapHandlers() error {
	pb.RegisterFileApiServiceServer(s.grpcServer, &ServiceServer{})

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", (&ServiceServer{}).httpEchoPing)

	rootMux := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.ProtoMajor == 2 && strings.Contains(req.Header.Get("Content-Type"), "application/grpc") {
			s.grpcServer.ServeHTTP(w, req)
		} else {
			httpMux.ServeHTTP(w, req)
		}
	})

	s.httpServer.Handler = h2c.NewHandler(rootMux, &http2.Server{})

	return nil
}

func (s *ServiceServer) GetFile(ctx context.Context, in *pb.DeliveryGetFile) (*pb.File, error) {
	return nil, nil
}

func (s *ServiceServer) IsAlive(ctx context.Context, in *pb.IsAliveRequest) (*pb.IsAliveResponse, error) {
	return nil, nil
}

func (s *ServiceServer) httpEchoPing(writer http.ResponseWriter, request *http.Request) {

}
