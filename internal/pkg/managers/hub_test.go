package managers

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/SierraSoftworks/heimdall/internal/pkg/hub"
	"github.com/SierraSoftworks/heimdall/pkg/models"
)

func TestHub(t *testing.T) {
	Convey("Hub", t, func(c C) {
		hub := hub.NewMemoryHub()
		So(hub, ShouldNotBeNil)

		rep := &models.Report{}

		handled := make(chan struct{})
		sub := NewReportSubscriber(func(r *models.Report) {
			handled <- struct{}{}
			c.So(r, ShouldNotBeNil)
			c.So(r, ShouldEqual, rep)
		})

		hub.Subscribe(sub)

		Convey("Notify", func() {
			hub.Notify(rep)
			select {
			case <-handled:
			case <-time.After(time.Second):
				So("not handled", ShouldBeNil)
			}
		})
	})
}
