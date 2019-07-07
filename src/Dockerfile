FROM golang:1.12.5-stretch as base
RUN mkdir /src
WORKDIR /src

COPY go.mod /src/go.mod
COPY go.sum /src/go.sum

RUN go mod download
COPY . /src/

FROM base as build
RUN go build -o airbag .

FROM build as tests
RUN go test ./...

FROM gcr.io/distroless/base as release
COPY --from=build /src/airbag /app/airbag

ENTRYPOINT [ "/app/airbag" ]