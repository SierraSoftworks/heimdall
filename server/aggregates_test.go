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

func TestAggregates(t *testing.T) {
	Convey("Aggregates", t, func() {
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

		Convey("/api/v1/aggregates", func() {
			url := fmt.Sprintf("%s/api/v1/aggregates", s.URL)

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				var data []models.Client
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldHaveLength, 1)
			})
		})

		Convey("/api/v1/aggregate/{aggregate}", func() {
			url := fmt.Sprintf("%s/api/v1/aggregate/%s", s.URL, "webservers")

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				var data models.AggregateDetails
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data.Name, ShouldEqual, "webservers")
			})

			Convey("DELETE", func() {
				r, err := http.NewRequest("DELETE", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				var data models.AggregateDetails
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data.Name, ShouldEqual, "webservers")
			})
		})

		Convey("/api/v1/aggregate/{aggregate}/checks", func() {
			url := fmt.Sprintf("%s/api/v1/aggregate/%s/checks", s.URL, "webservers")

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				var data []models.CheckDetails
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldHaveLength, 1)
			})
		})

		Convey("/api/v1/aggregate/{aggregate}/clients", func() {
			url := fmt.Sprintf("%s/api/v1/aggregate/%s/clients", s.URL, "webservers")

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
	})
}
