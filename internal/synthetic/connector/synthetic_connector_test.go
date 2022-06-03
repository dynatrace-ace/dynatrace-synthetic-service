package connector

import (
	"context"
	"os"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn-contrib/dynatrace-service/internal/credentials"
	"github.com/keptn-contrib/dynatrace-service/internal/dynatrace"
	"github.com/keptn-contrib/dynatrace-service/internal/env"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const mockSyntheticTestId = "TEST_SYNTHETIC_TEST_ID"
const mockProjectName = "TEST_PROJECT_NAME"
const mockServiceName = "TEST_SERVICE_NAME"
const mockStageName = "TEST_STAGE_NAME"
const mockBatchId = "TEST_BATCH_ID"
const mockGauge = 42

func init() {
	log.SetLevel(env.GetLogLevel())
}

func TestIngestSyntheticSuccessMetric(t *testing.T) {
	mockDynatraceCredentials, err := credentials.NewDynatraceCredentials(os.Getenv("DT_TENANT"), os.Getenv("DT_API_TOKEN"))
	assert.Nil(t, err)
	mockDtClient := dynatrace.NewClient(mockDynatraceCredentials)
	mockCtx := cloudevents.WithEncodingStructured(context.Background())

	mockPostData := []byte(generateMetricsIngestLine(mockSyntheticTestId, mockProjectName, mockServiceName, mockStageName, mockBatchId, mockGauge))

	_, err = mockDtClient.PostTextPlain(mockCtx, metricsIngestPath, mockPostData)
	assert.Nil(t, err)
}

func TestCalculateSuccessRate(t *testing.T) {
	mockBatchResponseBody := BatchResponseBody{
		TriggeredCount:       6,
		FailedCount:          1,
		FailedToExecuteCount: 1,
	}

	successRate, _ := calculateSuccessRate(mockBatchResponseBody)
	assert.Equal(t, 66.67, successRate)

	mockBatchResponseBody = BatchResponseBody{
		TriggeredCount:       0,
		FailedCount:          0,
		FailedToExecuteCount: 0,
	}

	successRate, _ = calculateSuccessRate(mockBatchResponseBody)
	assert.Equal(t, float64(0), successRate)
}
