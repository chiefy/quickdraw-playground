
csvurl := https://data.ny.gov/api/views/7sqk-ycpk/rows.csv?accessType=DOWNLOAD&bom=true&format=true&sorting=true

vendor:
	@dep ensure

.PHONY: csv
csv: quick-draw.csv

.PHONY: clean
clean:
	@-rm -rf vendor quick-draw.csv

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

.PHONY: makerange
makerange:
	@/bin/bash -c 'for n in {1..80}; do printf %i, $$n; done'