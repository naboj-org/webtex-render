FROM golang:1.20-alpine as build

WORKDIR /build/
COPY ./webtex_api ./webtex_api
COPY ./webtex_render ./webtex_render
COPY ./webtex_web ./webtex_web
COPY go.mod go.sum ./

RUN go build -o wr ./webtex_render
RUN go build -o wr_web ./webtex_web

FROM texlive/texlive:latest
RUN adduser texuser
USER texuser

WORKDIR /app
COPY --from=build /build/wr /app/wr
COPY --from=build /build/wr_web /app/wr_web
CMD ["/app/wr_web"]
