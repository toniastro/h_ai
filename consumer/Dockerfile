FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /consumer

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o consumer .

# Command to run the executable
CMD ./consumer