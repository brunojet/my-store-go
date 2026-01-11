package observability

import (
	"sync"
)

// ObservabilityParams describes how to build observability middlewares
// for a selected HTTP driver.
type ObservabilityParams struct {
	// Driver is the HTTP driver (e.g. "gin", "chi").
	Driver string
	// Config enables/disables features.
	Config ObservabilityConfig
}

// ObservabilityManager owns the computed middleware set for a given driver+config.
//
// It caches the computed slice(s) since middleware construction can allocate.
type ObservabilityManager struct {
	Params ObservabilityParams

	mu          sync.Mutex
	middlewares *ObservabilityMiddlewares
}

func NewObservabilityManager(params ObservabilityParams) *ObservabilityManager {
	return &ObservabilityManager{Params: params}
}

func (m *ObservabilityManager) Open() (ObservabilityMiddlewares, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.middlewares != nil {
		return *m.middlewares, nil
	}

	mw := BuildMiddlewares(m.Params.Driver, m.Params.Config)
	m.middlewares = &mw
	return mw, nil
}

func (m *ObservabilityManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.middlewares = nil
	return nil
}
