FROM golang:1.8-jessie
MAINTAINER <ApexDNS apex@apebits.com>

RUN mkdir -p /apexd
ADD cmd/apexd/statistics.toml /apexd
ADD apexd /apexd/

WORKDIR /apexd

CMD ./apexd
