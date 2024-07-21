FROM --platform=$BUILDPLATFORM golang:1.22 as builder

ENV GO111MODULE=on

WORKDIR /app

COPY * /app/

RUN make api_clean
RUN make api_dep_install
RUN make api_gen
RUN make build

FROM alpine:3.20.1

COPY --from=builder /bin/letscrum /bin/
COPY --from=builder /config /bin/config

EXPOSE 8081