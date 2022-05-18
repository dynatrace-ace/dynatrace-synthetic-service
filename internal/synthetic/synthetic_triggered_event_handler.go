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
	syntheticMonitorTag := eh.event.GetSyntheticMonitorTag()
	if syntheticMonitorTag == "" {
		log.Info("No monitor tag provided. Skipping handler...")
		return nil
	}

	err := eh.sendTriggerSyntheticStartedEvent()
	if err != nil {
		return err
	}

	sClient := connector.NewSyntheticConnector(eh.dtClient)

	executionData, err := sClient.Trigger(workCtx, syntheticMonitorTag)
	if err != nil {
		eh.sendFailedTriggerSyntheticFinishedEvent(executionData, err)
		return nil
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
