package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/keptn-contrib/dynatrace-service/internal/dynatrace"
	log "github.com/sirupsen/logrus"
)

const syntheticBatchBasePath = "/api/v2/synthetic/executions/batch"
const metricsIngestPath = "/api/v2/metrics/ingest"

const executionSuccessMetricKey = "ca.synthetic.execution_success_rate"

func getSyntheticBatchPath(batchId string) string {
	return fmt.Sprintf("%s/%s", syntheticBatchBasePath, batchId)
}

type SyntheticConnectorInterface interface {
	TriggerById(workCtx context.Context, monitorId string) (ExecutionData, error)
	TriggerByTag(workCtx context.Context, monitorTag string) (ExecutionData, error)
	WaitForBatchExecution(workCtx context.Context) (BatchResponseBody, float64, error)
}

type SyntheticConnector struct {
	dtClient      dynatrace.ClientInterface
	executionData ExecutionData
}

type ExecutionResponseBody struct {
	BatchId                   string                  `json:"batchId"`
	TriggeringProblemsCount   int16                   `json:"triggeringProblemsCount"`
	TriggeringProblemsDetails []ExecutionNotTriggered `json:"triggeringProblemsDetails"`
	TriggeredCount            int16                   `json:"triggeredCount"`
	Triggered                 []ExecutionTriggered    `json:"triggered"`
}

type ExecutionNotTriggered struct {
	EntityId   string `json:"entityId"`
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
	BatchId          string                   `json:"batchId"`
	ExecutionIds     []string                 `json:"executionIds"`
	FailedTriggers   []ExecutionNotTriggered  `json:"failedTriggers"`
	FailedExecutions []ExecutionNotSuccessful `json:"failedExecutions"`
	SuccessRate      float64                  `json:"successRate"`
}

type ExecutionNotSuccessful struct {
	ExecutionId        string `json:"executionId"`
	ExecutionStage     string `json:"executionStage"`
	ExecutionTimestamp int    `json:"executionTimestamp"`
	MonitorId          string `json:"monitorId"`
	LocationId         string `json:"locationId"`
}

type BatchResponseBody struct {
	BatchStatus          string                   `json:"batchStatus"`
	TriggeredCount       int                      `json:"triggeredCount"`
	ExecutedCount        int                      `json:"executedCount"`
	FailedCount          int                      `json:"failedCount"`
	FailedToExecuteCount int                      `json:"failedToExecuteCount"`
	FailedExecutions     []ExecutionNotSuccessful `json:"failedExecutions"`
}

type IngestResponseBody struct {
	LinesOk      int `json:"linesOk"`
	LinesInvalid int `json:"linesInvalid"`
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

func parseFailedTriggers(executionResponseBody ExecutionResponseBody) []ExecutionNotTriggered {
	executions := executionResponseBody.TriggeringProblemsDetails
	return executions
}

func parseBatchId(executionResponseBody ExecutionResponseBody) string {
	batchId := executionResponseBody.BatchId
	return batchId
}

func (sc *SyntheticConnector) TriggerById(workCtx context.Context, monitorId string) (ExecutionData, error) {
	jsonData := generateExecutionByIdEvent(monitorId)

	log.Debug("TriggerById")
	log.Debug(string(jsonData))

	return sc.trigger(workCtx, jsonData)
}

func (sc *SyntheticConnector) TriggerByTag(workCtx context.Context, monitorTag string) (ExecutionData, error) {
	jsonData := generateExecutionByTagEvent(monitorTag)

	log.Debug("TriggerByTag")
	log.Debug(string(jsonData))

	return sc.trigger(workCtx, jsonData)
}

func (sc *SyntheticConnector) trigger(workCtx context.Context, jsonData []byte) (ExecutionData, error) {
	resp, err := sc.dtClient.Post(workCtx, syntheticBatchBasePath, jsonData)
	if err != nil {
		return ExecutionData{}, err
	}

	executionResponseBody := ExecutionResponseBody{}
	err = json.Unmarshal(resp, &executionResponseBody)
	if err != nil {
		log.Error(err.Error())
		return ExecutionData{}, err
	}

	sc.executionData.BatchId = parseBatchId(executionResponseBody)
	sc.executionData.ExecutionIds = parseExecutionIds(executionResponseBody)
	sc.executionData.FailedTriggers = parseFailedTriggers(executionResponseBody)

	return sc.executionData, nil
}

func generateMetricsIngestLine(syntheticTestId string, projectName string, serviceName string, stageName string, batchId string, gauge float64) string {
	return fmt.Sprintf(
		"%s,dt.entity.synthetic_test=%s,ca.project.name=%s,ca.service.name=%s,ca.stage.name=%s,ca.synthetic.batch_id=%s gauge,%f",
		executionSuccessMetricKey,
		syntheticTestId,
		projectName,
		serviceName,
		stageName,
		batchId,
		gauge,
	)
}

func (sc *SyntheticConnector) IngestSyntheticSuccessMetric(workCtx context.Context, syntheticTestId string, projectName string, serviceName string, stageName string, batchId string, gauge float64) (IngestResponseBody, error) {
	postData := []byte(generateMetricsIngestLine(syntheticTestId, projectName, serviceName, stageName, batchId, gauge))

	resp, err := sc.dtClient.PostTextPlain(workCtx, metricsIngestPath, postData)
	if err != nil {
		return IngestResponseBody{}, err
	}

	log.Debug(string(resp))

	ingestResponseBody := IngestResponseBody{}
	err = json.Unmarshal(resp, &ingestResponseBody)
	if err != nil {
		log.Error(err.Error())
		return IngestResponseBody{}, err
	}

	return ingestResponseBody, nil
}

func (sc *SyntheticConnector) getBatchExecutionData(workCtx context.Context) (BatchResponseBody, error) {
	path := getSyntheticBatchPath(sc.executionData.BatchId)
	resp, err := sc.dtClient.Get(workCtx, path)
	if err != nil {
		return BatchResponseBody{}, err
	}

	log.Debug(string(resp))

	batchResponseBody := BatchResponseBody{}
	err = json.Unmarshal(resp, &batchResponseBody)
	if err != nil {
		log.Error(err.Error())
		return BatchResponseBody{}, err
	}

	return batchResponseBody, nil
}

// Calculates synthetic execution success rate
// triggeredCount = executedCount + failedToExecuteCount
// executedCount = failedCount + executions finished with SUCCESS
//
// -> triggeredCount = failedToExecuteCount + failedCount + executions finished with SUCCESS
//
func calculateSuccessRate(batchResponseBody BatchResponseBody) (float64, error) {
	if batchResponseBody.TriggeredCount < 1 {
		return 0, nil
	}

	totalFailedCount := batchResponseBody.FailedCount + batchResponseBody.FailedToExecuteCount
	totalTriggeredCount := batchResponseBody.TriggeredCount

	failureRate := float64(totalFailedCount) / float64(totalTriggeredCount) * 100
	successRate := 100.0 - failureRate

	successRateRounded := math.Round(successRate*100) / 100

	return successRateRounded, nil
}

// Waits for the last triggered batch to return an SUCCES or FAILED status.
//
// Attention: Does not wait for data retrieval
//
func (sc *SyntheticConnector) WaitForBatchExecution(workCtx context.Context) (BatchResponseBody, float64, error) {
	pollingStartTime := time.Now().UTC()
	pollingTimeout, err := time.ParseDuration("300s")
	if err != nil {
		return BatchResponseBody{}, 0, err
	}

	requestCounter := 1
	pollingInterval := 10 * time.Second

	for {
		currentPollingTime := time.Now().UTC()
		if currentPollingTime.Sub(pollingStartTime).Seconds() > pollingTimeout.Seconds() {
			return BatchResponseBody{}, 0, fmt.Errorf("could not retrieve data within %f seconds", pollingTimeout.Seconds())
		}

		log.Debug("Requesting data (", requestCounter, ")")

		batchResponseBody, err := sc.getBatchExecutionData(workCtx)
		if err != nil {
			return BatchResponseBody{}, 0, err
		}

		// batchStatus = RUNNING || SUCCESS || FAILED
		if batchResponseBody.BatchStatus == "SUCCESS" || batchResponseBody.BatchStatus == "FAILED" {
			successRate, _ := calculateSuccessRate(batchResponseBody)
			return batchResponseBody, successRate, nil
		}

		log.Debug("Waiting ", pollingInterval.Seconds(), " seconds...")
		time.Sleep(pollingInterval)

		requestCounter++
	}
}

// TBD
// func (sc *SyntheticConnector) WaitForBatchData(workCtx context.Context) error {}

func NewSyntheticConnector(dtClient dynatrace.ClientInterface) *SyntheticConnector {
	return &SyntheticConnector{
		dtClient: dtClient,
	}
}
