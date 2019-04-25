package managers

import (
	"github.com/SierraSoftworks/heimdall/internal/pkg/config"
	"github.com/SierraSoftworks/heimdall/internal/pkg/hub"
	"github.com/SierraSoftworks/heimdall/pkg/models"
)

func NewReportSubscriber(handler func(r *models.Report)) hub.Subscriber {
	return &reportSubscriber{
		handler: handler,
	}
}

type reportSubscriber struct {
	handler func(r *models.Report)
}

func (s *reportSubscriber) Receive(msg interface{}) {
	if report, ok := msg.(*models.Report); ok {
		go s.handler(report)
	}
}

func NewConfigSubscriber(handler func(c *config.Config)) hub.Subscriber {
	return &configSubscriber{
		handler: handler,
	}
}

type configSubscriber struct {
	handler func(r *config.Config)
}

func (s *configSubscriber) Receive(msg interface{}) {
	if report, ok := msg.(*config.Config); ok {
		go s.handler(report)
	}
}
