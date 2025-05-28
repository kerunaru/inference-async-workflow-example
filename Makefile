build:
	docker build --target server -t local/inference-async-workflow-example:server -f Containerfile .
	docker build --target queue -t local/inference-async-workflow-example:queue -f Containerfile .
