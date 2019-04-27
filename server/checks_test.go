package main

import (
	"fmt"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"

	"encoding/json"

	"github.com/SierraSoftworks/heimdall/models"
	"github.com/SierraSoftworks/heimdall/stores"
	. "github.com/smartystreets/goconvey/convey"
)

func TestChecks(t *testing.T) {
	Convey("Checks", t, func() {
		s := httptest.NewServer(router)
		defer s.Close()

		stores.GetStore().(*stores.Memory).Reset()

		client := &models.Client{
			Name: "test-client",
			Tags: map[string]string{
				"type": "test",
			},
		}

		check := &models.Check{
			Name:    "apache-port80",
			Command: "curl -sS -D - http://localhost:80/",
			Aggregates: []string{
				"webservers",
			},
			Interval: 60 * time.Second,
			Timeout:  5 * time.Second,
		}

		exec := &models.Execution{
			Scheduled: time.Now(),
			Executed:  time.Now(),
			Duration:  0,
			Status:    models.StatusOkay,
			Output:    "This is a quick test",
		}

		So(stores.GetStore().AddReport(&models.Report{
			Client:    client,
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

				var data []models.ClientDetails
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
