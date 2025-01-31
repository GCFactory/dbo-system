package grpc

import (
	"context"
	"fmt"
	"github.com/GCFactory/dbo-system/platform/pkg/csrf"
	"github.com/GCFactory/dbo-system/service/file-api/docs"
	apiMiddlewares "github.com/GCFactory/dbo-system/service/file-api/internal/middleware"
	pb "github.com/GCFactory/dbo-system/service/file-api/proto/api/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"strings"
)

type ServiceServer struct {
	pb.UnimplementedFileApiServiceServer
}

func (s *Server) MapHandlers(e *echo.Echo) error {

	mw := apiMiddlewares.NewMiddlewareManager(s.cfg, []string{"*"}, s.logger)
	//
	e.Use(mw.RequestLoggerMiddleware)

	if s.cfg.Docs.Enable {
		docs.SwaggerInfo.Title = s.cfg.Docs.Title
		docs.SwaggerInfo.Version = s.cfg.Version
		e.GET(fmt.Sprintf("/%s/*", s.cfg.Docs.Prefix), echoSwagger.WrapHandler)
	}

	if s.cfg.HTTPServer.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())
	//e.Use(mw.MetricsMiddleware(metrics))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("5M"))
	if s.cfg.HTTPServer.Debug {
		e.Use(mw.DebugMiddleware)
	}

	pb.RegisterFileApiServiceServer(s.grpcServer, &ServiceServer{})

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", (&ServiceServer{}).httpEchoPing)

	//rootMux := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	//	fmt.Println("Called loger")
	//	s.logger.Debug("Called")
	//	s.logger.Debugf("Proto major is", req.ProtoMajor)
	//	if req.ProtoMajor == 2 && strings.Contains(req.Header.Get("Content-Type"), "application/grpc") {
	//		s.logger.Debug("Called grpc")
	//		s.grpcServer.ServeHTTP(w, req)
	//	} else {
	//		s.logger.Debug("Called http")
	//		httpMux.ServeHTTP(w, req)
	//	}
	//})

	//e.Server.Handler = rootMux
	//e.Server.Handler = h2c.NewHandler(rootMux, &http2.Server{})
	//e.TLSServer.Handler = h2c.NewHandler(rootMux, &http2.Server{})
	//e.Any("/", echo.WrapHandler(rootMux))
	//s.httpServer.Handler = h2c.NewHandler(rootMux, &http2.Server{})
	e.Any("/*", func(c echo.Context) error {
		h2c.NewHandler(
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.ProtoMajor == 2 && strings.Contains(req.Header.Get("Content-Type"), "application/grpc") {
					s.logger.Debug("Called grpc")
					s.grpcServer.ServeHTTP(w, req)
				} else {
					s.logger.Debug("Called http")
					httpMux.ServeHTTP(w, req)
				}
			}),
			&http2.Server{}).ServeHTTP(c.Response(), c.Request())
		return nil
	})

	return nil
}

func (s *ServiceServer) GetFile(ctx context.Context, in *pb.DeliveryGetFile) (*pb.File, error) {
	return nil, nil
}

// @Summary      Ping Service
// @Success      200  {string} string ""
// @Router       /api/v1/fileapi/ready/live [get]
func (s *ServiceServer) IsAlive(ctx context.Context, in *pb.IsAliveRequest) (*pb.IsAliveResponse, error) {
	fmt.Println("Called is alive")
	return &pb.IsAliveResponse{}, nil
}

func (s *ServiceServer) httpEchoPing(writer http.ResponseWriter, request *http.Request) {

}
