FROM google/cloud-sdk:latest

# Upgrade pip to avoid connection issues
RUN pip install --upgrade pip

# Install any needed packages specified in requirements.txt
COPY requirements.txt ./requirements.txt
RUN pip install --no-cache-dir --trusted-host pypi.org --trusted-host files.pythonhosted.org -r ./requirements.txt

# Set the working directory in the container
WORKDIR /workspace

# Copy the current directory contents into the container at /workspace
COPY . /workspace

# Set the environment variable for the GCP project
ENV GCP_PROJECT=dvan-genai-poc
ENV GCP_PROJECT_NUMBER=651397568677
RUN gcloud config set project $GCP_PROJECT

ENV GOOGLE_APPLICATION_CREDENTIALS=/workspace/secrets/service_account_credentials.json

ENV PYTHONPATH "${PYTHONPATH}:/workspace"

# # Install dependencies for Node.js and npm
# RUN apt-get update && apt-get install -y \
#     curl \
#     gnupg \
#     && rm -rf /var/lib/apt/lists/*

# # Install Node.js and npm
# RUN curl -fsSL https://deb.nodesource.com/setup_14.x | bash - \
#     && apt-get install -y nodejs

# # Verify Node.js and npm installation
# RUN node -v && npm -v

# # Install global npm packages for Chrome extension development
# RUN npm install -g http-server

# Expose the necessary ports
EXPOSE 5000
EXPOSE 8080