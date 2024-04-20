# VirusScanGateway

VirusScanGateway is a Go-based web server designed to facilitate the secure uploading and scanning of files for viruses and malware. It integrates with the VirusTotal API to provide real-time scanning results, which are then stored and displayed from a PostgreSQL database.

## Features

- File upload and secure storage.
- Real-time virus and malware scanning using VirusTotal API.
- Persistent results storage with PostgreSQL.
- Dockerized environment for easy setup and deployment.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- You have installed Docker and Docker Compose.
- You have a VirusTotal API key. If you don't have one, you can register for an API key at [VirusTotal](https://www.virustotal.com/gui/join-us).

## Installation

To install VirusScanGateway, follow these steps:

1. Clone the repository:

```bash
git clone https://github.com/yourusername/VirusScanGateway.git
cd VirusScanGateway
```

2. Create a `.env.dev` file in the root directory with the following content, replacing `your_virustotal_api_key` with your actual API key:

```env
DATABASE_URL=postgres://cloudsine:password@localhost:5432/VirusScanGatewayDB?sslmode=disable
APP_ENV=DEV
VIRUSTOTAL_API_KEY=your_virustotal_api_key
```

3. Build and run the application using Docker Compose:

```bash
docker-compose up --build
```

## Usage

Once the application is running, you can upload files to be scanned via the provided web interface at `http://localhost:8080`.

Results from the scans will be available in the web interface, and stored in the PostgreSQL database for future reference.

## Development

To set up a development environment for contributing to VirusScanGateway, you can use the following additional steps:

1. Install Go locally on your machine.
2. Install a PostgreSQL database and configure it according to the application's requirements.

For local development, you can run the server directly using:

```bash
make start
```

Make sure your local environment variables are set according to the `.env.dev` file.

## License

Distributed under the MIT License. See `LICENSE` for more information.
