require 'net/http'
require 'uri'
require 'json'
require 'date'
require 'csv'

# ruby upload_rated.rb http://localhost:8080/api/1/events xoxo-ffffffff test.csv
# ruby upload_rated.rb http://stufff-review.appspot.com/api/1/events xoxo-ffffffff test.csv

endpoint = ARGV[0]
token = ARGV[1]
filename = ARGV[2]

# prepare the connection
uri = URI.parse endpoint

http = Net::HTTP.new(uri.host, uri.port)
http.use_ssl = false

req = Net::HTTP::Post.new(uri.path, {'Content-Type' =>'application/json',  'Authorization' => "Bearer #{token}"})

CSV.foreach(filename) do |row|
  req.body =  {
    "event" => "rated",
    "entity_type" => "user",
    "entity_id" => row[0],
    "target_entity_type" => "item",
    "target_entity_id" => row[1],
    "properties" => [ row[2]],
    "timestamp" => row[3].to_i
  }.to_json
  
  start = DateTime.now.strftime('%Q').to_i  
  res = http.request(req)
  stop = DateTime.now.strftime('%Q').to_i

  puts "Code: #{res.code} - #{stop - start} ms."

end
