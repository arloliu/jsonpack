DEBUG ?= 0
ifeq ($(DEBUG), 1)
	BUILDFLAG = -tags jsonpack_debug
	TESTFLAG = -tags jsonpack_debug
endif

BENCH_TARGET = $(TARGET)
ifeq ($(TARGET),)
TARGET = .
endif
PKGNAME = github.com/arloliu/jsonpack
.PHONY: test benchmark cpu-pprof mem-pprof protoc
test:
	go test $(PKGNAME)/internal/buffer -v $(TESTFLAG) -run $(TARGET)
	go test $(PKGNAME) -v $(TESTFLAG) -run $(TARGET)


benchmark:
	cd ./benchmark; go test benchmark -v $(TESTFLAG) -tags jsonpack_benchmark -benchmem -run ^$$ -bench $(TARGET)

cpu-pprof:
	$(eval LOGFILE := $(shell mktemp -u --suffix .prof))
	# go test $(PKGNAME) -cpuprofile $(LOGFILE) -v -benchmem -benchtime 5s -run ^$$ -bench BenchmarkComplex_Jsoniter_Marshal
	
	cd ./benchmark; go test benchmark -cpuprofile $(LOGFILE) -v -benchmem -benchtime 2s -run ^$$ -bench $(BENCH_TARGET)
	pprof -http=0.0.0.0:4231 $(LOGFILE)
mem-pprof:
	$(eval LOGFILE := $(shell mktemp -u --suffix .prof))
	# go test $(PKGNAME) -cpuprofile $(LOGFILE) -v -benchmem -benchtime 5s -run ^$$ -bench BenchmarkComplex_Jsoniter_Marshal
	
	cd ./benchmark; go test benchmark -memprofile $(LOGFILE) -v -benchmem -benchtime 2s -run ^$$ -bench $(BENCH_TARGET)
	pprof -http=0.0.0.0:4231 $(LOGFILE)

protoc:
	protoc --go_out=paths=source_relative:. testdata/*.proto