package synthetic

import (
	"context"

	"github.com/keptn-contrib/dynatrace-service/internal/adapter"
	"github.com/keptn-contrib/dynatrace-service/internal/dynatrace"
	"github.com/keptn-contrib/dynatrace-service/internal/keptn"
	"github.com/keptn-contrib/dynatrace-service/internal/synthetic/connector"
	log "github.com/sirupsen/logrus"
)

// SyntheticTriggerEventHandler handles a test triggered event.
type SyntheticTriggerEventHandler struct {
	event       SyntheticTriggerAdapterInterface
	dtClient    dynatrace.ClientInterface
	sClient     connector.SyntheticConnectorInterface
	kClient     keptn.ClientInterface
	eClient     keptn.EventClientInterface
	attachRules *dynatrace.AttachRules
}

// NewSyntheticTriggerEventHandler creates a new SyntheticTriggerEventHandler.
func NewSyntheticTriggerEventHandler(event SyntheticTriggerAdapterInterface, dtClient dynatrace.ClientInterface, sClient connector.SyntheticConnectorInterface, kClient keptn.ClientInterface, eClient keptn.EventClientInterface, attachRules *dynatrace.AttachRules) *SyntheticTriggerEventHandler {
	return &SyntheticTriggerEventHandler{
		event:       event,
		dtClient:    dtClient,
		sClient:     sClient,
		kClient:     kClient,
		eClient:     eClient,
		attachRules: attachRules,
	}
}

// HandleEvent handles a test triggered event.
func (eh *SyntheticTriggerEventHandler) HandleEvent(workCtx context.Context, replyCtx context.Context) error {
	syntheticMonitorId := eh.event.GetSyntheticMonitorId()
	syntheticMonitorTag := eh.event.GetSyntheticMonitorTag()

	isMonitorIdDefined := syntheticMonitorId != ""
	isMonitorTagDefined := syntheticMonitorTag != ""

	if !isMonitorIdDefined && !isMonitorTagDefined {
		log.Info("Neither monitor id nor tag provided. Skipping handler...")
		return nil
	}

	err := eh.sendTriggerSyntheticStartedEvent()
	if err != nil {
		return err
	}

	sClient := connector.NewSyntheticConnector(eh.dtClient)

	executionData := connector.ExecutionData{}

	if isMonitorTagDefined {
		executionData, err = sClient.TriggerByTag(workCtx, syntheticMonitorTag)
		if err != nil {
			eh.sendFailedTriggerSyntheticFinishedEvent(executionData, err)
			return nil
		}
	} else {
		executionData, err = sClient.TriggerById(workCtx, syntheticMonitorId)
		if err != nil {
			eh.sendFailedTriggerSyntheticFinishedEvent(executionData, err)
			return nil
		}
	}

	isWaitForExecutionRequested := eh.event.IsWaitForExecutionRequested()
	// isWaitForDataRequested := eh.event.IsWaitForDataRequested()

	if isWaitForExecutionRequested {
		batchResponseBody, successRate, err := sClient.WaitForBatchExecution(workCtx)
		if err != nil {
			eh.sendWarningfulTriggerSyntheticFinishedEvent(executionData, err)
			return err
		}

		executionData.FailedExecutions = batchResponseBody.FailedExecutions
		executionData.SuccessRate = successRate

		_, err = sClient.IngestSyntheticSuccessMetric(workCtx, syntheticMonitorId, eh.event.GetProject(), eh.event.GetService(), eh.event.GetStage(), executionData.BatchId, successRate)
		if err != nil {
			eh.sendWarningfulTriggerSyntheticFinishedEvent(executionData, err)
			return err
		}

		// } else if isWaitForDataRequested {
		// 	err = eh.sendSuccessfulTriggerSyntheticFinishedEvent(executionData)
		// 	if err != nil {
		//    eh.sendFailedTriggerSyntheticFinishedEvent(executionData, err)
		// 		return err
		// 	}

		// 	return nil
	}

	err = eh.sendSuccessfulTriggerSyntheticFinishedEvent(executionData)
	if err != nil {
		return err
	}

	return nil
}

func (eh *SyntheticTriggerEventHandler) sendTriggerSyntheticStartedEvent() error {
	return eh.sendEvent(NewSyntheticTriggerStartedEventFactory(eh.event))
}

func (eh *SyntheticTriggerEventHandler) sendSuccessfulTriggerSyntheticFinishedEvent(executionData connector.ExecutionData) error {
	return eh.sendEvent(NewSucceededSyntheticTriggerFinishedEventFactory(eh.event, executionData, nil))
}

func (eh *SyntheticTriggerEventHandler) sendWarningfulTriggerSyntheticFinishedEvent(executionData connector.ExecutionData, err error) error {
	return eh.sendEvent(NewWarningSyntheticTriggerFinishedEventFactory(eh.event, executionData, err))
}

func (eh *SyntheticTriggerEventHandler) sendFailedTriggerSyntheticFinishedEvent(executionData connector.ExecutionData, err error) error {
	return eh.sendEvent(NewErroredSyntheticTriggerFinishedEventFactory(eh.event, executionData, err))
}

func (eh *SyntheticTriggerEventHandler) sendEvent(factory adapter.CloudEventFactoryInterface) error {
	err := eh.kClient.SendCloudEvent(factory)
	if err != nil {
		log.WithError(err).Error("Could not send get sli cloud event")
		return err
	}

	return nil
}
