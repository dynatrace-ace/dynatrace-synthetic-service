package synthetic

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn-contrib/dynatrace-service/internal/adapter"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

type SyntheticTriggeredAdapterInterface interface {
	adapter.EventContentAdapter
	GetSyntheticMonitorId() string
}

type SyntheticTriggeredEventData struct {
	keptnv2.EventData
	MonitorId string `json:"monitorId"`
}

// SyntheticTriggeredAdapter is a content adaptor for events of type sh.keptn.event.test.triggered
type SyntheticTriggeredAdapter struct {
	event      SyntheticTriggeredEventData
	cloudEvent adapter.CloudEventAdapter
}

// NewSyntheticTriggeredAdapterFromEvent creates a new SyntheticTriggeredAdapter from a cloudevents Event
func NewSyntheticTriggeredAdapterFromEvent(e cloudevents.Event) (*SyntheticTriggeredAdapter, error) {
	ceAdapter := adapter.NewCloudEventAdapter(e)

	ttData := &SyntheticTriggeredEventData{}
	err := ceAdapter.PayloadAs(ttData)
	if err != nil {
		return nil, err
	}

	return &SyntheticTriggeredAdapter{
		event:      *ttData,
		cloudEvent: ceAdapter,
	}, nil
}

// GetShKeptnContext returns the shkeptncontext
func (a SyntheticTriggeredAdapter) GetShKeptnContext() string {
	return a.cloudEvent.GetShKeptnContext()
}

// GetSource returns the source specified in the CloudEvent context
func (a SyntheticTriggeredAdapter) GetSource() string {
	return a.cloudEvent.GetSource()
}

// GetEvent returns the event type
func (a SyntheticTriggeredAdapter) GetEvent() string {
	return keptnv2.GetFinishedEventType(keptnv2.TestTaskName)
}

// GetProject returns the project
func (a SyntheticTriggeredAdapter) GetProject() string {
	return a.event.Project
}

// GetStage returns the stage
func (a SyntheticTriggeredAdapter) GetStage() string {
	return a.event.Stage
}

// GetService returns the service
func (a SyntheticTriggeredAdapter) GetService() string {
	return a.event.Service
}

// GetDeployment returns the name of the deployment
func (a SyntheticTriggeredAdapter) GetDeployment() string {
	return ""
}

// GetTestStrategy returns the used test strategy
func (a SyntheticTriggeredAdapter) GetTestStrategy() string {
	return ""
}

// GetMonitorId returns the used synthetic monitor id
func (a SyntheticTriggeredAdapter) GetSyntheticMonitorId() string {
	return a.event.MonitorId
}

// GetDeploymentStrategy returns the used deployment strategy
func (a SyntheticTriggeredAdapter) GetDeploymentStrategy() string {
	return ""
}

// GetLabels returns a map of labels
func (a SyntheticTriggeredAdapter) GetLabels() map[string]string {
	return a.event.Labels
}
