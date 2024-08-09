# Use the official Ubuntu image from the Docker Hub
FROM ubuntu:latest

# Update the package list and install curl
RUN apt-get update && apt-get install -y curl

# Set the entrypoint to keep the container running
ENTRYPOINT ["tail", "-f", "/dev/null"]
