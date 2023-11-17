FROM centos
COPY webbook /
COPY config/config.yaml /config/
NetWork webbook_default
CMD  ./webbook
