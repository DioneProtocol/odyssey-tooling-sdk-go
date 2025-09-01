#!/usr/bin/env bash

# Check if Docker is installed
if ! command -v docker &>/dev/null; then
	echo "Docker is not installed. Please install Docker and try again."
	exit 1
fi

# Check if Docker Compose is installed
if ! docker compose version &>/dev/null; then
	echo "Docker Compose is not installed or not configured correctly. Please install Docker Compose and try again."
	exit 1
fi

# Provide docker-compose systemctl unit file
cat <<EOF | sudo tee /etc/systemd/system/odyssey-cli-docker.service
[Unit]
Description=Odyssey CLI Docker Compose Service
Requires=docker.service
After=docker.service

[Service]
User=ubuntu
Group=ubuntu
Restart=on-failure
ExecStart=/usr/bin/docker compose -f /home/ubuntu/.odyssey-cli/services/docker-compose.yml up 
ExecStop=/usr/bin/docker compose -f /home/ubuntu/.odyssey-cli/services/docker-compose.yml down

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd manager configuration
sudo systemctl daemon-reload

# Enable the new service
sudo systemctl enable odyssey-cli-docker.service

echo "Service created and enabled successfully."
