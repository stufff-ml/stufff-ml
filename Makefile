
.PHONY: all
all: deploy_functions deploy_app

.PHONY: deploy_functions
deploy_functions: deploy_submit

.PHONY: deploy_submit
deploy_submit:
	cd functions && gcloud functions deploy func_submit --region=europe-west1 --trigger-http --entry-point=handle_request --memory=128MB --runtime=python37 --source=func_submit --quiet

.PHONY: deploy_app
deploy_app:
	cd cmd/app && gcloud app deploy . --quiet