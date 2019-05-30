#
# Builder Image
#
FROM vaporio/golang:1.11 as builder

#
# Final Image
#
FROM scratch

LABEL org.label-schema.schema-version="1.0" \
      org.label-schema.name="vaporio/ipmi-plugin" \
      org.label-schema.vcs-url="https://github.com/vapor-ware/synse-ipmi-plugin" \
      org.label-schema.vendor="Vapor IO"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy the executable.
COPY synse-ipmi-plugin ./plugin

EXPOSE 5001
ENTRYPOINT ["./plugin"]
