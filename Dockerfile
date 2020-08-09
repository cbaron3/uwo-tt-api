FROM golang:1.14 as builder

# Set the Current Working Directory inside the container
WORKDIR /app


# # We want to populate the module cache based on the go.{mod,sum} files.
# COPY go.mod .
# COPY go.sum .

# RUN go mod download -x

# COPY --from=itinance/swag /root/swag /usr/local/bin
# COPY . .

# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
COPY go.* ./
RUN go mod download -x

# Copy local code to the container image.
COPY --from=itinance/swag /root/swag /usr/local/bin
COPY . ./
# COPY ./.env /

RUN swag init

# Build the binary.
# -mod=readonly ensures immutable go.mod and go.sum in container builds.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server
# COPY --from=builder /app/.env .env

# # Run the binary program produced by `go install`
# COPY ./.env /
RUN touch .env

# Run the web service on container startup.
CMD ["/server"]