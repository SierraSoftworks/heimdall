package managers

import (
	"sync"

	"github.com/SierraSoftworks/heimdall/internal/pkg/api"
	"github.com/SierraSoftworks/heimdall/internal/pkg/config"
	"github.com/SierraSoftworks/heimdall/internal/pkg/hub"
)

func init() {
	RegisterManager(&APIManager{})
}

type APIManager struct {
	hub hub.Hub

	api *api.API
	l   sync.Mutex
}

func (m *APIManager) Setup(h hub.Hub) error {
	m.hub = h
	m.hub.Subscribe(NewConfigSubscriber(m.ConfigUpdated))

	return nil
}

func (m *APIManager) ConfigUpdated(conf *config.Config) {
	m.l.Lock()
	defer m.l.Unlock()

}

func (m *APIManager) Shutdown() {

}
