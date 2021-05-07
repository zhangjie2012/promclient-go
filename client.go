package promclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	url     string
	timeout time.Duration

	httpClient *http.Client
}

func NewClient(url string, timeout time.Duration) *Client {
	if timeout == 0 {
		timeout = 15 * time.Second
	}

	c := Client{
		url:     strings.TrimRight(url, "/"),
		timeout: timeout,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
	return &c
}

func (c *Client) request(api string, params url.Values) (interface{}, error) {
	// resp, err := c.httpClient.Post(api, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	fullUrl := fmt.Sprintf("%s?%s", api, params.Encode())
	// fmt.Println(fullUrl)
	resp, err := c.httpClient.Get(fullUrl)
	if err != nil {
		return nil, err
	}

	if err := HttpCodeTrans(resp.StatusCode); err != nil {
		return nil, err
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawResp := PromResponse{}
	if err := json.Unmarshal(bs, &rawResp); err != nil {
		return nil, err
	}

	// status: "success"/"error"
	if rawResp.Status != "success" {
		return nil, fmt.Errorf("type=%s, error=%s", rawResp.ErrorType, rawResp.Error)
	}

	if rawResp.Data == nil {
		return nil, ErrorNoData
	}

	return rawResp.Data, nil
}

// Query instant query at a single point in time.
// - timestamp: Unix timestamp in seconds, with optional decimal places for sub-second precision
// docs: https://prometheus.io/docs/prometheus/latest/querying/api/#instant-queries
func (c *Client) Query(q string, timestamp float64) ([]QueryPoint, error) {
	params := url.Values{}
	params.Add("query", q)
	params.Add("time", strconv.FormatFloat(timestamp, 'f', -1, 64))
	params.Add("timeout", fmt.Sprintf("%ds", int(c.timeout.Seconds())))

	api := fmt.Sprintf("%s%s", c.url, "/api/v1/query")

	rawData, err := c.request(api, params)
	if err != nil {
		return nil, err
	}

	points, err := ParseQueryResult(rawData)

	return points, err
}

// QueryRange query over a range of time
// - startTs, endTs: time range (second)
// - step: query resolution step; if 0, default 30s
// - timeout: query timeout; if 0, default c.timeout
// docs: https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries
func (c *Client) QueryRange(q string, startTs, endTs float64, step, timeout time.Duration) ([]QueryPoint, error) {
	if step == 0 {
		step = 30 * time.Second
	}
	if timeout == 0 {
		timeout = c.timeout
	}

	params := url.Values{}
	params.Add("query", q)
	params.Add("start", strconv.FormatFloat(startTs, 'f', -1, 64))
	params.Add("end", strconv.FormatFloat(endTs, 'f', -1, 64))
	params.Add("step", fmt.Sprintf("%ds", int(step.Seconds())))
	params.Add("timeout", fmt.Sprintf("%ds", int(timeout.Seconds())))

	api := fmt.Sprintf("%s%s", c.url, "/api/v1/query_range")

	rawData, err := c.request(api, params)
	if err != nil {
		return nil, err
	}

	points, err := ParseQueryResult(rawData)

	return points, err
}

// QuerySeries finding series by label matchers
// - matches: series selector
// - startTs, endTs: time range (second)
// docs: https://prometheus.io/docs/prometheus/latest/querying/api/#finding-series-by-label-matchers
func (c *Client) QuerySeries(matches []string, startTs, endTs float64) ([]Series, error) {
	params := url.Values{}
	for _, m := range matches {
		params.Add("match[]", m)
	}
	params.Add("start", strconv.FormatFloat(startTs, 'f', -1, 64))
	params.Add("end", strconv.FormatFloat(endTs, 'f', -1, 64))

	api := fmt.Sprintf("%s%s", c.url, "/api/v1/series")

	rawData, err := c.request(api, params)
	if err != nil {
		return nil, err
	}

	series, err := ParseSeries(rawData)
	return series, err
}

// QueryLabels getting label names
// - matches: series selector
// - startTs, endTs: time range (second)
// docs: https://prometheus.io/docs/prometheus/latest/querying/api/#getting-label-names
func (c *Client) QueryLabels(matches []string, startTs, endTs float64) (Labels, error) {
	params := url.Values{}
	for _, m := range matches {
		params.Add("match[]", m)
	}
	params.Add("start", strconv.FormatFloat(startTs, 'f', -1, 64))
	params.Add("end", strconv.FormatFloat(endTs, 'f', -1, 64))

	api := fmt.Sprintf("%s%s", c.url, "/api/v1/labels")

	rawData, err := c.request(api, params)
	if err != nil {
		return nil, err
	}

	labels, err := ParseLabels(rawData)
	return labels, err
}

// QueryLabelValues returns a list of label names
// - matches: series selector
// - startTs, endTs: time range (second)
// docs: https://prometheus.io/docs/prometheus/latest/querying/api/#getting-label-names
func (c *Client) QueryLabelValues(labelName string, matches []string, startTs, endTs float64) (LabelValues, error) {
	params := url.Values{}
	for _, m := range matches {
		params.Add("match[]", m)
	}
	params.Add("start", strconv.FormatFloat(startTs, 'f', -1, 64))
	params.Add("end", strconv.FormatFloat(endTs, 'f', -1, 64))

	api := fmt.Sprintf("%s%s", c.url, fmt.Sprintf("/api/v1/label/%s/values", labelName))

	rawData, err := c.request(api, params)
	if err != nil {
		return nil, err
	}

	labels, err := ParseLabelValues(rawData)
	return labels, err
}

// QueryTargets return target discovery state
func (c *Client) QueryTargets(state TargetStateT) (*Target, error) {
	if state == "" {
		state = TargetStateAny
	}
	params := url.Values{}
	params.Add("state", string(state))

	api := fmt.Sprintf("%s%s", c.url, "/api/v1/targets")

	rawData, err := c.request(api, params)
	if err != nil {
		return nil, err
	}

	target, err := ParseTarget(rawData)
	return target, err
}
