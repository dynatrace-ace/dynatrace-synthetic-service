package dashboard

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keptn-contrib/dynatrace-service/internal/adapter"
	"github.com/keptn-contrib/dynatrace-service/internal/common"
	"github.com/keptn-contrib/dynatrace-service/internal/dynatrace"
	"github.com/keptn-contrib/dynatrace-service/internal/test"
)

func TestFindDynatraceDashboardSuccess(t *testing.T) {
	keptnEvent := createKeptnEvent(QUALITYGATE_PROJECT, QUALITYGATE_STAGE, QUALTIYGATE_SERVICE)

	handler := test.NewFileBasedURLHandler(t)
	handler.AddExact("/api/config/v1/dashboards", "./testdata/test_get_dashboards.json")

	dh, teardown := createDashboardRetrieval(t, keptnEvent, handler)
	defer teardown()

	dashboardID, err := dh.findDynatraceDashboard(context.TODO())

	assert.NoError(t, err)
	assert.EqualValues(t, dashboardID, QUALITYGATE_DASHBOARD_ID)
}

func TestFindDynatraceDashboardNoneExistingDashboard(t *testing.T) {
	keptnEvent := createKeptnEvent("BAD PROJECT", QUALITYGATE_STAGE, QUALTIYGATE_SERVICE)

	handler := test.NewFileBasedURLHandler(t)
	handler.AddExact("/api/config/v1/dashboards", "./testdata/test_get_dashboards.json")

	dh, teardown := createDashboardRetrieval(t, keptnEvent, handler)
	defer teardown()

	dashboardID, err := dh.findDynatraceDashboard(context.TODO())

	assert.Error(t, err)
	assert.Empty(t, dashboardID)
}

func TestLoadDynatraceDashboardWithQUERY(t *testing.T) {
	keptnEvent := createKeptnEvent(QUALITYGATE_PROJECT, QUALITYGATE_STAGE, QUALTIYGATE_SERVICE)

	handler := test.NewFileBasedURLHandler(t)
	handler.AddExact("/api/config/v1/dashboards", "./testdata/test_get_dashboards.json")
	handler.AddExact("/api/config/v1/dashboards/12345678-1111-4444-8888-123456789012", "./testdata/test_get_dashboards_id.json")

	dh, teardown := createDashboardRetrieval(t, keptnEvent, handler)
	defer teardown()

	dashboard, dashboardID, err := dh.Retrieve(context.TODO(), common.DynatraceConfigDashboardQUERY)

	assert.NoError(t, err)
	assert.NotNil(t, dashboard)
	assert.EqualValues(t, QUALITYGATE_DASHBOARD_ID, dashboardID)
}

func TestLoadDynatraceDashboardWithID(t *testing.T) {
	keptnEvent := createKeptnEvent(QUALITYGATE_PROJECT, QUALITYGATE_STAGE, QUALTIYGATE_SERVICE)

	handler := test.NewFileBasedURLHandler(t)
	handler.AddExact("/api/config/v1/dashboards/12345678-1111-4444-8888-123456789012", "./testdata/test_get_dashboards_id.json")

	dh, teardown := createDashboardRetrieval(t, keptnEvent, handler)
	defer teardown()

	dashboard, dashboardID, err := dh.Retrieve(context.TODO(), QUALITYGATE_DASHBOARD_ID)

	assert.NoError(t, err)
	assert.NotNil(t, dashboard)
	assert.EqualValues(t, QUALITYGATE_DASHBOARD_ID, dashboardID)
}

func TestLoadDynatraceDashboardWithEmptyDashboard(t *testing.T) {
	keptnEvent := createKeptnEvent(QUALITYGATE_PROJECT, QUALITYGATE_STAGE, QUALTIYGATE_SERVICE)

	handler := test.NewFileBasedURLHandler(t)

	dh, teardown := createDashboardRetrieval(t, keptnEvent, handler)
	defer teardown()

	dashboardJSON, dashboard, err := dh.Retrieve(context.TODO(), "")

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid 'dashboard'")
	}
	assert.Nil(t, dashboardJSON)
	assert.Empty(t, dashboard)
}

func createDashboardRetrieval(t *testing.T, eventData adapter.EventContentAdapter, handler http.Handler) (*Retrieval, func()) {
	httpClient, url, teardown := test.CreateHTTPSClient(handler)

	retrieval := NewRetrieval(
		dynatrace.NewClientWithHTTP(createDynatraceCredentials(t, url), httpClient),
		eventData)

	return retrieval, teardown
}
