FROM golang:1.16-alpine AS builder
ARG BIN
ARG VERSION
LABEL stage=builder

WORKDIR /src
COPY . .
RUN apk update && apk add make && make build


FROM alpine
ARG BIN
ARG VERSION
LABEL maintainer="jdologl@cslab.ece.ntua.gr" version=$VERSION

COPY --from=builder /src/$BIN /bin/$BIN
ENV BIN $BIN

CMD ["/bin/$BIN"]
