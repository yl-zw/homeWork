FROM centos
COPY webbook /
COPY config/config.yaml /config/
CMD  ./webbook
