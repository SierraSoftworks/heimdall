package stores

import (
	"sort"
	"time"

	"github.com/SierraSoftworks/heimdall/models"
)

type Memory struct {
	clients    map[string]*memoryClient
	checks     map[string]*memoryCheck
	executions map[memoryExecutionKey][]models.Execution
	aggregates map[string]*memoryAggregate
}

type memoryExecutionKey struct {
	Client string
	Check  string
}

type memoryClient struct {
	Client *models.Client
	Checks map[string]*memoryExecState
}

type memoryCheck struct {
	Check   *models.Check
	Clients map[string]*memoryExecState
}

type memoryAggregate struct {
	Executions map[memoryExecutionKey]*memoryExecState
}

type memoryExecState struct {
	Status   models.Status
	Executed time.Time
}

func NewMemory() *Memory {
	return &Memory{
		clients:    map[string]*memoryClient{},
		checks:     map[string]*memoryCheck{},
		executions: map[memoryExecutionKey][]models.Execution{},
		aggregates: map[string]*memoryAggregate{},
	}
}

func (s *Memory) Reset() {
	s.clients = map[string]*memoryClient{}
	s.checks = map[string]*memoryCheck{}
	s.executions = map[memoryExecutionKey][]models.Execution{}
	s.aggregates = map[string]*memoryAggregate{}
}

func (s *Memory) AddReport(r *models.Report) error {
	e := &memoryExecState{
		Executed: r.Execution.Executed,
		Status:   r.Execution.Status,
	}

	s.updateClient(r, e)
	s.updateCheck(r, e)
	s.updateExecs(r, e)
	s.updateAggregates(r, e)

	return nil
}

func (s *Memory) updateClient(r *models.Report, e *memoryExecState) {
	c, ok := s.clients[r.Client.Name]
	if !ok {
		c = &memoryClient{
			Checks: map[string]*memoryExecState{},
		}
	}

	c.Client = r.Client
	c.Checks[r.Check.Name] = e

	s.clients[r.Client.Name] = c
}

func (s *Memory) updateCheck(r *models.Report, e *memoryExecState) {
	c, ok := s.checks[r.Check.Name]
	if !ok {
		c = &memoryCheck{
			Clients: map[string]*memoryExecState{},
		}
	}

	c.Check = r.Check
	c.Clients[r.Client.Name] = e

	s.checks[r.Check.Name] = c
}

func (s *Memory) updateExecs(r *models.Report, e *memoryExecState) {
	es, ok := s.executions[memoryExecutionKey{r.Client.Name, r.Check.Name}]
	if !ok {
		es = []models.Execution{*r.Execution}
	} else if len(es) == 10 {
		es = append(es[1:], *r.Execution)
	} else {
		es = append(es, *r.Execution)
	}

	sort.Sort(memoryOrderedCheckExecutions(es))

	s.executions[memoryExecutionKey{r.Client.Name, r.Check.Name}] = es
}

func (s *Memory) updateAggregates(r *models.Report, e *memoryExecState) {
	for _, an := range r.Check.Aggregates {
		a, ok := s.aggregates[an]
		if !ok {
			a = &memoryAggregate{
				Executions: map[memoryExecutionKey]*memoryExecState{},
			}
		}

		a.Executions[memoryExecutionKey{r.Client.Name, r.Check.Name}] = e

		s.aggregates[an] = a
	}
}

func (s *Memory) GetClients(q *ClientsQuery) ([]models.Client, error) {
	if q == nil {
		q = &ClientsQuery{}
	}

	cs := []models.Client{}
	for _, c := range s.clients {
		match := true
		for t, v := range q.Tags {
			if cv, ok := c.Client.Tags[t]; !ok || cv != v {
				match = false
				break
			}
		}

		if match {
			cs = append(cs, *c.Client)
		}
	}

	return cs, nil
}

func (s *Memory) GetClient(name string) (*models.ClientDetails, error) {
	c, ok := s.clients[name]
	if !ok {
		return nil, nil
	}

	t := time.Time{}
	for _, k := range c.Checks {
		if k.Executed.After(t) {
			t = k.Executed
		}
	}

	return &models.ClientDetails{
		Client:   c.Client,
		LastSeen: t,
	}, nil
}

func (s *Memory) GetClientChecks(client string) ([]models.CheckDetails, error) {
	c, ok := s.clients[client]
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

func (s *Memory) RemoveClient(name string) (*models.Client, error) {
	c, ok := s.clients[name]
	if !ok {
		return nil, nil
	}

	delete(s.clients, name)
	for cn := range c.Checks {
		ch := s.checks[cn]
		if _, ok := ch.Clients[name]; ok {
			delete(ch.Clients, name)
		}

		mk := memoryExecutionKey{name, cn}
		if _, ok := s.executions[mk]; ok {
			delete(s.executions, mk)
		}

		for _, an := range ch.Check.Aggregates {
			if _, ok := s.aggregates[an].Executions[mk]; ok {
				delete(s.aggregates[an].Executions, mk)
			}
		}
	}

	return c.Client, nil
}

func (s *Memory) GetChecks(q *ChecksQuery) ([]models.Check, error) {
	if q == nil {
		q = &ChecksQuery{}
	}

	hasStatusQuery := q.Status != nil && len(q.Status) > 0

	cs := []models.Check{}
	for _, c := range s.checks {
		match := true

		if hasStatusQuery {
			match = false
			for _, cl := range c.Clients {
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

	for _, cs := range c.Clients {
		if cs.Status.IsWorseThan(cd.Status) {
			cd.Status = cs.Status
		}

		if cs.Executed.After(cd.Executed) {
			cd.Executed = cs.Executed
		}
	}

	return cd, nil
}

func (s *Memory) GetCheckClients(check string) ([]models.ClientDetails, error) {
	c, ok := s.checks[check]
	if !ok {
		return nil, nil
	}

	cs := []models.ClientDetails{}
	for cn, cd := range c.Clients {
		if s.clients[cn] == nil {
			continue
		}

		cs = append(cs, models.ClientDetails{
			Client:   s.clients[cn].Client,
			Status:   cd.Status,
			LastSeen: cd.Executed,
		})
	}

	return cs, nil
}

func (s *Memory) GetCheckExecutions(client, check string) ([]models.Execution, error) {
	return s.executions[memoryExecutionKey{client, check}], nil
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
			ClientName: ck.Client,
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

func (s *Memory) GetAggregateClients(name string) ([]models.Client, error) {
	a, ok := s.aggregates[name]
	if !ok {
		return nil, nil
	}

	iterated := map[string]struct{}{}
	cs := []models.Client{}

	for ck := range a.Executions {
		_, ok := iterated[ck.Client]
		if ok {
			continue
		}

		iterated[ck.Client] = struct{}{}
		c, ok := s.clients[ck.Client]
		if !ok {
			continue
		}

		cs = append(cs, *c.Client)
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

type memoryOrderedCheckExecutions []models.Execution

func (m memoryOrderedCheckExecutions) Len() int {
	return len(m)
}

func (m memoryOrderedCheckExecutions) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m memoryOrderedCheckExecutions) Less(i, j int) bool {
	return m[i].Executed.Before(m[j].Executed)
}
