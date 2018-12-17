from google.cloud import storage
from io import BytesIO
import pandas as pd
import turicreate as tc

# constants but configurable
source_bucket_name = 'models.stufff.review'
target_bucket_name = 'models.stufff.review'
training_percentage = 0.8
training_seed = 5
max_recommendations = 10

# should come from ENV, job specific
model_id = 'foo1233.default'
rev = '2'

# calculate based on ENV
source_name = model_id + '/' + model_id + '_' + rev + '.csv'
target_name = model_id + '/' + model_id + '_model_' + rev + '.csv'

# access to Cloud Storage
client = storage.Client()
source_bucket = client.get_bucket(source_bucket_name)

# Load the training data as a panda dataframe
source_blob = source_bucket.get_blob(source_name)
content = source_blob.download_as_string()
_data = pd.read_csv(BytesIO(content))
_data.columns = ['event','entity_type','entity_id','target_entity_type','target_entity_id','timestamp','properties']

# create the model
data = tc.SFrame( _data, format='dataframe')
data_train, data_test = data.random_split(training_percentage, seed=training_seed)
model = tc.recommender.create(data_train, 'entity_id', 'target_entity_id')

# get a recommendation for all entities
entities = _data.entity_id.unique()
rec = model.recommend(entities, max_recommendations)

# prepare to upload the result
target_bucket = client.get_bucket(target_bucket_name)
target_blob = target_bucket.blob(target_name)

rec.export_csv('upload.csv',',',header=False)
target_blob.upload_from_filename('upload.csv')
