
Name=webbook
build:
	@rm -f webbook|| true
	@go build  $(Name) .
	@docker rmi -f  $(Name):v1
	@docker build  -t  $(Name):v1 .
	@docker rm -f $(Name)
	@docker run -d --name  $(Name) --network webbook_default $(Name):v1
	@docker logs -f  $(Name)
git:
	@go mod tidy
	@git add .
	@git commit -m "makefile"
	@git pull
	@git push origin main