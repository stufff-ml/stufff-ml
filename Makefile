
.PHONY: all
all: deploy_functions deploy_app

.PHONY: deploy_functions
deploy_functions: train_model

.PHONY: train_model
train_model:
	cd functions && gcloud functions deploy train_model --region=europe-west1 --trigger-http --entry-point=handle_request --memory=128MB --runtime=python37 --source=train_model --quiet

.PHONY: deploy_app
deploy_app:
	cd cmd/app && gcloud app deploy . --quiet