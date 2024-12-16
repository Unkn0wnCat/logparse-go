package resultcollector

import (
	"fmt"
	"strings"
)

type Scope string

const (
	ScopeFile Scope = "file"
	ScopeLine Scope = "line"
)

type Result struct {
	Scope    Scope  `json:"scope"`
	Success  bool   `json:"success"`
	Filename string `json:"filename"`
	Line     int    `json:"line"`
	Message  string `json:"message"`
}

type ResultCollector struct {
	results []Result
}

func New() *ResultCollector {
	return &ResultCollector{
		results: []Result{},
	}
}

func (rc *ResultCollector) GetAll() []Result {
	return rc.results
}

func (rc *ResultCollector) DumpString() string {
	var builder strings.Builder

	for i, result := range rc.results {
		builder.WriteString(fmt.Sprintf("%05d - %s:%05d - Success: %t - Message: %s\n", i, result.Filename, result.Line, result.Success, result.Message))
	}

	return builder.String()
}

func (rc *ResultCollector) Add(result Result) {
	rc.results = append(rc.results, result)
}

func (rc *ResultCollector) Clear() {
	rc.results = []Result{}
}
