package utils

// ShutdownManager makes it simple to manage the lifecycles
// of a number of targets, specifically in environments
// where composition is commonplace.
type ShutdownManager struct {
	Targets []Shutdownable
}

// NewShutdownManager creates a new group which will be shutdown
// when the group itself is shutdown.
func NewShutdownManager() *ShutdownManager {
	return &ShutdownManager{
		Targets: []Shutdownable{},
	}
}

// AddTarget adds an entry to the list of items to be shutdown
// when this manager is told to shutdown.
func (m *ShutdownManager) AddTarget(target Shutdownable) {
	m.Targets = append(m.Targets, target)
}

// Shutdown instructs this manager to propagate the shutdown
// command to all of its registered targets.
func (m *ShutdownManager) Shutdown() {
	for _, target := range m.Targets {
		target.Shutdown()
	}
}

// Shutdownable represents something which can be instructed
// to shutdown any running tasks through a Shutdown() method.
type Shutdownable interface {
	Shutdown()
}
