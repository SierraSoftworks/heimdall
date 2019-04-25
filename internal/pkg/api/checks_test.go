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

func TestChecks(t *testing.T) {
	Convey("Checks", t, func() {
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

		Convey("/api/v1/checks", func() {
			url := fmt.Sprintf("%s/api/v1/checks", s.URL)

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				var data []models.Check
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldHaveLength, 1)
			})
		})

		Convey("/api/v1/check/{check}", func() {
			url := fmt.Sprintf("%s/api/v1/check/%s", s.URL, "apache-port80")

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				var data models.CheckDetails
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data.Name, ShouldEqual, "apache-port80")
			})
		})

		Convey("/api/v1/check/{check}/clients", func() {
			url := fmt.Sprintf("%s/api/v1/check/%s/clients", s.URL, "apache-port80")

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				var data []models.Source
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldHaveLength, 1)
			})
		})

		Convey("/api/v1/check/{check}/client/{client}/executions", func() {
			url := fmt.Sprintf("%s/api/v1/check/%s/client/%s/executions", s.URL, "apache-port80", "test-client")

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				var data []models.Execution
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldHaveLength, 1)
			})
		})

		Convey("/api/v1/client/{client}/check/{check}/executions", func() {
			url := fmt.Sprintf("%s/api/v1/client/%s/check/%s/executions", s.URL, "test-client", "apache-port80")

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				var data []models.Execution
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldHaveLength, 1)
			})
		})
	})
}
