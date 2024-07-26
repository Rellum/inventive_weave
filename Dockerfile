# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS base
WORKDIR /src
COPY go.mod go.sum .
RUN go mod download
COPY . .

FROM base AS build-creators
RUN go build -o /bin/creators ./svc/creators

FROM base AS build-fe
RUN go build -o /bin/fe ./svc/fe

FROM scratch AS creators
COPY --from=build-creators /bin/creators /bin/
ENTRYPOINT [ "/bin/creators" ]

FROM scratch AS fe
COPY --from=build-fe /bin/fe /bin/
ENTRYPOINT [ "/bin/fe" ]
