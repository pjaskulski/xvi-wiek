swagger:
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models

test: 
	go test -v ./cmd/web

cover: 
	go test -cover ./cmd/web

coverage:
	go test -race -coverprofile=coverage.out ./cmd/web 
	go tool cover -html coverage.out -o coverage.html

cyclo:
	gocyclo -over 15 ./cmd/web

vet:
	go vet ./cmd/web

gosec:
	gosec -quiet -exclude G104 ./cmd/web

unparam:
	unparam ./cmd/web

run: 
	go run ./cmd/web

build: 
	go build -o xvi-wiek ./cmd/web

ebook:
	./build_ebook.sh
	
view:
	evince ./ui/static/pdf/xvi-wiek.pdf 
