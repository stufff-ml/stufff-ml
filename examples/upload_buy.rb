require 'net/http'
require 'uri'
require 'json'
require 'date'
require 'csv'

# ruby upload_buy.rb http://localhost:8080/api/1/events data_small.csv xoxo-...
# ruby upload_buy.rb http://stufff-review.appspot.com/api/1/events data_small.csv xoxo-...

endpoint = ARGV[0]
filename = ARGV[1]
token = ARGV[2]
batch_size = 200
use_timestamp = false

# prepare the connection
uri = URI.parse endpoint

http = Net::HTTP.new(uri.host, uri.port)
http.use_ssl = false

req = Net::HTTP::Post.new(uri.path, {'Content-Type' =>'application/json',  'Authorization' => "Bearer #{token}"})

payload = []
n = 0

CSV.foreach(filename) do |row|
  if row[2] == "5.0"
    payload << {
      "event" => "buy",
      "entity_type" => "user",
      "entity_id" => row[0],
      "target_entity_type" => "item",
      "target_entity_id" => row[1],
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
end

# send the remaining records
if payload.size > 0
  req.body =  payload.to_json
      
  start = DateTime.now.strftime('%Q').to_i  
  res = http.request(req)
  stop = DateTime.now.strftime('%Q').to_i

  puts "Code: #{res.code} - #{stop - start} ms."

end