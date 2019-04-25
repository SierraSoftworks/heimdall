package memory

import (
	"net/url"
	"testing"

	"time"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
	"github.com/SierraSoftworks/heimdall/pkg/duration"
	"github.com/SierraSoftworks/heimdall/pkg/models"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMemoryStore(t *testing.T) {
	Convey("Memory", t, func() {
		u, err := url.Parse("memory://memory")
		So(err, ShouldBeNil)

		s, err := NewMemoryStore(&driver.Driver{
			Type: "memory",
			URL:  u,
		})

		So(err, ShouldBeNil)
		So(s, ShouldNotBeNil)

		source := &models.Source{
			Name: "test-source",
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

		So(s.AddReport(&models.Report{
			Source:    source,
			Check:     check,
			Execution: exec,
		}), ShouldBeNil)

		Convey("GetSources", func() {
			Convey("Nil Query", func() {
				cs, err := s.GetSources(nil)
				So(err, ShouldBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Source{*source})
			})

			Convey("Matching Tags", func() {
				cs, err := s.GetSources(&plugins.SourcesQuery{
					Tags: map[string]string{"type": "test"},
				})
				So(err, ShouldBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Source{*source})
			})

			Convey("Unmatched Tags", func() {
				cs, err := s.GetSources(&plugins.SourcesQuery{
					Tags: map[string]string{"unknown": "tag"},
				})
				So(err, ShouldBeNil)
				So(cs, ShouldHaveLength, 0)
				So(cs, ShouldResemble, []models.Source{})
			})
		})

		Convey("GetSource", func() {
			Convey("Known Source", func() {
				c, err := s.GetSource("test-source")
				So(err, ShouldBeNil)
				So(c, ShouldNotBeNil)
				So(c, ShouldResemble, source)
			})

			Convey("Unknown Source", func() {
				c, err := s.GetSource("unknown-source")
				So(err, ShouldBeNil)
				So(c, ShouldBeNil)
			})
		})

		Convey("GetSourceChecks", func() {
			Convey("Known Source", func() {
				cs, err := s.GetSourceChecks("test-source")
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

			Convey("Unknnown Source", func() {
				cs, err := s.GetSourceChecks("unknown-source")
				So(err, ShouldBeNil)
				So(cs, ShouldBeNil)
			})
		})

		Convey("RemoveSource", func() {
			Convey("Known Source", func() {
				c, err := s.RemoveSource("test-source")
				So(err, ShouldBeNil)
				So(c, ShouldNotBeNil)
				So(c, ShouldResemble, source)

				cd, err := s.GetSource("test-source")
				So(err, ShouldBeNil)
				So(cd, ShouldBeNil)
			})

			Convey("Unknown Source", func() {
				c, err := s.RemoveSource("unknown-source")
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
				cs, err := s.GetChecks(&plugins.ChecksQuery{Status: []models.Status{}})
				So(err, ShouldBeNil)
				So(cs, ShouldNotBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Check{*check})
			})

			Convey("Matching Query", func() {
				cs, err := s.GetChecks(&plugins.ChecksQuery{Status: []models.Status{models.StatusOkay}})
				So(err, ShouldBeNil)
				So(cs, ShouldNotBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Check{*check})
			})

			Convey("Unmatched Query", func() {
				cs, err := s.GetChecks(&plugins.ChecksQuery{Status: []models.Status{models.StatusCrit}})
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

		Convey("GetCheckSources", func() {
			Convey("Known Check", func() {
				cs, err := s.GetCheckSources("apache-port80")
				So(err, ShouldBeNil)
				So(cs, ShouldNotBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Source{*source})
			})

			Convey("Unknown Check", func() {
				cs, err := s.GetCheckSources("unknown-check")
				So(err, ShouldBeNil)
				So(cs, ShouldBeNil)
			})
		})

		Convey("GetCheckExecutions", func() {
			Convey("Known Check", func() {
				cs, err := s.GetCheckExecutions("test-source", "apache-port80")
				So(err, ShouldBeNil)
				So(cs, ShouldNotBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Execution{
					*exec,
				})
			})

			Convey("Unknown Source", func() {
				cs, err := s.GetCheckExecutions("unknown-source", "apache-port80")
				So(err, ShouldBeNil)
				So(cs, ShouldBeNil)
			})

			Convey("Unknown Check", func() {
				cs, err := s.GetCheckExecutions("test-source", "unknown-check")
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
							ClientName: "test-source",
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

		Convey("GetAggregateSources", func() {
			Convey("Known Aggregate", func() {
				cs, err := s.GetAggregateSources("webservers")
				So(err, ShouldBeNil)
				So(cs, ShouldNotBeNil)
				So(cs, ShouldHaveLength, 1)
				So(cs, ShouldResemble, []models.Source{*source})
			})

			Convey("Unknown Aggregate", func() {
				cs, err := s.GetAggregateSources("unknown")
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
