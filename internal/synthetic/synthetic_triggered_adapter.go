package synthetic

import (
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn-contrib/dynatrace-service/internal/adapter"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

type SyntheticTriggerAdapterInterface interface {
	adapter.EventContentAdapter
	adapter.TriggeredCloudEventContentAdapter

	GetSyntheticMonitorId() string
	GetSyntheticMonitorTag() string
	IsWaitForDataRequested() bool
	IsWaitForExecutionRequested() bool
}

type TestEventData struct {
	MonitorTag string `json:"monitorTag"`
	MonitorId  string `json:"monitorId"`
	WaitFor    string `json:"waitFor"`
}

type SyntheticTriggerEventData struct {
	keptnv2.EventData
	MonitorTag string        `json:"monitorTag"`
	MonitorId  string        `json:"monitorId"`
	WaitFor    string        `json:"waitFor"`
	Test       TestEventData `json:"test"`
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

// GetSyntheticMonitorId returns the used synthetic monitor id
func (a SyntheticTriggerAdapter) GetSyntheticMonitorId() string {
	isDefinedInTestAttribute := a.event.Test.MonitorId != ""
	if isDefinedInTestAttribute {
		return a.event.Test.MonitorId
	} else {
		return a.event.MonitorId
	}
}

// GetSyntheticMonitorTag returns the used synthetic monitor tag
func (a SyntheticTriggerAdapter) GetSyntheticMonitorTag() string {
	isDefinedInTestAttribute := a.event.Test.MonitorTag != ""
	if isDefinedInTestAttribute {
		return a.event.Test.MonitorTag
	} else {
		return a.event.MonitorTag
	}
}

// IsWaitForDataRequested returns whether the synthetic monitor shall wait for data retrieval
func (a SyntheticTriggerAdapter) IsWaitForDataRequested() bool {
	isDefinedInTestAttribute := a.event.Test.WaitFor != ""
	if isDefinedInTestAttribute {
		return strings.ToLower(a.event.Test.WaitFor) == "data"
	} else {
		return strings.ToLower(a.event.WaitFor) == "data"
	}
}

// IsWaitForExecutionRequested returns whether the synthetic monitor shall wait for synthetic execution
func (a SyntheticTriggerAdapter) IsWaitForExecutionRequested() bool {
	isDefinedInTestAttribute := a.event.Test.WaitFor != ""
	if isDefinedInTestAttribute {
		return strings.ToLower(a.event.Test.WaitFor) == "execution"
	} else {
		return strings.ToLower(a.event.WaitFor) == "execution"
	}
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
