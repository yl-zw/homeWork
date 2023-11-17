FROM centos
COPY webbook /
COPY config/config.yaml /config/
EXPOSE 8080/http
CMD  ./webbook
