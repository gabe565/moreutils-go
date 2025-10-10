FROM alpine AS source
WORKDIR /app
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/moreutils .
RUN ./moreutils install .

FROM alpine
WORKDIR /data
COPY --from=source /app /usr/bin
ENTRYPOINT ["moreutils"]
