from google.cloud import storage

client = storage.Client()

bucket = client.get_bucket('exports.stufff.review')


# Then do other things...
blob = bucket.get_blob('foo1233.default/foo1233.default_1544791560.csv')
print(blob.download_as_string())
