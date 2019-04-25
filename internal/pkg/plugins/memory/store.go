package memory

import (
	"time"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
	"github.com/SierraSoftworks/heimdall/pkg/models"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"
)

type Memory struct {
	conf *driver.Driver

	sources    map[string]*memorySource
	checks     map[string]*memoryCheck
	executions map[memoryExecutionKey][]models.Execution
	aggregates map[string]*memoryAggregate
}

type memoryExecutionKey struct {
	Source string
	Check  string
}

type memorySource struct {
	Source *models.Source
	Checks map[string]*memoryExecState
}

type memoryCheck struct {
	Check   *models.Check
	Sources map[string]*memoryExecState
}

type memoryAggregate struct {
	Executions map[memoryExecutionKey]*memoryExecState
}

type memoryExecState struct {
	Status   models.Status
	Executed time.Time
}

func NewMemoryStore(cfg *driver.Driver) (plugins.Store, error) {
	return &Memory{
		conf: cfg,

		sources:    map[string]*memorySource{},
		checks:     map[string]*memoryCheck{},
		executions: map[memoryExecutionKey][]models.Execution{},
		aggregates: map[string]*memoryAggregate{},
	}, nil
}

func (c *Memory) Driver() *driver.Driver {
	return c.conf
}

func (s *Memory) Reset() {
	s.sources = map[string]*memorySource{}
	s.checks = map[string]*memoryCheck{}
	s.executions = map[memoryExecutionKey][]models.Execution{}
	s.aggregates = map[string]*memoryAggregate{}
}

func (s *Memory) AddReport(r *models.Report) error {
	e := &memoryExecState{
		Executed: r.Execution.Executed,
		Status:   r.Execution.Status,
	}

	s.updateSource(r, e)
	s.updateCheck(r, e)
	s.updateExecs(r, e)
	s.updateAggregates(r, e)

	return nil
}

func (s *Memory) updateSource(r *models.Report, e *memoryExecState) {
	c, ok := s.sources[r.Source.Name]
	if !ok {
		c = &memorySource{
			Checks: map[string]*memoryExecState{},
		}
	}

	c.Source = r.Source
	c.Checks[r.Check.Name] = e

	s.sources[r.Source.Name] = c
}

func (s *Memory) updateCheck(r *models.Report, e *memoryExecState) {
	c, ok := s.checks[r.Check.Name]
	if !ok {
		c = &memoryCheck{
			Sources: map[string]*memoryExecState{},
		}
	}

	c.Check = r.Check
	c.Sources[r.Source.Name] = e

	s.checks[r.Check.Name] = c
}

func (s *Memory) updateExecs(r *models.Report, e *memoryExecState) {
	es, ok := s.executions[memoryExecutionKey{r.Source.Name, r.Check.Name}]
	if !ok {
		es = []models.Execution{*r.Execution}
	} else if len(es) == 10 {
		es = append(es[1:], *r.Execution)
	} else {
		es = append(es, *r.Execution)
	}
	s.executions[memoryExecutionKey{r.Source.Name, r.Check.Name}] = es
}

func (s *Memory) updateAggregates(r *models.Report, e *memoryExecState) {
	for _, an := range r.Check.Collections {
		a, ok := s.aggregates[an]
		if !ok {
			a = &memoryAggregate{
				Executions: map[memoryExecutionKey]*memoryExecState{},
			}
		}

		a.Executions[memoryExecutionKey{r.Source.Name, r.Check.Name}] = e

		s.aggregates[an] = a
	}
}

func (s *Memory) GetSources(q *plugins.SourcesQuery) ([]models.Source, error) {
	if q == nil {
		q = &plugins.SourcesQuery{}
	}

	cs := []models.Source{}
	for _, c := range s.sources {
		match := true
		for t, v := range q.Tags {
			if cv, ok := c.Source.Tags[t]; !ok || cv != v {
				match = false
				break
			}
		}

		if match {
			cs = append(cs, *c.Source)
		}
	}

	return cs, nil
}

func (s *Memory) GetSource(name string) (*models.Source, error) {
	c, ok := s.sources[name]
	if !ok {
		return nil, nil
	}

	t := time.Time{}
	for _, k := range c.Checks {
		if k.Executed.After(t) {
			t = k.Executed
		}
	}

	return c.Source, nil
}

func (s *Memory) GetSourceChecks(source string) ([]models.CheckDetails, error) {
	c, ok := s.sources[source]
	if !ok {
		return nil, nil
	}

	cs := []models.CheckDetails{}
	for cn, cd := range c.Checks {
		if s.checks[cn] == nil {
			continue
		}

		cs = append(cs, models.CheckDetails{
			Check:    s.checks[cn].Check,
			Status:   cd.Status,
			Executed: cd.Executed,
		})
	}

	return cs, nil
}

func (s *Memory) RemoveSource(name string) (*models.Source, error) {
	c, ok := s.sources[name]
	if !ok {
		return nil, nil
	}

	delete(s.sources, name)
	for cn := range c.Checks {
		ch := s.checks[cn]
		if _, ok := ch.Sources[name]; ok {
			delete(ch.Sources, name)
		}

		mk := memoryExecutionKey{name, cn}
		if _, ok := s.executions[mk]; ok {
			delete(s.executions, mk)
		}

		for _, an := range ch.Check.Collections {
			if _, ok := s.aggregates[an].Executions[mk]; ok {
				delete(s.aggregates[an].Executions, mk)
			}
		}
	}

	return c.Source, nil
}

func (s *Memory) GetChecks(q *plugins.ChecksQuery) ([]models.Check, error) {
	if q == nil {
		q = &plugins.ChecksQuery{}
	}

	hasStatusQuery := q.Status != nil && len(q.Status) > 0

	cs := []models.Check{}
	for _, c := range s.checks {
		match := true

		if hasStatusQuery {
			match = false
			for _, cl := range c.Sources {
				for _, qs := range q.Status {
					if cl.Status == qs {
						match = true
						break
					}
				}

				if match {
					break
				}
			}
		}

		if match {
			cs = append(cs, *c.Check)
		}
	}

	return cs, nil
}

func (s *Memory) GetCheck(name string) (*models.CheckDetails, error) {
	c, ok := s.checks[name]
	if !ok {
		return nil, nil
	}

	cd := &models.CheckDetails{
		Check: c.Check,
	}

	for _, cs := range c.Sources {
		if cs.Status.IsWorseThan(cd.Status) {
			cd.Status = cs.Status
		}

		if cs.Executed.After(cd.Executed) {
			cd.Executed = cs.Executed
		}
	}

	return cd, nil
}

func (s *Memory) GetCheckSources(check string) ([]models.Source, error) {
	c, ok := s.checks[check]
	if !ok {
		return nil, nil
	}

	cs := []models.Source{}
	for cn := range c.Sources {
		if s.sources[cn] == nil {
			continue
		}

		cs = append(cs, *s.sources[cn].Source)
	}

	return cs, nil
}

func (s *Memory) GetCheckExecutions(source, check string) ([]models.Execution, error) {
	return s.executions[memoryExecutionKey{source, check}], nil
}

func (s *Memory) GetAggregates() ([]models.Aggregate, error) {
	as := []models.Aggregate{}

	for an, ad := range s.aggregates {
		a := models.Aggregate{
			Name:   an,
			Status: models.StatusOkay,
		}

		for _, cd := range ad.Executions {
			if cd.Status.IsWorseThan(a.Status) {
				a.Status = cd.Status
			}
		}

		as = append(as, a)
	}

	return as, nil
}

func (s *Memory) GetAggregate(name string) (*models.AggregateDetails, error) {
	a, ok := s.aggregates[name]
	if !ok {
		return nil, nil
	}

	ad := &models.AggregateDetails{
		Aggregate: &models.Aggregate{
			Name:   name,
			Status: models.StatusOkay,
		},
		Entries: []models.AggregateEntry{},
	}

	for ck, cd := range a.Executions {
		if cd.Status.IsWorseThan(ad.Status) {
			ad.Status = cd.Status
		}

		ad.Entries = append(ad.Entries, models.AggregateEntry{
			ClientName: ck.Source,
			CheckName:  ck.Check,
			Status:     cd.Status,
			Executed:   cd.Executed,
		})
	}

	return ad, nil
}

func (s *Memory) GetAggregateChecks(name string) ([]models.Check, error) {
	a, ok := s.aggregates[name]
	if !ok {
		return nil, nil
	}

	iterated := map[string]struct{}{}
	cs := []models.Check{}

	for ck := range a.Executions {
		_, ok := iterated[ck.Check]
		if ok {
			continue
		}

		iterated[ck.Check] = struct{}{}
		c, ok := s.checks[ck.Check]
		if !ok {
			continue
		}

		cs = append(cs, *c.Check)
	}

	return cs, nil
}

func (s *Memory) GetAggregateSources(name string) ([]models.Source, error) {
	a, ok := s.aggregates[name]
	if !ok {
		return nil, nil
	}

	iterated := map[string]struct{}{}
	cs := []models.Source{}

	for ck := range a.Executions {
		_, ok := iterated[ck.Source]
		if ok {
			continue
		}

		iterated[ck.Source] = struct{}{}
		c, ok := s.sources[ck.Source]
		if !ok {
			continue
		}

		cs = append(cs, *c.Source)
	}

	return cs, nil
}

func (s *Memory) RemoveAggregate(name string) (*models.Aggregate, error) {
	a, ok := s.aggregates[name]
	if !ok {
		return nil, nil
	}

	ad := &models.Aggregate{
		Name:   name,
		Status: models.StatusOkay,
	}

	for _, cd := range a.Executions {
		if cd.Status.IsWorseThan(ad.Status) {
			ad.Status = cd.Status
		}
	}

	delete(s.aggregates, name)
	return ad, nil
}
