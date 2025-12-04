# Using this image to ensure compatibility with the converter-service
# can probably be better
FROM alpine:3.23

RUN apk --no-cache add python3 ca-certificates

LABEL authors="valeriovinciarelli"

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /opt/converter

COPY converter-routine converter-routine

RUN mkdir /opt/converter/plugins

RUN chown -R appuser:appgroup /opt/converter

USER appuser:appgroup

CMD ["./converter-routine"]
