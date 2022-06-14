package synthetic

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn-contrib/dynatrace-service/internal/adapter"
	"github.com/keptn-contrib/dynatrace-service/internal/synthetic/connector"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

type SyntheticTriggerStartedEventData struct {
	keptnv2.EventData
}

type SyntheticExecution struct {
	BatchId          string                             `json:"batchId"`
	ExecutionIds     []string                           `json:"executionIds"`
	FailedTriggers   []connector.ExecutionNotTriggered  `json:"failedTriggers"`
	FailedExecutions []connector.ExecutionNotSuccessful `json:"failedExecutions"`
	SuccessRate      float64                            `json:"successRate"`
}

type SyntheticTriggerFinishedEventData struct {
	keptnv2.EventData
	SyntheticExecution SyntheticExecution `json:"syntheticExecution"`
}

// SyntheticTriggerStartedEventFactory is a factory for test.started cloud events.
type SyntheticTriggerStartedEventFactory struct {
	event SyntheticTriggerAdapterInterface
}

// NewSyntheticTriggerStartedEventFactory creates a new SyntheticTriggerStartedEventFactory.
func NewSyntheticTriggerStartedEventFactory(event SyntheticTriggerAdapterInterface) *SyntheticTriggerStartedEventFactory {
	return &SyntheticTriggerStartedEventFactory{
		event: event,
	}
}

// CreateCloudEvent creates a cloud event based on the factory or returns an error if this can't be done.
func (f *SyntheticTriggerStartedEventFactory) CreateCloudEvent() (*cloudevents.Event, error) {
	startedEvent := SyntheticTriggerStartedEventData{
		EventData: keptnv2.EventData{
			Project: f.event.GetProject(),
			Stage:   f.event.GetStage(),
			Service: f.event.GetService(),
			Labels:  f.event.GetLabels(),
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
		},
	}

	return adapter.NewCloudEventFactory(f.event, keptnv2.GetStartedEventType(keptnv2.TestTaskName), startedEvent).CreateCloudEvent()

}

// SyntheticTriggerFinishedEventFactory is a factory for test.finished cloud events.
type SyntheticTriggerFinishedEventFactory struct {
	event         SyntheticTriggerAdapterInterface
	status        keptnv2.StatusType
	result        keptnv2.ResultType
	err           error
	executionData connector.ExecutionData
}

// NewSucceededSyntheticTriggerFinishedEventFactory creates a new SyntheticTriggerFinishedEventFactory with status succeeded.
func NewSucceededSyntheticTriggerFinishedEventFactory(event SyntheticTriggerAdapterInterface, executionData connector.ExecutionData, err error) *SyntheticTriggerFinishedEventFactory {
	return &SyntheticTriggerFinishedEventFactory{
		event:         event,
		status:        keptnv2.StatusSucceeded,
		result:        keptnv2.ResultPass,
		err:           err,
		executionData: executionData,
	}
}

// NewErroredSyntheticTriggerFinishedEventFactory creates a new SyntheticTriggerFinishedEventFactory with status errored.
func NewErroredSyntheticTriggerFinishedEventFactory(event SyntheticTriggerAdapterInterface, executionData connector.ExecutionData, err error) *SyntheticTriggerFinishedEventFactory {
	return &SyntheticTriggerFinishedEventFactory{
		event:         event,
		status:        keptnv2.StatusErrored,
		result:        keptnv2.ResultFailed,
		err:           err,
		executionData: executionData,
	}
}

// NewWarningSyntheticTriggerFinishedEventFactory creates a new SyntheticTriggerFinishedEventFactory with status unknown, result warning.
func NewWarningSyntheticTriggerFinishedEventFactory(event SyntheticTriggerAdapterInterface, executionData connector.ExecutionData, err error) *SyntheticTriggerFinishedEventFactory {
	return &SyntheticTriggerFinishedEventFactory{
		event:         event,
		status:        keptnv2.StatusUnknown,
		result:        keptnv2.ResultWarning,
		err:           err,
		executionData: executionData,
	}
}

// CreateCloudEvent creates a cloud event based on the factory or returns an error if this can't be done.
func (f *SyntheticTriggerFinishedEventFactory) CreateCloudEvent() (*cloudevents.Event, error) {
	msg := ""

	if f.err != nil {
		msg = f.err.Error()
	}

	finishedEvent := SyntheticTriggerFinishedEventData{
		EventData: keptnv2.EventData{
			Project: f.event.GetProject(),
			Stage:   f.event.GetStage(),
			Service: f.event.GetService(),
			Labels:  f.event.GetLabels(),
			Status:  f.status,
			Result:  f.result,
			Message: msg,
		},
		SyntheticExecution: SyntheticExecution{
			BatchId:          f.executionData.BatchId,
			ExecutionIds:     f.executionData.ExecutionIds,
			FailedTriggers:   f.executionData.FailedTriggers,
			FailedExecutions: f.executionData.FailedExecutions,
			SuccessRate:      f.executionData.SuccessRate,
		},
	}

	return adapter.NewCloudEventFactory(f.event, keptnv2.GetFinishedEventType(keptnv2.TestTaskName), finishedEvent).CreateCloudEvent()
}
