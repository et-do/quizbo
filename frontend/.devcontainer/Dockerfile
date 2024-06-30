# Use the official Node.js image as a base
FROM node:18

# Set environment variables
ENV PORT=3000

# Set the current working directory inside the container
WORKDIR /workspace/frontend

# Install necessary global packages (e.g., create-react-app)
RUN npm install -g create-react-app
RUN npm install -g firebase-tools

# Expose port 3000 for the app
EXPOSE 3000

USER root
# Set the default command for the container
CMD ["bash"]
