
local-run: 
	go run cmd/main.go 

docker-build:
	docker build -t url-shortener . 

docker-run: docker-build
	docker run -p 8080:8080 url-shortener

