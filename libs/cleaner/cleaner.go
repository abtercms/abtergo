package cleaner

import (
	"log/slog"
	"sync"
)

// Fn represents a callback function which can be used for cleaning up resources on server shutdown.
type Fn func() error

// Cleaner is a registry used to maintain and call callbacks for cleaning up resources on server shutdown.
type Cleaner struct {
	lock     sync.Mutex
	registry map[string]Fn
	logger   *slog.Logger
}

// New creates a new Cleaner instance.
func New(logger *slog.Logger) *Cleaner {
	return &Cleaner{
		registry: make(map[string]Fn),
		logger:   logger,
	}
}

// Register adds a new callback for cleanup registry on service shutdown.
func (c *Cleaner) Register(id string, fn Fn) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.registry[id] = fn
}

// Unregister removes a callback from being used for cleanup on service shutdown.
func (c *Cleaner) Unregister(id string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.registry, id)
}

// Run runs all registered cleanup callbacks present in the registry.
func (c *Cleaner) Run() {
	c.lock.Lock()
	defer c.lock.Unlock()

	var wg sync.WaitGroup
	for id, fn := range c.registry {
		wg.Add(1)

		id := id
		fn := fn

		go func() {
			defer wg.Done()
			err := fn()
			if err != nil {
				c.logger.Error("clean up failed", slog.Attr{Key: "err", Value: slog.StringValue(err.Error())}, slog.Attr{Key: "id", Value: slog.StringValue(id)})
			}
		}()
	}

	wg.Wait()
}
