# task-manager-in-docker
Discover Task Manager in Docker! 🐳 Embrace the power of multistage builds to shrink image sizes and optimize resource utilization. Join us on a journey of streamlined deployment, enhanced efficiency, and seamless scalability. Let's revolutionize DevOps together! 💼✨ #MultistageBuilds #Containerization #Efficiency

## CI/CD: GitHub Actions -> Docker Hub -> EC2

This repository now includes a workflow at `.github/workflows/deploy-ec2.yml`.

It does the following on every push to `main`:

1. Builds the Docker image.
2. Pushes the image to Docker Hub with the `latest` tag.
3. Connects to your EC2 instance over SSH.
4. Pulls the latest image, starts a candidate container on a temporary port, checks health, then replaces the old container.

### Required GitHub Secrets

Add these in your GitHub repository under `Settings -> Secrets and variables -> Actions`:

- `DOCKERHUB_USERNAME`: your Docker Hub username.
- `DOCKERHUB_TOKEN`: Docker Hub access token (not password).
- `DOCKER_IMAGE_NAME`: full Docker Hub image name, for example `yourname/task-manager`.
- `EC2_HOST`: public IP or DNS of the EC2 instance.
- `EC2_USER`: SSH user (for Ubuntu AMI usually `ubuntu`).
- `EC2_SSH_KEY`: private key content (`.pem`) used by GitHub Actions to SSH.
- `APP_NAME`: deployment name, for example `task-manager`.
- `HOST_PORT`: host port to publish (optional for this CLI app).
- `APP_PORT`: container port to publish (optional for this CLI app).
- `TEMP_PORT`: temporary host port used while the new container is checked.
- `HEALTHCHECK_PATH`: HTTP path used to verify the new container, for example `/health`.

### Notes About This Project

This app is currently a CLI program (interactive stdin) and not an HTTP web server.

- If you are only running it as a background container, keep `HOST_PORT`, `APP_PORT`, `TEMP_PORT`, and `HEALTHCHECK_PATH` empty.
- If you later convert it into an API/web service, set `APP_PORT`, `TEMP_PORT`, and `HEALTHCHECK_PATH` (for example `8080`, `8081`, and `/health`).

### EC2 One-Time Checks

Run these once on EC2:

```bash
sudo systemctl enable docker
sudo systemctl start docker
sudo usermod -aG docker "$USER"
```

Then reconnect to the server so group changes apply.
