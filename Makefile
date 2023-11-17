.PHONY:docker
docker:
	@rm -f webbook|| true
	@go mod tidy
	@git add .
	@git commit -m "makefile"
	@git push origin main
	@go build  webbook .
	@docker rmi -f webbook:v1
	@docker build  -t webbook:v1 .
	@docker run -d --name webbook webbook:v1
