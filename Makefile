swagger:
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models

test: 
	go test -v ./cmd/web

cover: 
	go test -cover ./cmd/web

run: 
	go run ./cmd/web

build: 
	go build -o xvi-wiek ./cmd/web

ebook:
	./build_ebook.sh
	
view:
	evince ./ui/static/pdf/xvi-wiek.pdf 
