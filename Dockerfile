FROM iron/go:dev as builder
WORKDIR /go/src/github.com/vapor-ware/synse-ipmi-plugin
COPY . .
RUN make build


FROM iron/go
LABEL maintainer="vapor@vapor.io"

WORKDIR /plugin

RUN apk --update --no-cache add ipmitool

COPY --from=builder /go/src/github.com/vapor-ware/synse-ipmi-plugin/build/plugin ./plugin
COPY config.yml .
COPY config/proto /etc/synse/plugin/config/proto

EXPOSE 5001

CMD ["./plugin"]