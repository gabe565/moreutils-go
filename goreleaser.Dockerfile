FROM alpine AS source
WORKDIR /app
COPY moreutils .
RUN ./moreutils install .

FROM scratch
LABEL org.opencontainers.image.source="https://github.com/gabe565/moreutils"
COPY --from=source /app /
ENTRYPOINT ["/moreutils"]
