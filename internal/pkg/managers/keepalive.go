package managers

import (
	"fmt"
	"sync"
	"time"

	"github.com/SierraSoftworks/heimdall/internal/pkg/config"
	"github.com/SierraSoftworks/heimdall/internal/pkg/hub"
	"github.com/SierraSoftworks/heimdall/pkg/duration"
	"github.com/SierraSoftworks/heimdall/pkg/models"
	"github.com/SierraSoftworks/scheduler"
)

func init() {
	RegisterManager(&KeepaliveManager{})
}

type KeepaliveManager struct {
	hub hub.Hub

	schedule *scheduler.ActiveTask
	l        sync.Mutex
}

func (m *KeepaliveManager) Setup(h hub.Hub) error {
	m.hub = h

	m.hub.Subscribe(NewConfigSubscriber(m.ConfigUpdated))

	return nil
}

func (m *KeepaliveManager) ConfigUpdated(conf *config.Config) {
	m.l.Lock()
	defer m.l.Unlock()

	if m.schedule != nil {
		m.schedule.Cancel()
	}

	m.schedule = scheduler.Do(func(t time.Time) error {
		m.hub.Notify(&models.Report{
			Check: &models.Check{
				Name:     "keepalive",
				Command:  "",
				Interval: duration.Duration(30 * time.Second),
				Timeout:  0,
			},
			Source: conf.Source,
			Execution: &models.Execution{
				Scheduled: t,
				Executed:  t,
				Duration:  0,
				Status:    models.StatusOkay,
				Output:    m.generateOutput(conf),
			},
		})

		return nil
	}).Every(conf.Keepalive.Interval).Schedule()
}

func (m *KeepaliveManager) Shutdown() {
	m.l.Lock()
	defer m.l.Unlock()

	m.schedule.Cancel()
	m.schedule = nil
}

func (m *KeepaliveManager) generateOutput(conf *config.Config) string {
	return fmt.Sprintf(
		"OK - %d checks, %d collectors, %d publishers",
		len(conf.Checks),
		len(conf.Collectors),
		len(conf.Publishers),
	)
}
