[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/arloliu/jsonpack)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/arloliu/jsonpack/main/LICENSE)

Fast and space efficiency JSON serialization golang library. It is a schema oriented design which leverages schema definition to encode JSON document into compact binary encoded format, and decodes back into JSON document.

# Introduction
When we want to exchange data between services or over network, the JSON is most popular format to do it.
In golang world, the most convenient way is using official `encoding/json` package to marshal and unmarshal JSON document from/to struct or map. For usual RESTful web service scenario, JSON format is quite convenience and representative, but for real-time message exchanging or other scenarios that has small footprint data and low-letency requirement, JSON is a bit too heavyweight, not only data footprint is not space saving, but also has heavy loading in encode/decode procedure.

So if we want to a compact, small footprint data for exchanging over network, and also leverages the 
convenience of JSON, we need an encoding format that removes "property name" and other notations likes ':', '[', '{'...etc from original JSON document, and leaves "value" only.

To achieve this goal, we need a schematic to define our JSON document and provide enough information for serialization engine to know sequence and data type of every properties in document. it's the reason why `JSONPack` is a *schema oriented design* library.

# Key Features
* Similar Marshal / Unmarshal API to standard `encoding/json` package.
* Space saving encoded format, the size of encoded data is similar to Protocol Buffers, can be 30-80% of original JSON document, depends on data.
* Blazing fast, provides about 3.x decoding speed compared to `protobuf` and many times than other JSON packages.
* Memory saving design, avoids any un-neccessary memory allocations, suitable for embedded environment.
* Has production ready javascript implementation [Buffer Plus](https://github.com/arloliu/buffer-plus), can be used in node.js and Web browser environment.
* No need to write schema definition by hand, `jsonpack` will generate schema definition from golang struct automatically.

# How to Get
```
go get https://github.com/arloliu/jsonpack
```
# Usage

Example of add schema definition
```go
import "github.com/arloliu/jsonpack"

type Info struct {
	Name string `json:"name"`
	Area uint32 `json:"area"`
	// omit this field
	ExcludeField string `-`
}

jsonPack := jsonpack.NewJSONPack()
sch, err := jsonPack.AddSchema("Info", Info{}, jsonpack.LittleEndian)
```

Example of encoding data with `Info` struct
```go
infoStruct := map[string]interface{} {
	"name": "example name",
	"area": uint32(888),
}
// encodedResult1 contains encoded data,
encodedResult1, err := jsonPack.Marshal("Info", infoStruct)
```

Example of encoding data with golang map
```go
infoMap := map[string]interface{} {
	"name": "example name",
	"area": uint32(888),
}

encodedResult2, err := jsonPack.Marshal("Info", infoMap)
```

Example of decoding data
```go
decodeInfoStruct = Info{}
err := jsonPack.Decode("Info", encodedResult1, &decodeInfoStruct)

decodeInfoMap = make(map[string]interface{})
err := jsonPack.Decode("Info", encodedResult2, &decodeInfoMap)
```

# Benchmark
> The benchmark result is a important reference but not always suitable for every scenarios.

Test environment: Intel i7-9700K CPU@3.60GHz.

Benchmark code is [here](https://github.com/arloliu/jsonpack/blob/main/benchmark/benchmark_test.go)

*Sorts from fastest to slowest in the following.*
## Encode from golang map
|           | ns/op       | allocation bytes | allocation times |
|-----------|-------------|------------------|------------------|
| jsonpack  | 1933 ns/op  | 752 B/op         | 2 allocs/op      |
| jsoniter  | 10134 ns/op | 3320 B/op        | 46 allocs/op     |
| std. json | 23560 ns/op | 8610 B/op        | 171 allocs/op    |
| goccy     | 75298 ns/op | 82639 B/op       | 651 allocs/op    |

## Decode into golang map
|           | ns/op       | allocation bytes | allocation times |
|-----------|-------------|------------------|------------------|
| jsonpack  | 6461 ns/op  | 6512 B/op        | 96 allocs/op     |
| jsoniter  | 17436 ns/op | 9666 B/op        | 290 allocs/op    |
| std. json | 18949 ns/op | 8864 B/op        | 228 allocs/op    |
| goccy     | 19985 ns/op | 15900 B/op       | 316 allocs/op    |


## Encode from golang struct
|           | ns/op      | allocation bytes | allocation times |
|-----------|------------|------------------|------------------|
| jsonpack  | 1834 ns/op | 800 B/op         | 3 allocs/op      |
| protobuf  | 1972 ns/op | 896 B/op         | 1 allocs/op      |
| goccy     | 2166 ns/op | 1280 B/op        | 1 allocs/op      |
| jsoniter  | 3372 ns/op | 1296 B/op        | 3 allocs/op      |
| std. json | 3578 ns/op | 1280 B/op        | 1 allocs/op      |

## Decode into golang struct
|           | ns/op       | allocation bytes | allocation times |
|-----------|-------------|------------------|------------------|
| jsonpack  | 1475 ns/op  | 96 B/op          | 2 allocs/op      |
| goccy     | 3284 ns/op  | 2215 B/op        | 5 allocs/op      |
| jsoniter  | 4680 ns/op  | 1072 B/op        | 79 allocs/op     |
| protobuf  | 5075 ns/op  | 3152 B/op        | 84 allocs/op     |
| std. json | 18378 ns/op | 1232 B/op        | 69 allocs/op     |

The benchmark result indicates jsonpack keeps constant performance on both encoding and encoding side, and keeps very low memory allocation size and times.

The benchmark result also delivers an important message.

**The performance of operating with golang map sucks**

So it's better to use struct if possible. :)


