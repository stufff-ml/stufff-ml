from flask import request, jsonify

from oauth2client.client import GoogleCredentials
from googleapiclient import discovery
from googleapiclient.discovery_cache.base import Cache

# In order to fix some weirdness, see https://github.com/googleapis/google-api-python-client/issues/325
class MemoryCache(Cache):
    _CACHE = {}

    def get(self, url):
        return MemoryCache._CACHE.get(url)

    def set(self, url, content):
        MemoryCache._CACHE[url] = content


def handle_request(request):
  from flask import abort
  from google.cloud import error_reporting
  
  client = error_reporting.Client()

  if request.method == 'GET':
    return abort(405)
  elif request.method == 'PUT':
        return abort(405)

  j = request.json
  if j is None:
      return make_response({'status':'error','message':'Missing payload'}, 400)

  project_id = 'stufff-review'
  job_id = j.get('id')
  training_input = j.get('payload')

  # Store your full project ID in a variable in the format the API needs.
  parent_project = 'projects/{}'.format(project_id)

  # Build a representation of the Cloud ML API.
  ml = discovery.build('ml', 'v1', cache=MemoryCache())

  # Create a dictionary with the fields from the request body.
  request_dict = {
    'jobId': job_id,
    'trainingInput': training_input
  }

  # Create a request to submit a model for training
  request = ml.projects().jobs().create(parent=parent_project, body=request_dict)

  try:
    resp = request.execute()
    return jsonify({'status':'ok'})
  except:
    client.report_exception()
    return jsonify({'status':'error'}), 500
  