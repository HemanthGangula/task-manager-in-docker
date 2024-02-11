###########################################
# BASE IMAGE FOR BUILD STAGE
###########################################

FROM ubuntu:latest AS build

# Install Go
RUN apt-get update && apt-get install -y golang-go

# Set the environment variable to disable Go modules
ENV GO111MODULE=off

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 go build -o /app/task_manager .

###########################################
# FINAL IMAGE
###########################################

FROM scratch

# Copy the executable from the build stage to the final image
COPY --from=build /app/task_manager /app/task_manager

# Set the entrypoint for the container
ENTRYPOINT ["/app/task_manager"]

