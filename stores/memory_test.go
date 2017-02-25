package stores

import (
	"testing"

	"time"

	"github.com/SierraSoftworks/heimdall/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMemoryStore(t *testing.T) {
	Convey("Memory", t, func() {
		s := NewMemory()

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

		So(s.AddReport(&models.Report{
			Client:    client,
			Check:     check,
			Execution: exec,
		}), ShouldBeNil)

		Convey("GetClients", func() {
			Convey("Nil Query", func() {
				cs, err := s.GetClients(nil)
				So(err, ShouldBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Client{*client})
			})

			Convey("Matching Tags", func() {
				cs, err := s.GetClients(&ClientsQuery{
					Tags: map[string]string{"type": "test"},
				})
				So(err, ShouldBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Client{*client})
			})

			Convey("Unmatched Tags", func() {
				cs, err := s.GetClients(&ClientsQuery{
					Tags: map[string]string{"unknown": "tag"},
				})
				So(err, ShouldBeNil)
				So(cs, ShouldHaveLength, 0)
				So(cs, ShouldResemble, []models.Client{})
			})
		})

		Convey("GetClient", func() {
			Convey("Known Client", func() {
				c, err := s.GetClient("test-client")
				So(err, ShouldBeNil)
				So(c, ShouldNotBeNil)
				So(c, ShouldResemble, &models.ClientDetails{
					Client:   client,
					LastSeen: exec.Executed,
					Status:   models.StatusOkay,
				})
			})

			Convey("Unknown Client", func() {
				c, err := s.GetClient("unknown-client")
				So(err, ShouldBeNil)
				So(c, ShouldBeNil)
			})
		})

		Convey("GetClientChecks", func() {
			Convey("Known Client", func() {
				cs, err := s.GetClientChecks("test-client")
				So(err, ShouldBeNil)
				So(cs, ShouldNotBeNil)
				So(cs, ShouldResemble, []models.CheckDetails{
					models.CheckDetails{
						Check:    check,
						Status:   models.StatusOkay,
						Executed: exec.Executed,
					},
				})
			})

			Convey("Unknnown Client", func() {
				cs, err := s.GetClientChecks("unknown-client")
				So(err, ShouldBeNil)
				So(cs, ShouldBeNil)
			})
		})

		Convey("RemoveClient", func() {
			Convey("Known Client", func() {
				c, err := s.RemoveClient("test-client")
				So(err, ShouldBeNil)
				So(c, ShouldNotBeNil)
				So(c, ShouldResemble, client)

				cd, err := s.GetClient("test-client")
				So(err, ShouldBeNil)
				So(cd, ShouldBeNil)
			})

			Convey("Unknown Client", func() {
				c, err := s.RemoveClient("unknown-client")
				So(err, ShouldBeNil)
				So(c, ShouldBeNil)
			})
		})

		Convey("GetChecks", func() {
			Convey("Nil Query", func() {
				cs, err := s.GetChecks(nil)
				So(err, ShouldBeNil)
				So(cs, ShouldNotBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Check{*check})
			})

			Convey("Empty Query", func() {
				cs, err := s.GetChecks(&ChecksQuery{Status: []models.Status{}})
				So(err, ShouldBeNil)
				So(cs, ShouldNotBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Check{*check})
			})

			Convey("Matching Query", func() {
				cs, err := s.GetChecks(&ChecksQuery{Status: []models.Status{models.StatusOkay}})
				So(err, ShouldBeNil)
				So(cs, ShouldNotBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Check{*check})
			})

			Convey("Unmatched Query", func() {
				cs, err := s.GetChecks(&ChecksQuery{Status: []models.Status{models.StatusCrit}})
				So(err, ShouldBeNil)
				So(cs, ShouldNotBeNil)
				So(cs, ShouldHaveLength, 0)
				So(cs, ShouldResemble, []models.Check{})
			})
		})

		Convey("GetCheck", func() {
			Convey("Known Check", func() {
				c, err := s.GetCheck("apache-port80")
				So(err, ShouldBeNil)
				So(c, ShouldNotBeNil)
				So(c, ShouldResemble, &models.CheckDetails{
					Check:    check,
					Status:   exec.Status,
					Executed: exec.Executed,
				})
			})

			Convey("Unknown Check", func() {
				c, err := s.GetCheck("unknown-check")
				So(err, ShouldBeNil)
				So(c, ShouldBeNil)
			})
		})

		Convey("GetCheckClients", func() {
			Convey("Known Check", func() {
				cs, err := s.GetCheckClients("apache-port80")
				So(err, ShouldBeNil)
				So(cs, ShouldNotBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.ClientDetails{
					models.ClientDetails{
						Client:   client,
						Status:   exec.Status,
						LastSeen: exec.Executed,
					},
				})
			})

			Convey("Unknown Check", func() {
				cs, err := s.GetCheckClients("unknown-check")
				So(err, ShouldBeNil)
				So(cs, ShouldBeNil)
			})
		})

		Convey("GetCheckExecutions", func() {
			Convey("Known Check", func() {
				cs, err := s.GetCheckExecutions("test-client", "apache-port80")
				So(err, ShouldBeNil)
				So(cs, ShouldNotBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Execution{
					*exec,
				})
			})

			Convey("Unknown Client", func() {
				cs, err := s.GetCheckExecutions("unknown-client", "apache-port80")
				So(err, ShouldBeNil)
				So(cs, ShouldBeNil)
			})

			Convey("Unknown Check", func() {
				cs, err := s.GetCheckExecutions("test-client", "unknown-check")
				So(err, ShouldBeNil)
				So(cs, ShouldBeNil)
			})
		})

		Convey("GetAggregates", func() {
			as, err := s.GetAggregates()
			So(err, ShouldBeNil)
			So(as, ShouldNotBeNil)
			So(as, ShouldHaveLength, 1)
			So(as, ShouldResemble, []models.Aggregate{
				models.Aggregate{
					Name:   "webservers",
					Status: exec.Status,
				},
			})
		})

		Convey("GetAggregate", func() {
			Convey("Known Aggregate", func() {
				a, err := s.GetAggregate("webservers")
				So(err, ShouldBeNil)
				So(a, ShouldNotBeNil)
				So(a, ShouldResemble, &models.AggregateDetails{
					Aggregate: &models.Aggregate{
						Name:   "webservers",
						Status: exec.Status,
					},
					Entries: []models.AggregateEntry{
						models.AggregateEntry{
							CheckName:  "apache-port80",
							ClientName: "test-client",
							Status:     exec.Status,
							Executed:   exec.Executed,
						},
					},
				})
			})

			Convey("Unknown Aggregate", func() {
				a, err := s.GetAggregate("unknown")
				So(err, ShouldBeNil)
				So(a, ShouldBeNil)
			})
		})

		Convey("GetAggregateChecks", func() {
			Convey("Known Aggregate", func() {
				cs, err := s.GetAggregateChecks("webservers")
				So(err, ShouldBeNil)
				So(cs, ShouldNotBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Check{*check})
			})

			Convey("Unknown Aggregate", func() {
				cs, err := s.GetAggregateChecks("unknown")
				So(err, ShouldBeNil)
				So(cs, ShouldBeNil)
			})
		})

		Convey("GetAggregateClients", func() {
			Convey("Known Aggregate", func() {
				cs, err := s.GetAggregateClients("webservers")
				So(err, ShouldBeNil)
				So(cs, ShouldNotBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Client{*client})
			})

			Convey("Unknown Aggregate", func() {
				cs, err := s.GetAggregateClients("unknown")
				So(err, ShouldBeNil)
				So(cs, ShouldBeNil)
			})
		})

		Convey("RemoveAggregate", func() {
			Convey("Known Aggregate", func() {
				a, err := s.RemoveAggregate("webservers")
				So(err, ShouldBeNil)
				So(a, ShouldNotBeNil)
				So(a, ShouldResemble, &models.Aggregate{
					Name:   "webservers",
					Status: exec.Status,
				})

				ad, err := s.GetAggregate("webservers")
				So(err, ShouldBeNil)
				So(ad, ShouldBeNil)
			})

			Convey("Unknown Aggregate", func() {
				a, err := s.RemoveAggregate("unknown")
				So(err, ShouldBeNil)
				So(a, ShouldBeNil)
			})
		})
	})
}
