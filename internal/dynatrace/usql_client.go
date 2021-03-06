package dynatrace

import (
	"context"
	"encoding/json"
	"time"

	"github.com/keptn-contrib/dynatrace-service/internal/common"
	"github.com/keptn-contrib/dynatrace-service/internal/sli/usql"
)

const USQLPath = "/api/v1/userSessionQueryLanguage/table"

// USQLRequiredDelay is delay required between the end of a timeframe and an USQL API request using it.
const USQLRequiredDelay = 6 * time.Minute

// USQLMaximumWait is maximum acceptable wait time between the end of a timeframe and an USQL API request using it.
const USQLMaximumWait = 8 * time.Minute

const (
	queryKey             = "query"
	explainKey           = "explain"
	addDeepLinkFieldsKey = "addDeepLinkFields"
	startTimestampKey    = "startTimestamp"
	endTimestampKey      = "endTimestamp"
)

// USQLClientQueryParameters encapsulates the query parameters for the USQLClient's GetByQuery method.
type USQLClientQueryParameters struct {
	query     usql.Query
	timeframe common.Timeframe
}

// NewUSQLClientQueryParameters creates new USQLClientQueryParameters.
func NewUSQLClientQueryParameters(query usql.Query, timeframe common.Timeframe) USQLClientQueryParameters {
	return USQLClientQueryParameters{
		query:     query,
		timeframe: timeframe,
	}
}

// encode encodes USQLClientQueryParameters into a URL-encoded string.
func (q *USQLClientQueryParameters) encode() string {
	queryParameters := newQueryParameters()
	queryParameters.add(queryKey, q.query.GetQuery())
	queryParameters.add(explainKey, "false")
	queryParameters.add(addDeepLinkFieldsKey, "false")
	queryParameters.add(startTimestampKey, common.TimestampToUnixMillisecondsString(q.timeframe.Start()))
	queryParameters.add(endTimestampKey, common.TimestampToUnixMillisecondsString(q.timeframe.End()))
	return queryParameters.encode()
}

// DTUSQLResult struct
type DTUSQLResult struct {
	ExtrapolationLevel int             `json:"extrapolationLevel"`
	ColumnNames        []string        `json:"columnNames"`
	Values             [][]interface{} `json:"values"`
}

type USQLClient struct {
	client ClientInterface
}

func NewUSQLClient(client ClientInterface) *USQLClient {
	return &USQLClient{
		client: client,
	}
}

// GetByQuery executes the passed USQL API query, validates that the call returns data and returns the data set.
func (uc *USQLClient) GetByQuery(ctx context.Context, parameters USQLClientQueryParameters) (*DTUSQLResult, error) {
	err := NewTimeframeDelay(parameters.timeframe, USQLRequiredDelay, USQLMaximumWait).Wait(ctx)
	if err != nil {
		return nil, err
	}

	body, err := uc.client.Get(ctx, USQLPath+"?"+parameters.encode())
	if err != nil {
		return nil, err
	}

	// parse response json
	var result DTUSQLResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
