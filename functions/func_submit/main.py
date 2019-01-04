from flask import jsonify, request

from googleapiclient import discovery
from googleapiclient import errors

def handle_request(request):
  payload = request.json
  if payload is None:
      return jsonify({'status':'error'})
      
  job_id = payload.get('jobId')
  project_id = payload.get('projectId')
  
  # remove the following attributes; they are not part of the ML API call
  payload.pop('jobId')
  payload.pop('projectId')

  resp = submit_job(project_id, job_id, payload)
  return jsonify(resp)


def submit_job(project_id, job_id, training_input):

  # Store your full project ID in a variable in the format the API needs.
  parent_project = 'projects/{}'.format(project_id)

  # Build a representation of the Cloud ML API.
  ml = discovery.build('ml', 'v1')

  # Create a dictionary with the fields from the request body.
  request_dict = {
    'jobId': job_id,
    'trainingInput': training_input
  }

  # Create a request to submit a model for training
  request = ml.projects().jobs().create(parent=parent_project, body=request_dict)

  # Make the call.
  resp = request.execute()
  return resp
  
