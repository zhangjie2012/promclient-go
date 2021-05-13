package promclient

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrorNoData = errors.New("Error: NoData")
	// ErrorNotCompatible: api return not compatible official document(maybe prometheus version upgrade ?)
	ErrorNotCompatible   = errors.New("Error: NotCompatible")
	ErrorUnsupportedType = errors.New("Error: UnsupportedType")
)

// PromResponse https://prometheus.io/docs/prometheus/latest/querying/api/#format-overview
type PromResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`

	// Only set if status is "error".
	// The data field may still hold additional data.
	ErrorType string `json:"errorType,omitempty"`
	Error     string `json:"error,omitempty"`

	// Only if there were warnings while executing the request.
	// There will still be data in the data field.
	Warning []string `json:"warnings"`
}

// TSValue value on some timestamp
type TSValue struct {
	Timestamp float64 `json:"timestamp"`
	Value     float64 `json:"value"`
}

func (t TSValue) String() string {
	return fmt.Sprintf("timestamp=%f value=%f", t.Timestamp, t.Value)
}

// QueryPoint one query result item
type QueryPoint struct {
	Metric map[string]string
	Values []TSValue
}

func (q QueryPoint) String() string {
	return fmt.Sprintf("metric=%v values=%v", q.Metric, q.Values)
}

type Series map[string]string

type Labels []string

type LabelValues []string

type TargetStateT string

const (
	TargetStateActive  TargetStateT = "active"
	TargetStateDropped TargetStateT = "dropped"
	TargetStateAny     TargetStateT = "any"
)

type ActiveTarget struct {
	DiscoveredLabels   map[string]interface{} `json:"discoveredLabels"`
	Labels             map[string]interface{} `json:"labels"`
	ScrapePool         string                 `json:"scrapePool"`
	ScrapeUrl          string                 `json:"scrapeUrl"`
	GlobalUrl          string                 `json:"globalUrl"`
	LastError          string                 `json:"lastError"`
	LastScrape         time.Time              `json:"lastScrape"`
	LastScrapeDuration float64                `json:"lastScrapeDuration"`
	Health             string                 `json:"health"`
}

type DroppedTarget struct {
	DiscoveredLabels map[string]interface{} `json:"discoveredLabels"`
}

type Target struct {
	ActiveTargets  []ActiveTarget  `json:"activeTargets"`
	DroppedTargets []DroppedTarget `json:"droppedTargets"`
}
