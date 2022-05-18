package synthetic

import (
	"context"

	"github.com/keptn-contrib/dynatrace-service/internal/dynatrace"
	"github.com/keptn-contrib/dynatrace-service/internal/keptn"
	log "github.com/sirupsen/logrus"
)

// SyntheticTriggeredEventHandler handles a test triggered event.
type SyntheticTriggeredEventHandler struct {
	event       SyntheticTriggeredAdapterInterface
	dtClient    dynatrace.ClientInterface
	kClient     keptn.ClientInterface
	eClient     keptn.EventClientInterface
	attachRules *dynatrace.AttachRules
}

// NewSyntheticTriggeredEventHandler creates a new SyntheticTriggeredEventHandler.
func NewSyntheticTriggeredEventHandler(event SyntheticTriggeredAdapterInterface, dtClient dynatrace.ClientInterface, kClient keptn.ClientInterface, eClient keptn.EventClientInterface, attachRules *dynatrace.AttachRules) *SyntheticTriggeredEventHandler {
	return &SyntheticTriggeredEventHandler{
		event:       event,
		dtClient:    dtClient,
		kClient:     kClient,
		eClient:     eClient,
		attachRules: attachRules,
	}
}

// HandleEvent handles a test triggered event.
func (eh *SyntheticTriggeredEventHandler) HandleEvent(workCtx context.Context, replyCtx context.Context) error {
	syntheticMonitorId := eh.event.GetSyntheticMonitorId()
	if syntheticMonitorId == "" {
		log.Info("No monitor ID provided. Skipping handler...")
		return nil
	}

	// TBD: Send test started event
	log.Info("TBD: Trigger synthetic %s", syntheticMonitorId)

	return nil
}
