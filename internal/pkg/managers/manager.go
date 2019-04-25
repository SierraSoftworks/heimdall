package managers

import (
	"github.com/SierraSoftworks/heimdall/internal/pkg/hub"
)

var managers = []Manager{}

func GetManagers() []Manager {
	return managers
}

func RegisterManager(manager Manager) {
	managers = append(managers, manager)
}

type Manager interface {
	Setup(h hub.Hub) error
	Shutdown()
}
