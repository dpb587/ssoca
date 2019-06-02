package filter

import "fmt"

type defaultManager struct {
	filters map[string]Filter
}

func NewDefaultManager() defaultManager {
	return defaultManager{
		filters: map[string]Filter{},
	}
}

func (m *defaultManager) Add(name string, f Filter) {
	m.filters[name] = f
}

func (m defaultManager) Filters() []string {
	filters := []string{}

	for name, _ := range m.filters {
		filters = append(filters, name)
	}

	return filters
}

func (m *defaultManager) Get(name string) (Filter, error) {
	filter, ok := m.filters[name]
	if !ok {
		return nil, fmt.Errorf("unrecognized filter: %s", name)
	}

	return filter, nil
}
