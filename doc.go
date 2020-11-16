/*
Fast and space efficiency JSON serialization golang library.
It is a schema oriented design which leverages schema definition to encode
JSON document into compact binary encoded format, and decodes back into JSON document.

Introduction

When we want to exchange data between services or over network, the JSON is most popular
format to do it.
In golang world, the most convenient way is using official `encoding/json` package to
marshal and unmarshal JSON document from/to struct or map. For usual RESTful web service
scenario, JSON format is quite convenience and representative, but for real-time message
exchanging or other scenarios that has small footprint data and low-letency requirement,
JSON is a bit too heavyweight, not only data footprint is not space saving, but also has
heavy loading in encode/decode procedure.

So if we want to a compact, small footprint data for exchanging over network, and also
leverages the convenience of JSON, we need an encoding format that removes "property name"
and other notations likes ':', '[', '{'...etc from original JSON document, and leaves
"value" only.

To achieve this goal, we need a schematic to define our JSON document and provide enough
information for serialization engine to know sequence and data type of every properties
in document. it's the reason why jsonpack is a schema oriented design library.

Key Features

* Similar Marshal / Unmarshal API to standard `encoding/json` package.

* Space saving encoded format, the size of encoded data is similar to Protocol Buffers, can be 30-80% of original JSON document, depends on data.

* Blazing fast, provides about 3.x decoding speed compared to `protobuf` and many times than other JSON packages.

* Memory saving design, avoids any un-neccessary memory allocations, suitable for embedded environment.

* Has production ready javascript implementation https://github.com/arloliu/buffer-plus, can be used in node.js and Web browser environment.

* No need to write schema definition by hand, `jsonpack` will generate schema definition from golang struct automatically.
*/
package jsonpack
