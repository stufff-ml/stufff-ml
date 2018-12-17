require 'net/http'
require 'uri'
require 'json'
require 'date'
require 'csv'

# ruby upload_rated.rb http://localhost:8080/api/1/events xoxo-ffffffff test.csv
# ruby upload_rated.rb http://stufff-review.appspot.com/api/1/events xoxo-ffffffff test.csv
# ruby upload_rated.rb http://stufff-review.appspot.com/api/1/events xoxo-foobarxx test.csv

endpoint = ARGV[0]
token = ARGV[1]
filename = ARGV[2]
batch_size = 100
use_timestamp = false

# prepare the connection
uri = URI.parse endpoint

http = Net::HTTP.new(uri.host, uri.port)
http.use_ssl = false

req = Net::HTTP::Post.new(uri.path, {'Content-Type' =>'application/json',  'Authorization' => "Bearer #{token}"})

payload = []
n = 0

CSV.foreach(filename) do |row|
  payload << {
    "event" => "rated",
    "entity_type" => "user",
    "entity_id" => row[0],
    "target_entity_type" => "item",
    "target_entity_id" => row[1],
    "properties" => [ row[2]],
    "timestamp" => use_timestamp ? row[3].to_i : Time.now.getutc.to_i
  }
  n = n + 1

  if n == batch_size
    req.body =  payload.to_json
  
    start = DateTime.now.strftime('%Q').to_i  
    res = http.request(req)
    stop = DateTime.now.strftime('%Q').to_i

    puts "Code: #{res.code} - #{stop - start} ms."

    # reset
    n = 0
    payload = []
  end

end

# send the remaining records
if payload.size > 0
  req.body =  payload.to_json
      
  start = DateTime.now.strftime('%Q').to_i  
  res = http.request(req)
  stop = DateTime.now.strftime('%Q').to_i

  puts "Code: #{res.code} - #{stop - start} ms."

end