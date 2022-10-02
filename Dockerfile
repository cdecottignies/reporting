FROM scratch
LABEL maintainer="eng@vadesecure.com"
ARG VERSION=dev
LABEL org.opencontainers.image.version=$VERSION
ARG REPORTING_BIN=build/cmd/reporting/reporting
ADD ./conf/ca-certificates.crt /etc/ssl/certs/
COPY $REPORTING_BIN  /reporting
COPY conf /conf

ENTRYPOINT ["/reporting"]
