.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f haha9705/webook:v0.0.1
	@docker build -t haha9705/webook:v0.0.1 .