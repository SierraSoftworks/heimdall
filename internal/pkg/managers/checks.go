package managers

import (
	"sync"
	"time"

	"github.com/SierraSoftworks/heimdall/internal/pkg/config"
	"github.com/SierraSoftworks/heimdall/internal/pkg/hub"
	"github.com/SierraSoftworks/heimdall/internal/pkg/runner"
	"github.com/SierraSoftworks/heimdall/pkg/models"
	"github.com/SierraSoftworks/scheduler"
)

func init() {
	RegisterManager(&CheckManager{
		runner:          runner.NewCheckRunner(),
		scheduledChecks: []*scheduler.ActiveTask{},
	})
}

type CheckManager struct {
	hub hub.Hub

	runner          *runner.Runner
	scheduledChecks []*scheduler.ActiveTask
	l               sync.Mutex
}

func (m *CheckManager) Setup(h hub.Hub) error {
	m.hub = h

	m.hub.Subscribe(NewConfigSubscriber(m.ConfigUpdated))

	return nil
}

func (m *CheckManager) ConfigUpdated(conf *config.Config) {
	m.l.Lock()
	defer m.l.Unlock()

	for _, sched := range m.scheduledChecks {
		sched.Cancel()
	}

	m.scheduledChecks = make([]*scheduler.ActiveTask, 0)

	for _, check := range conf.Checks {
		func(check models.Check) {
			sched := scheduler.Do(func(t time.Time) error {
				m.hub.Notify(&models.Report{
					Execution: m.runner.ExecuteCheck(&check, conf.Source),
				})

				return nil
			}).Every(time.Duration(check.Interval)).Schedule()

			m.scheduledChecks = append(m.scheduledChecks, sched)
		}(check)
	}
}

func (m *CheckManager) Shutdown() {
	m.l.Lock()
	defer m.l.Unlock()

	for _, sched := range m.scheduledChecks {
		sched.Cancel()
	}

	m.scheduledChecks = make([]*scheduler.ActiveTask, 0)
}
