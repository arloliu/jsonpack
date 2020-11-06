#jsonpack
Fast and space efficiency JSON serialization golang library. It is a schema oriented implementation which leverages schema definition to encode JSON format text data into compact binary encoded data, and then decodes back into JSON format.

# Introduction
When we want to exchange data between services or over network, the JSON is most popular format to do it.
In golang world, the most convenient way is using official `encoding/json` package to marshal and unmarshal JSON data from/to struct or map.
