# Step 1: Use an official Go image
FROM golang:1.22.3

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy your Go code into the container
COPY . .

# Step 4: Download dependencies
RUN go mod download

# Step 5: Build the Go binary
RUN go build -o myapp

# Step 6: Command to run your app
CMD ["./myapp"]
