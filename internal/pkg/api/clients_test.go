package api

import (
	"fmt"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"
	"net/url"

	"encoding/json"

	"github.com/SierraSoftworks/heimdall/internal/pkg/plugins/memory"
	"github.com/SierraSoftworks/heimdall/pkg/driver"
	"github.com/SierraSoftworks/heimdall/pkg/duration"
	"github.com/SierraSoftworks/heimdall/pkg/models"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSources(t *testing.T) {
	Convey("Sources", t, func() {
		u, err := url.Parse("memory://memory")
		So(err, ShouldBeNil)

		store, err := plugins.GetStore(&driver.Driver{
			Type: "memory",
			URL:  u,
		})
		So(err, ShouldBeNil)
		So(store, ShouldNotBeNil)

		api, err := NewAPI(store)
		So(err, ShouldBeNil)

		s := httptest.NewServer(api.Router())
		defer s.Close()

		api.Store().(*memory.Memory).Reset()

		client := &models.Source{
			Name: "test-client",
			Tags: map[string]string{
				"type": "test",
			},
		}

		check := &models.Check{
			Name:    "apache-port80",
			Command: "curl -sS -D - http://localhost:80/",
			Collections: []string{
				"webservers",
			},
			Interval: duration.Duration(60 * time.Second),
			Timeout:  duration.Duration(5 * time.Second),
		}

		exec := &models.Execution{
			Scheduled: time.Now(),
			Executed:  time.Now(),
			Duration:  0,
			Status:    models.StatusOkay,
			Output:    "This is a quick test",
		}

		So(api.Store().AddReport(&models.Report{
			Source:    client,
			Check:     check,
			Execution: exec,
		}), ShouldBeNil)

		Convey("/api/v1/clients", func() {
			url := fmt.Sprintf("%s/api/v1/clients", s.URL)

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				data := []models.Source{}
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldNotBeNil)
				So(data, ShouldHaveLength, 1)
			})
		})

		Convey("/api/v1/client/{client}", func() {
			url := fmt.Sprintf("%s/api/v1/client/%s", s.URL, "test-client")

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				data := models.Source{}
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldNotBeNil)
				So(data.Name, ShouldEqual, "test-client")
			})

			Convey("DELETE", func() {
				r, err := http.NewRequest("DELETE", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				data := models.Source{}
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldNotBeNil)
				So(data.Name, ShouldEqual, "test-client")

				cl, err := api.Store().GetSource("test-client")
				So(err, ShouldBeNil)
				So(cl, ShouldBeNil)
			})
		})

		Convey("/api/v1/client/{client}/checks", func() {
			url := fmt.Sprintf("%s/api/v1/client/%s/checks", s.URL, "test-client")

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				data := []models.Source{}
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldNotBeNil)
				So(data, ShouldHaveLength, 1)
			})
		})
	})
}
