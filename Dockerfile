FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bantercore .

FROM alpine:3.20

RUN adduser -D -H -u 10001 appuser
COPY --from=build /bantercore /bantercore
USER appuser
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/bantercore"]