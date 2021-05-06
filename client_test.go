package promclient

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	testClient *Client
	testQExp1  = "up"
	testQExp2  = "apiserver_request_latencies_sum"

	currentTs = float64(time.Now().Unix())
	startTs   = float64(time.Now().Unix()) - 5*60
	endTs     = float64(time.Now().Unix())
)

func TestMain(m *testing.M) {
	testClient = NewClient(os.Getenv("URL"), 0)
	os.Exit(m.Run())
}

func TestQuery1(t *testing.T) {
	// points, err := testClient.Query(testQExp1, currentTs)
	_, err := testClient.Query(testQExp1, currentTs)
	require.Nil(t, err)

	// for _, p := range points {
	// 	fmt.Println(p)
	// }
}

func TestQuery2(t *testing.T) {
	// points, err := testClient.Query(testQExp1, currentTs)
	_, err := testClient.Query(testQExp1, currentTs)
	require.Nil(t, err)

	// for _, p := range points {
	// 	fmt.Println(p)
	// }
}

func TestQueryRange1(t *testing.T) {
	// points, err := testClient.QueryRange(testQExp1, startTs, endTs, 0, 0)
	_, err := testClient.QueryRange(testQExp1, startTs, endTs, 0, 0)
	require.Nil(t, err)

	// for _, p := range points {
	// 	fmt.Println(p)
	// }
}

func TestQueryRange2(t *testing.T) {
	// points, err := testClient.QueryRange(testQExp2, startTs, endTs, 2*time.Minute, 0)
	_, err := testClient.QueryRange(testQExp2, startTs, endTs, 2*time.Minute, 0)
	require.Nil(t, err)

	// for _, p := range points {
	// 	fmt.Println(p)
	// }
}

func TestQuerySeries(t *testing.T) {
	matches := []string{
		"up",
		"apiserver_request_latencies_sum",
	}
	// ss, err := testClient.QuerySeries(matches, startTs, endTs)
	_, err := testClient.QuerySeries(matches, startTs, endTs)
	require.Nil(t, err)

	// for _, s := range ss {
	// 	fmt.Println(s)
	// }
}

func TestQueryLabels(t *testing.T) {
	// labels, err := testClient.QueryLabels([]string{"up"}, startTs, endTs)
	_, err := testClient.QueryLabels([]string{"up"}, startTs, endTs)
	require.Nil(t, err)

	// for _, l := range labels {
	// 	fmt.Println(l)
	// }
}

func TestQueryLabelValues(t *testing.T) {
	// values, err := testClient.QueryLabelValues("instance", []string{}, startTs, endTs)
	_, err := testClient.QueryLabelValues("instance", []string{}, startTs, endTs)
	require.Nil(t, err)

	// for _, v := range values {
	// 	fmt.Println(v)
	// }
}

func TestQueryTarget(t *testing.T) {
	// target, err := testClient.QueryTargets("")
	_, err := testClient.QueryTargets("")
	require.Nil(t, err)

	// fmt.Println(target.ActiveTargets)
	// fmt.Println(target.DroppedTargets)
}
