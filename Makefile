
Name=webbook
docker:
	@rm -f webbook|| true
	@go mod tidy
	@git add .
	@git commit -m "makefile"
	@git push origin main
	@go build  $(Name) .
	@docker rmi -f  $(Name):v1
	@docker build  -t  $(Name):v1 .
	@docker rm -f $(Name)
	@docker run -d --name  $(Name) -network=webbook_default $(Name):v1
	@docker logs -f  $(Name)
build:
	@docker compose up -d