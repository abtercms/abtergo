package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/fib"
	"github.com/abtergo/abtergo/pkg/block"
	"github.com/abtergo/abtergo/pkg/page"
	"github.com/abtergo/abtergo/pkg/redirect"
	"github.com/abtergo/abtergo/pkg/renderer"
	"github.com/abtergo/abtergo/pkg/template"
)

type cleaner interface {
	Run()
}

// HTTPServer represents an HTTP server and contains all dependencies necessary.
type HTTPServer struct {
	logger  *zap.Logger
	cleaner cleaner
	fiber   *fiber.App
}

// NewHTTPServer creates a new HTTPServer instance.
func NewHTTPServer(logger *zap.Logger, cleaner cleaner) *HTTPServer {
	f := fiber.New(fiber.Config{
		CaseSensitive:      true,
		EnableIPValidation: true,
		EnablePrintRoutes:  true,             // This might be worth disabling
		ErrorHandler:       fib.ErrorHandler, // TODO: Make this more useful
		Immutable:          true,             // This allows passing around the parsed payload nicely, but is slower
		JSONDecoder:        json.Unmarshal,
		JSONEncoder:        json.Marshal,
	})

	return &HTTPServer{
		logger:  logger,
		cleaner: cleaner,
		fiber:   f,
	}
}

// SetupMiddleware sets up middleware to be used.
func (s *HTTPServer) SetupMiddleware(cCtx *cli.Context) *HTTPServer {
	// Add middleware
	s.useLoggerMiddleware(s.fiber)
	s.usePProfMiddleware(cCtx, s.fiber)
	s.useLimiterMiddleware(cCtx, s.fiber)
	s.useCompressMiddleware(cCtx, s.fiber)
	s.useRecoverMiddleware(s.fiber)

	return s
}

// SetupHandlers sets up handlers for each module.
func (s *HTTPServer) SetupHandlers() *HTTPServer {
	// Add API handlers
	api := s.fiber.Group("/api")
	s.addRedirectRoutes(api)
	s.addTemplateRoutes(api)
	s.addPageRoutes(api)
	s.addBlockRoutes(api)

	s.addCatchAll()

	return s
}

// Start starts a new HTTPServer.
func (s *HTTPServer) Start(cCtx *cli.Context) error {
	// Listen from a different goroutine
	go func() {
		if err := s.fiber.Listen(fmt.Sprintf(":%d", cCtx.Uint("port"))); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c
	s.logger.Info("server shutting down gracefully...")
	err := s.fiber.Shutdown()
	if err != nil {
		s.logger.Error(errors.Wrap(err, "shutting down the HTTP server failed").Error())

		return err
	}

	s.logger.Info("running cleanup tasks...")

	s.cleaner.Run()

	s.logger.Info("fiber was successful shutdown.")

	return nil
}

const (
	limiterMaxRequestsFlag = "max-requests"
	limiterTimeframeFlag   = "timeframe"
	usePprofFlag           = "pprof"
	pprofPrefixFlag        = "pprof-prefix"
	compressionLevelFlag   = "compress-level"
)

const (
	localhost = "127.0.0.1"
)

func (s *HTTPServer) useLimiterMiddleware(cCtx *cli.Context, router fiber.Router) {
	if !cCtx.Bool("use-pprof") {
		return
	}

	max := cCtx.Int(limiterMaxRequestsFlag)
	expiration := time.Duration(cCtx.Int(limiterTimeframeFlag)) * time.Second

	limiterMiddleware := limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == localhost
		},
		Max:        max,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			// TODO: Make this based on auth for services
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	})

	router.Use(limiterMiddleware)
}

func (s *HTTPServer) useLoggerMiddleware(r fiber.Router) {
	r.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
}

func (s *HTTPServer) usePProfMiddleware(cCtx *cli.Context, r fiber.Router) {
	pprofEnabled := cCtx.Bool(usePprofFlag)
	pprofPrefix := cCtx.String(pprofPrefixFlag)

	if pprofEnabled {
		return
	}

	r.Use(pprof.New(pprof.Config{Prefix: pprofPrefix}))
}

func (s *HTTPServer) useCompressMiddleware(cCtx *cli.Context, r fiber.Router) {
	level := cCtx.Int(compressionLevelFlag)

	r.Use(compress.New(compress.Config{
		Level: compress.Level(level),
	}))
}

func (s *HTTPServer) useRecoverMiddleware(r fiber.Router) {
	r.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			s.logger.Error(fmt.Sprintf("%v", e), zap.String("stack", string(debug.Stack())))
		},
	}))
}

func (s *HTTPServer) addRedirectRoutes(r fiber.Router) {
	repo := redirect.NewInMemoryRepo()
	service := redirect.NewService(s.logger, repo)
	handler := redirect.NewHandler(s.logger, service)

	handler.AddAPIRoutes(r)
}

func (s *HTTPServer) addTemplateRoutes(r fiber.Router) {
	repo := template.NewInMemoryRepo()
	service := template.NewService(s.logger, repo)
	handler := template.NewHandler(s.logger, service)

	handler.AddAPIRoutes(r)
}

func (s *HTTPServer) addPageRoutes(r fiber.Router) {
	repo := page.NewInMemoryRepo()
	updater := page.NewUpdater()
	service := page.NewService(s.logger, repo, updater)
	handler := page.NewHandler(s.logger, service)

	handler.AddAPIRoutes(r)
}

func (s *HTTPServer) addBlockRoutes(r fiber.Router) {
	repo := block.NewInMemoryRepo()
	service := block.NewService(s.logger, repo)
	handler := block.NewHandler(s.logger, service)

	handler.AddAPIRoutes(r)
}

func (s *HTTPServer) addCatchAll() {
	pageRepo := page.NewInMemoryRepo()
	templateRepo := template.NewInMemoryRepo()
	blockRepo := block.NewInMemoryRepo()
	redirectRepo := redirect.NewInMemoryRepo()

	r := renderer.NewRenderer()
	service := renderer.NewService(r, pageRepo, templateRepo, blockRepo, redirectRepo)
	handler := renderer.NewHandler(s.logger, service)

	handler.AddRoutes(s.fiber)
}
