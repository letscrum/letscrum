FROM --platform=$BUILDPLATFORM golang:1.22 as builder

ENV GO111MODULE=on

WORKDIR /app

COPY . .

# install protoc
RUN apt-get update && apt-get install -y protobuf-compiler

RUN make api_clean
RUN make api_dep_install
RUN make api_gen
RUN make build

FROM alpine:3.20.1

ARG TARGETARCH

COPY --from=builder /app/bin/letscrum /bin/
COPY --from=builder /app/config/config.yaml /etc/letscrum/config.yaml

CMD ["/bin/letscrum", "server", "--config", "/bin/config/config.yaml"]

EXPOSE 8081