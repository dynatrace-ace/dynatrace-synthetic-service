.PHONY: checktag build push deploy helmchart

IMAGE := "dynatraceace/keptn-dt-synthetic-service"

build: checktag
	@docker build -t "${IMAGE}:${tag}" .
	@echo "\nSuccesfully built ${IMAGE}:${tag}!"

push: build
	@docker push "${IMAGE}:${tag}"
	@echo "\nSuccesfully pushed ${IMAGE}:${tag}!"

deploy: checktag
	@helm upgrade --install -n keptn dynatrace-synthetic-service --set "dynatraceService.image.tag=${tag}" chart/

helmchart: checktag
	@tar -czf dt-synthetic-service-${tag}.tar.gz chart/
	@echo "\nSuccesfully bundled chart dt-synthetic-service-${tag}.tar.gz!" 

checktag:
ifndef tag
	$(error tag is not set)
endif
