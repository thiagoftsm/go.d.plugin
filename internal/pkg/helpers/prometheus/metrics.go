package prometheus

import (
	"math"
	"sort"

	"github.com/prometheus/prometheus/pkg/labels"
)

type (
	// Metric is a pair of label set and value
	Metric struct {
		Labels labels.Labels
		Value  float64
	}

	// Metrics is a list of Metric
	Metrics []Metric
)

// Name the __name__ label value
func (m Metric) Name() string {
	return m.Labels[0].Value
}

// Add Append a metric.
func (m *Metrics) Add(kv Metric) {
	*m = append(*m, kv)
}

// Reset Clear all data but reuse the memory alloced.
func (m *Metrics) Reset() {
	*m = (*m)[:0]
}

// Sort Sort
func (m Metrics) Sort() {
	sort.Sort(m)
}

func (m Metrics) Len() int {
	return len(m)
}

func (m Metrics) Less(i, j int) bool {
	return m[i].Name() < m[j].Name()
}

func (m Metrics) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// FindByName Find metrics where it's __name__ label matches given name.
// It expect the metrics is sorted.
// Complexity: O(log(N))
func (m Metrics) FindByName(name string) Metrics {
	from := sort.Search(len(m), func(i int) bool {
		return m[i].Name() >= name
	})
	if from == len(m) || m[from].Name() != name { // not found
		return Metrics{}
	}
	until := from + 1
	for until < len(m) && m[until].Name() == name {
		until++
	}
	return m[from:until]
}

// Match Find metrics where it's label matches given matcher.
// It do NOT expect the metrics is sorted.
// Complexity: O(N)
func (m Metrics) Match(matcher *labels.Matcher) Metrics {
	res := Metrics{}
	for _, kv := range m {
		value := kv.Labels.Get(matcher.Name)
		if matcher.Matches(value) {
			res.Add(kv)
		}
	}
	return res
}

// Max Return the max value.
// It do NOT expect the metrics is sorted.
// Complexity: O(N)
func (m Metrics) Max() float64 {
	max := -math.MaxFloat64
	for _, kv := range m {
		if max < kv.Value {
			max = kv.Value
		}
	}
	return max
}