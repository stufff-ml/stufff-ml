# stufff-ml
stufff-ml


admin:full	- Allows full read/write admin access to the API.
api:full		- Allows full read/write access to the API with all permissions to administrate the app


https://github.com/GoogleCloudPlatform/python-docs-samples


https://godoc.org/cloud.google.com/go/storage

gcloud functions deploy func_submit --region=europe-west1 --trigger-http --entry-point=handle_request --memory=128MB --runtime=python37 --source=func_submit 

training_input = {
    'scaleTier': 'BASIC',
    'packageUris': ['gs://models.stufff.review/packages/default-1/default-1.tar.gz'],
    'pythonModule': 'model.task',
    'args': [
      '--model-id', '26144595808e',
      '--model-rev','1'
    ],
    'region': 'europe-west1',
    "jobDir": 'gs://models.stufff.review/26144595808e/26144595808e/default-1',
    'runtimeVersion': '1.12',
    'pythonVersion': '2.7'
  }
	
{
	"projectId":"stufff-review",
	"jobId":"M26144595808e_default_5",
	"scaleTier": "BASIC",
    "packageUris": ["gs://models.stufff.review/packages/default-1/default-1.tar.gz"],
    "pythonModule": "model.task",
    "args": [
      "--model-id", "26144595808e",
      "--model-rev","1"
    ],
    "region": "europe-west1",
    "jobDir": "gs://models.stufff.review/26144595808e/26144595808e/default-1",
    "runtimeVersion": "1.12",
    "pythonVersion": "2.7"
}

[
	{
		"event": "sign-up",
		"entity_type": "user",
		"entity_id": "1",
		"properties": ["a","a"]
		},
		{
		"event": "sign-up",
		"entity_type": "user",
		"entity_id": "2",
		"properties": ["b","b"]
	}
]




[
	{
		"entity_id": "6",
		"domain": "buy",
		"items": [
			{
				"item":"316",
				"score":0.3
			},
			{
				"item":"318",
				"score":0.14
			}
		]
	},
	{
		"entity_id": "457",
		"domain": "buy",
		"items": [
			{
				"item":"316",
				"score":0.37
			},
			{
				"item":"318",
				"score":0.22
			},
			{
				"item":"527",
				"score":0.2
			}
		]
	}
]
