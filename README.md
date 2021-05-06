[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/zhangjie2012/promclient-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/zhangjie2012/promclient-go)](https://goreportcard.com/report/github.com/zhangjie2012/promclient-go)

# prometheus-client-go

a [Prometheus](https://prometheus.io/) go client (SDK), a simple Prometheus [HTTP API](https://prometheus.io/docs/prometheus/latest/querying/api/) wrapper.

## Progress

- [X] Query
- [X] QueryRange
- [X] QuerySeries
- [X] QueryLabels
- [X] QueryLabelValues
- [X] QueryTargets
- [ ] QueryRules
- [ ] QueryAlerts
- [ ] QueryTargetMetadata
- [ ] QueryMetricMetadata
- [ ] QueryAlertManagers
- [ ] QueryStatusConfig
- [ ] QueryStatusFlags
- [ ] QueryStatusRuntimeInfo
- [ ] QueryStatusBuildInfo
- [ ] QueryStatusTSDB
- [ ] TSDB Admin APIs

## Usages

```
go get github.com/zhangjie2012/promclient-go
```

```go
import promclient "github.com/zhangjie2012/promclient-go"

c = promclient.NewClient(PrometheusUrl, 0)
c.Query("up", float64(time.Now().Unix()))
```
