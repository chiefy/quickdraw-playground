csvurl := https://data.ny.gov/api/views/7sqk-ycpk/rows.csv?accessType=DOWNLOAD&bom=true&format=true&sorting=true


$(GOPATH)/bin/dep:
	@curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

.PHONY: all
all: clean csv vendor build-ui 
	@docker-compose rm -fv 
	@docker-compose build 
	@HOST=$(HOST) docker-compose up

.PHONY: stop 
stop:
	@docker-compose stop; docker-compose rm -fv

vendor: $(GOPATH)/bin/dep
	@dep ensure

.PHONY: csv
csv: quick-draw.csv

.PHONY: clean
clean:
	@-rm -rf vendor quick-draw.csv ui/dist/spa

quick-draw.csv: 
	@curl -LSs "$(csvurl)" | awk -F, '{seen[$$1,$$2]++;seen[$$2,$$1]++}seen[$$1,$$2]==1 && seen[$$2,$$1]==1' > $@

.PHONY: start-postgres
start-postgres:
	@docker-compose up -d

.PHONY: import-csv 
import-csv: vendor csv 
	@go run . --import

.PHONY: serve
serve:
	@go run . --serve

.PHONY: ui 
ui:
	@cd ui && quasar dev

.PHONY: build-ui
build-ui: ui/dist/spa 

ui/dist/spa:
	@cd ui && quasar build 

.PHONY: makerange
makerange:
	@/bin/bash -c 'for n in {1..80}; do printf %i, $$n; done'