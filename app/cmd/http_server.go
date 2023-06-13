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
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/fib"
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
	errorHandler := fib.NewErrorHandler(logger)

	f := fiber.New(fiber.Config{
		CaseSensitive:      true,
		EnableIPValidation: true,
		EnablePrintRoutes:  true, // This might be worth disabling
		ErrorHandler:       errorHandler.Handle,
		Immutable:          true, // This allows passing around the parsed payload nicely, but is slower
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
	s.useHelmetMiddleware()
	s.useZapLoggerMiddleware()
	s.usePProfMiddleware(cCtx)
	s.useLimiterMiddleware(cCtx)
	s.useCompressMiddleware(cCtx)
	s.useRecoverMiddleware()

	return s
}

// SetupHandlers sets up handlers for each module.
func (s *HTTPServer) SetupHandlers() *HTTPServer {
	// Add API handlers
	api := s.fiber.Group("/api")

	createRedirectHandler(s.logger).AddAPIRoutes(api)
	createTemplateHandler(s.logger).AddAPIRoutes(api)
	createPageHandler(s.logger).AddAPIRoutes(api)
	createBlockHandler(s.logger).AddAPIRoutes(api)
	createRendererHandler(s.logger).AddRoutes(api)

	api.Get("/healthz", func(cCtx *fiber.Ctx) error { panic(errors.New("hello")) })

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

func (s *HTTPServer) useLimiterMiddleware(cCtx *cli.Context) {
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

	s.fiber.Use(limiterMiddleware)
}

func (s *HTTPServer) useHelmetMiddleware() {
	s.fiber.Use(helmet.New())
}

func (s *HTTPServer) usePProfMiddleware(cCtx *cli.Context) {
	pprofEnabled := cCtx.Bool(usePprofFlag)
	pprofPrefix := cCtx.String(pprofPrefixFlag)

	if pprofEnabled {
		return
	}

	s.fiber.Use(pprof.New(pprof.Config{Prefix: pprofPrefix}))
}

func (s *HTTPServer) useCompressMiddleware(cCtx *cli.Context) {
	level := cCtx.Int(compressionLevelFlag)

	s.fiber.Use(compress.New(compress.Config{
		Level: compress.Level(level),
	}))
}

func (s *HTTPServer) useRecoverMiddleware() {
	s.fiber.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			s.logger.Error(fmt.Sprintf("%v", e), zap.String("stack", string(debug.Stack())))
		},
	}))
}

func (s *HTTPServer) useZapLoggerMiddleware() {
	s.fiber.Use(fiberzap.New(fiberzap.Config{
		Logger: s.logger,
	}))
}
