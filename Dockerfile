FROM golang AS builder
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -o /bin/wikibot ./cmd/wikibot

FROM scratch
WORKDIR /
COPY --from=builder /bin/wikibot /bin/wikibot
USER nobody
ENTRYPOINT ["/bin/wikibot"]
