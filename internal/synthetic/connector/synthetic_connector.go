package connector

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/keptn-contrib/dynatrace-service/internal/dynatrace"
	log "github.com/sirupsen/logrus"
)

const syntheticTriggerPath = "/api/v2/synthetic/monitors/execute"

type SyntheticConnectorInterface interface {
	TriggerById(workCtx context.Context, monitorId string) (ExecutionData, error)
	TriggerByTag(workCtx context.Context, monitorTag string) (ExecutionData, error)
}

type SyntheticConnector struct {
	dtClient      dynatrace.ClientInterface
	executionData ExecutionData
}

type ExecutionResponseBody struct {
	BatchId           string                  `json:"batchId"`
	NotTriggeredCount int16                   `json:"notTriggeredCount"`
	NotTriggered      []ExecutionNotTriggered `json:"notTriggered"`
	TriggeredCount    int16                   `json:"triggeredCount"`
	Triggered         []ExecutionTriggered    `json:"triggered"`
}

type ExecutionNotTriggered struct {
	MonitorId  string `json:"monitorId"`
	LocationId string `json:"locationId"`
	Cause      string `json:"cause"`
}

type ExecutionTriggered struct {
	MonitorId  string `json:"monitorId"`
	Executions []struct {
		ExecutionId string `json:"executionId"`
		LocationId  string `json:"locationId"`
	} `json:"executions"`
}

type ExecutionData struct {
	BatchId          string                  `json:"batchId"`
	ExecutionIds     []string                `json:"executionIds"`
	FailedExecutions []ExecutionNotTriggered `json:"failedExecutions"`
}

func generateExecutionByIdEvent(monitorId string) []byte {
	jsonData := []byte(fmt.Sprintf(`{
		"monitors": [
			{
				"monitorId": "%s",
				"locations": []
			}
		]
	}`, monitorId))

	return jsonData
}

func generateExecutionByTagEvent(monitorTag string) []byte {
	jsonData := []byte(fmt.Sprintf(`{
		"group": {
			"tags": [
				"%s"
			]
		}
	}`, monitorTag))

	return jsonData
}

func parseExecutionIds(executionResponseBody ExecutionResponseBody) []string {
	executionIds := []string{}

	for _, triggered := range executionResponseBody.Triggered {
		for _, execution := range triggered.Executions {
			executionIds = append(executionIds, execution.ExecutionId)
		}
	}

	return executionIds
}

func parseFailedExecutions(executionResponseBody ExecutionResponseBody) []ExecutionNotTriggered {
	executions := executionResponseBody.NotTriggered
	return executions
}

func parseBatchId(executionResponseBody ExecutionResponseBody) string {
	batchId := executionResponseBody.BatchId
	return batchId
}

func (sc *SyntheticConnector) TriggerById(workCtx context.Context, monitorId string) (ExecutionData, error) {
	jsonData := generateExecutionByIdEvent(monitorId)
	return sc.trigger(workCtx, jsonData)
}

func (sc *SyntheticConnector) TriggerByTag(workCtx context.Context, monitorTag string) (ExecutionData, error) {
	jsonData := generateExecutionByTagEvent(monitorTag)
	return sc.trigger(workCtx, jsonData)
}

func (sc *SyntheticConnector) trigger(workCtx context.Context, jsonData []byte) (ExecutionData, error) {
	triggerData := ExecutionData{}

	resp, err := sc.dtClient.Post(workCtx, syntheticTriggerPath, jsonData)
	if err != nil {
		return ExecutionData{}, err
	}

	executionResponseBody := ExecutionResponseBody{}
	err = json.Unmarshal(resp, &executionResponseBody)
	if err != nil {
		log.Error(err.Error())
		return ExecutionData{}, err
	}

	triggerData.BatchId = parseBatchId(executionResponseBody)
	triggerData.ExecutionIds = parseExecutionIds(executionResponseBody)
	triggerData.FailedExecutions = parseFailedExecutions(executionResponseBody)

	return ExecutionData{
		BatchId:          parseBatchId(executionResponseBody),
		ExecutionIds:     parseExecutionIds(executionResponseBody),
		FailedExecutions: parseFailedExecutions(executionResponseBody),
	}, nil
}

func NewSyntheticConnector(dtClient dynatrace.ClientInterface) *SyntheticConnector {
	return &SyntheticConnector{
		dtClient: dtClient,
	}
}
