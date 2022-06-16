# Dynatrace-synthetic-service

## This project was blatantly copied from [dynatrace-service (0.22.0)](https://github.com/keptn-contrib/dynatrace-service).

## Overview

The dynatrace-synthetic-service allows you to trigger Dynatrace Synthetic executions as part of a Keptn sequence. This way Synthetic HTTP or clickpath test results can be integrated in a quality gate.

## Installation

Apart from using the dynatrace-synthetic-service repository and Helm chart, please follow the [dynatrace-service installation guides](https://github.com/keptn-contrib/dynatrace-service/blob/0.22.0/documentation/installation.md).

E.g. chart installation:
```
helm upgrade --install dynatrace-synthetic-service -n keptn chart/ \
  --set dynatraceService.config.keptnApiUrl=$KEPTN_ENDPOINT \
  --set dynatraceService.config.keptnBridgeUrl=$KEPTN_BRIDGE_URL
```

## Project setup

Is step is ony required if it hasn't been completed as part of the *dynatrace-service* project setup!

A project can be setup according to the [dynatrace-service project setup guides](https://github.com/keptn-contrib/dynatrace-service/blob/0.22.0/documentation/project-setup.md).

## Triggering a synthetic test

A synthetic test can be triggered by publishing an event with the following attributes:

```
{
  "type": "sh.keptn.event.dev.test.triggered",
  "contenttype": "application/json",
  "specversion": "0.2",
  "source": "<event source>",
  "data": {
    "project": "<Keptn Project>",
    "service": "<Keptn Service>",
    "stage": "<Keptn Stage>",
    "monitorTag": "<Synthetic Monitor tag>", # Either monitorTag or monitorId is required
	  "monitorId": "<Synthetic Monitor id>",   # Either monitorTag or monitorId is required
	  "waitFor": "EXECUTION"                   # Optional
  }
}
```

|Attribute|Comment|
|---|---|
|monitorTag|Service triggers execution of all Synthetic Monitors tagged with the value of *monitorTag*. Either monitorTag or monitorId has to be specified.|
|monitorId|Service triggers execution of the particular Synthetic Monitor which id matches *monitorId*. Either monitorTag or monitorId has to be specified|
|waitFor|Optional: By default, a synthetic test is triggered without waiting for any results. The attribute can be set to "EXECUTION" which makes the serice wait for synthetic execution results, i.e. successful/failed|
