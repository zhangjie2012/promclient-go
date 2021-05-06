package promclient

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
)

// HttpCodeTrans translate request status code to human-readable
// https://prometheus.io/docs/prometheus/latest/querying/api/#format-overview
func HttpCodeTrans(status int) error {
	switch status {
	case 200:
		return nil
	case 400:
		return fmt.Errorf("400 Bad Request")
	case 404:
		return fmt.Errorf("404 Not Found")
	case 422:
		return fmt.Errorf("422 Unprocessable Entity")
	case 503:
		return fmt.Errorf("503 Service Unavailable")
	default:
		return fmt.Errorf("%d %s", status, http.StatusText(status))
	}
}

// ParseQueryResult: parse prometheus query API return data
// "/api/v1/query" and "/api/v1/query_range"
func ParseQueryResult(raw interface{}) ([]QueryPoint, error) {
	m := raw.(map[string]interface{})

	// must have "resultType" and "result" fields
	resultType, ok1 := m["resultType"]
	result, ok2 := m["result"]
	if !(ok1 && ok2) {
		return nil, ErrorNotCompatible
	}

	switch resultType {
	case "vector":
		return parseVectorResult(result)
	case "matrix":
		return parseMatrixResult(result)
	default:
		return nil, ErrorUnsupportedType
	}
}

// parseVectorResult parse resultType="vector"
func parseVectorResult(raw interface{}) ([]QueryPoint, error) {
	points := []QueryPoint{}

	results := raw.([]interface{})
	for _, result := range results {
		data := result.(map[string]interface{})

		rMetric, ok1 := data["metric"]
		rValue, ok2 := data["value"]
		if !(ok1 && ok2) {
			return nil, ErrorNotCompatible
		}

		rValues := rValue.([]interface{})

		metric := ToMetric(rMetric.(map[string]interface{}))
		timestamp := rValues[0].(float64)
		value := ToValue(rValues[1].(string))

		points = append(points, QueryPoint{
			Metric: metric,
			Values: []TSValue{{Timestamp: timestamp, Value: value}},
		})
	}

	return points, nil
}

// parseMatrixResult parse resultType="matrix"
func parseMatrixResult(raw interface{}) ([]QueryPoint, error) {
	points := []QueryPoint{}

	results := raw.([]interface{})
	for _, result := range results {
		data := result.(map[string]interface{})
		rMetric, ok1 := data["metric"]
		rValue, ok2 := data["values"]
		if !(ok1 && ok2) {
			return nil, ErrorNotCompatible
		}
		rValues := rValue.([]interface{})

		metric := ToMetric(rMetric.(map[string]interface{}))
		tsValues := []TSValue{}
		for _, v := range rValues {
			vs := v.([]interface{})
			timestamp := vs[0].(float64)
			value := ToValue(vs[1].(string))
			tsValues = append(tsValues, TSValue{Timestamp: timestamp, Value: value})
		}

		points = append(points, QueryPoint{
			Metric: metric,
			Values: tsValues,
		})
	}

	return points, nil
}

func ParseSeries(raw interface{}) ([]Series, error) {
	data := raw.([]interface{})
	ss := []Series{}
	for _, d := range data {
		s := Series{}
		m := d.(map[string]interface{})
		for k, v := range m {
			s[k] = v.(string)
		}
		ss = append(ss, s)
	}
	return ss, nil
}

func ParseLabels(raw interface{}) (Labels, error) {
	data := raw.([]interface{})
	labels := Labels{}
	for _, d := range data {
		labels = append(labels, d.(string))
	}
	return labels, nil
}

func ParseLabelValues(raw interface{}) (LabelValues, error) {
	data := raw.([]interface{})
	values := LabelValues{}
	for _, d := range data {
		values = append(values, d.(string))
	}
	return values, nil
}

func ParseTarget(raw interface{}) (*Target, error) {
	target := Target{}
	bs, _ := json.Marshal(raw)
	if err := json.Unmarshal(bs, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

func ToMetric(raw map[string]interface{}) map[string]string {
	m := map[string]string{}
	for key, val := range raw {
		switch v := val.(type) {
		case string:
			m[key] = v
		default:
			m[key] = "Unsupported"
		}
	}
	return m
}

func ToValue(raw string) float64 {
	if raw == "+Inf" {
		return math.Inf(1)
	}
	v, _ := strconv.ParseFloat(raw, 64)
	return v
}
