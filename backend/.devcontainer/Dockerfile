# Use the official Golang image as a base
FROM golang:1.22.4

# Install necessary tools
RUN apt-get update && apt-get install -y \
    git \
    curl \
    vim \
    apt-transport-https \
    ca-certificates \
    gnupg

# Add the Cloud SDK distribution URI as a package source
RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] http://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list

# Import the Google Cloud Platform public key
RUN curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -

# Update and install the Cloud SDK
RUN apt-get update && apt-get install -y google-cloud-sdk

# Set environment variables
ENV PORT=8080

# Set the current working directory inside the container
WORKDIR /workspace

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if go.mod and go.sum are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Set the environment variable for the GCP project
ENV GCP_PROJECT=read-robin-2e150
ENV GCP_PROJECT_NUMBER=31059740280
ENV ENV=development
ENV GOLANG_SAMPLES_PROJECT_ID=read-robin-2e150
ENV GOPROXY=https://proxy.golang.org,direct


RUN gcloud config set project $GCP_PROJECT

ENV GOOGLE_APPLICATION_CREDENTIALS=/workspace/secrets/service_account_credentials.json

# Expose port 8080 for the app
EXPOSE 8080

# Set the default command for the container
CMD ["go", "run", "main.go"]