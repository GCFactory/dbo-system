package server

import (
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "net/http/pprof"
	"os"
)

const (
	certFile       = "ssl/Server.crt"
	keyFile        = "ssl/Server.pem"
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

// Server struct
type Server struct {
	echo   *echo.Echo
	cfg    *config.Config
	db     *sqlx.DB
	logger logger.Logger
	// Channel to control goroutines
	shutdownChan chan os.Signal
	httpChan     chan int
}

func NewServer(cfg *config.Config, db *sqlx.DB, logger logger.Logger) *Server {
	return &Server{echo: echo.New(), cfg: cfg, db: db, logger: logger, shutdownChan: make(chan os.Signal, 3), httpChan: make(chan int, 3)}
}

func (s *Server) Run() error {

}
