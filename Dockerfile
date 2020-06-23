FROM golang

# Set the Current Working Directory inside the container
WORKDIR /tmp/uwo-tt-api

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download -x

COPY --from=itinance/swag /root/swag /usr/local/bin
COPY . .

RUN swag init
RUN go build -o ./out/uwo-tt-api .

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
COPY ./.env /

#ENTRYPOINT CompileDaemon -exclude-dir=.git -exclude-dir=docs --build="make build-dev" --command=./main
CMD ["/tmp/uwo-tt-api/out/uwo-tt-api"]