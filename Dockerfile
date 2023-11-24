FROM golang:alpine AS build
WORKDIR /build
COPY . .
RUN go build -ldflags "-s -w" main.go

FROM jrottenberg/ffmpeg:6-alpine 
COPY --from=build /build/main /bin/
EXPOSE 8090
ENTRYPOINT ["/bin/main"]
