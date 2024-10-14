FROM alpine AS source
WORKDIR /app
COPY moreutils .
RUN ./moreutils install .

FROM alpine
WORKDIR /data
LABEL org.opencontainers.image.source="https://github.com/gabe565/moreutils-go"
COPY --from=source /app /usr/bin
ENTRYPOINT ["moreutils"]
