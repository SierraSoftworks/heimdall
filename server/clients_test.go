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

func TestClients(t *testing.T) {
	Convey("Clients", t, func() {
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

		Convey("/api/v1/clients", func() {
			url := fmt.Sprintf("%s/api/v1/clients", s.URL)

			Convey("GET", func() {
				r, err := http.NewRequest("GET", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				data := []models.Client{}
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

				data := models.ClientDetails{}
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldNotBeNil)
				So(data.Client, ShouldNotBeNil)
				So(data.Name, ShouldEqual, "test-client")
			})

			Convey("DELETE", func() {
				r, err := http.NewRequest("DELETE", url, nil)
				So(err, ShouldBeNil)

				res, err := http.DefaultClient.Do(r)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(res.StatusCode, ShouldEqual, 200)

				data := models.ClientDetails{}
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldNotBeNil)
				So(data.Client, ShouldNotBeNil)
				So(data.Name, ShouldEqual, "test-client")

				cl, err := stores.GetStore().GetClient("test-client")
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

				data := []models.ClientDetails{}
				So(json.NewDecoder(res.Body).Decode(&data), ShouldBeNil)
				So(data, ShouldNotBeNil)
				So(data, ShouldHaveLength, 1)
			})
		})
	})
}
