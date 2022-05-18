package synthetic

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn-contrib/dynatrace-service/internal/adapter"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

type SyntheticTriggerAdapterInterface interface {
	adapter.EventContentAdapter
	adapter.TriggeredCloudEventContentAdapter

	GetSyntheticMonitorTag() string
}

type SyntheticTriggerEventData struct {
	keptnv2.EventData
	MonitorTag string `json:"monitorTag"`
}

// SyntheticTriggerAdapter is a content adaptor for events of type sh.keptn.event.test.triggered
type SyntheticTriggerAdapter struct {
	event      SyntheticTriggerEventData
	cloudEvent adapter.CloudEventAdapter
}

// NewSyntheticTriggerAdapterFromEvent creates a new SyntheticTriggerAdapter from a cloudevents Event
func NewSyntheticTriggerAdapterFromEvent(e cloudevents.Event) (*SyntheticTriggerAdapter, error) {
	ceAdapter := adapter.NewCloudEventAdapter(e)

	ttData := &SyntheticTriggerEventData{}
	err := ceAdapter.PayloadAs(ttData)
	if err != nil {
		return nil, err
	}

	return &SyntheticTriggerAdapter{
		event:      *ttData,
		cloudEvent: ceAdapter,
	}, nil
}

// GetShKeptnContext returns the shkeptncontext
func (a SyntheticTriggerAdapter) GetShKeptnContext() string {
	return a.cloudEvent.GetShKeptnContext()
}

// GetSource returns the source specified in the CloudEvent context
func (a SyntheticTriggerAdapter) GetSource() string {
	return a.cloudEvent.GetSource()
}

// GetEvent returns the event type
func (a SyntheticTriggerAdapter) GetEvent() string {
	return keptnv2.GetFinishedEventType(keptnv2.TestTaskName)
}

// GetProject returns the project
func (a SyntheticTriggerAdapter) GetProject() string {
	return a.event.Project
}

// GetStage returns the stage
func (a SyntheticTriggerAdapter) GetStage() string {
	return a.event.Stage
}

// GetService returns the service
func (a SyntheticTriggerAdapter) GetService() string {
	return a.event.Service
}

// GetDeployment returns the name of the deployment
func (a SyntheticTriggerAdapter) GetDeployment() string {
	return ""
}

// GetTestStrategy returns the used test strategy
func (a SyntheticTriggerAdapter) GetTestStrategy() string {
	return ""
}

// GetMonitorId returns the used synthetic monitor id
func (a SyntheticTriggerAdapter) GetSyntheticMonitorTag() string {
	return a.event.MonitorTag
}

// GetDeploymentStrategy returns the used deployment strategy
func (a SyntheticTriggerAdapter) GetDeploymentStrategy() string {
	return ""
}

// GetLabels returns a map of labels
func (a SyntheticTriggerAdapter) GetLabels() map[string]string {
	return a.event.Labels
}

func (a SyntheticTriggerAdapter) GetEventID() string {
	return a.cloudEvent.GetEventID()
}
