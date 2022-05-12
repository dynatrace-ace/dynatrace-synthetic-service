# Release Notes 0.7.0

## New Features

- Enrich Events sent to Dynatrace with Labels & Source (#116)[https://github.com/keptn-contrib/dynatrace-service/issues/116]
- Support multiple Dynatrace environments (for different stages/projects)
- Support non-keptn deployed services (#115)[https://github.com/keptn-contrib/dynatrace-service/issues/115]

## Fixed Issues

- Only handle cloud-events coming from the Dynatrace tenant (#127)[https://github.com/keptn-contrib/dynatrace-service/issues/127]
- Set Dynatrace OneAgent Operator version to 0.6.0 to avoid incompatibilities with older Kubernetes versions (#125)[https://github.com/keptn-contrib/dynatrace-service/issues/125]

## Known Limitations
- Dynatrace Kubernetes OneAgent operator is now limited to version 0.6.0, which is the last version that supports Kubernetes 1.13 (#132)[https://github.com/keptn-contrib/dynatrace-service/issues/132)]
- Prior to 0.7.0, Dynatrace-Service has created an alerting profile with a filter that blocks certain problem notifications. There is a known work-around for this described in (#125)[https://github.com/keptn-contrib/dynatrace-service/issues/125]
- Alerting profiles are not overwritten (#134)[https://github.com/keptn-contrib/dynatrace-service/issues/134] - also related to (#125)[https://github.com/keptn-contrib/dynatrace-service/issues/125]
- When using Container-Optimized OS (COS) based GKE clusters, the deployed OneAgent has to be updated after the installation of Dynatrace
