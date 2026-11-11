BIN       := quantile
DEMO_FILE := $(CURDIR)/demo.jsonl
LOAD_FILE := $(CURDIR)/1M.jsonl
LOAD_N    ?= 1000000

.PHONY: all setup build test demo load-gen load-run clean
.DEFAULT_GOAL := demo

all: demo

setup:
	brew install go@1.27
	@echo 'Add to PATH: export PATH="$$(brew --prefix go@1.27)/bin:$$PATH"'

build:
	go build -o $(BIN) .

test:
	go test -race ./...

demo: build
	@echo $(DEMO_FILE) | ./$(BIN) -q 0.75 -f latency_ms

$(LOAD_FILE):
	awk -v n=$(LOAD_N) 'BEGIN{ srand(42); \
	  split("acme globex",c); split("us-east us-west",r); split("/api/v1/query /api/v1/insert",e); \
	  for(i=1;i<=n;i++) printf "{\"timestamp\": \"2026-06-18T10:%02d:%02dZ\", \"customer_id\": \"%s\", \"region\": \"%s\", \"endpoint\": \"%s\", \"latency_ms\": %d, \"bytes_sent\": %d}\n", \
	    int(i/60)%60, i%60, c[int(rand()*2)+1], r[int(rand()*2)+1], e[int(rand()*2)+1], int(rand()*10)+1, 5000+int(rand()*5000) }' > $@

load-gen: $(LOAD_FILE)

load-run: build $(LOAD_FILE)
	echo $(LOAD_FILE) | /usr/bin/time -l ./$(BIN) -q 0.75 -f bytes_sent -v

clean:
	rm -f $(BIN) $(LOAD_FILE)
