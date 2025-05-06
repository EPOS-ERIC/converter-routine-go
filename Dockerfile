# Using this image to ensure compatibility with the converter-service
# can probably be better
FROM alpine:3.20
RUN apk --no-cache add python3

LABEL authors="valeriovinciarelli"

WORKDIR /opt/converter

COPY converter-routine converter-routine

CMD ["./converter-routine"]
