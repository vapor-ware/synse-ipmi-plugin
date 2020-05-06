
FROM vaporio/foundation:bionic

LABEL org.label-schema.schema-version="1.0" \
      org.label-schema.name="vaporio/ipmi-plugin" \
      org.label-schema.vcs-url="https://github.com/vapor-ware/synse-ipmi-plugin" \
      org.label-schema.vendor="Vapor IO"

RUN apt-get update \
 && apt-get install -y --no-install-recommends ipmitool ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY synse-ipmi-plugin ./plugin

EXPOSE 5001
ENTRYPOINT ["./plugin"]
