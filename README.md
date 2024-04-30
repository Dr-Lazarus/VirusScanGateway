# VirusScanGateway
- ![example workflow](https://github.com/Dr-Lazarus/VirusScanGateway/actions/workflows/EC2-container-webapp.yml/badge.svg)
- ![example workflow](https://github.com/Dr-Lazarus/VirusScanGateway/actions/workflows/Testing-environment.yml/badge.svg)
  
VirusScanGateway is a Go-based web server designed to facilitate the secure uploading and scanning of files for viruses and malware. It integrates with the VirusTotal API to provide real-time scanning results, which are then stored and displayed from a PostgreSQL database.

## Features

- File upload and secure storage of results in PostGre SQL DB.
- Real-time virus and malware scanning using VirusTotal API.
- Hosted on EC2 with custom domain: www.virusscanapi.lat
- Domain secured with HTTPs.
- Connected to AWS RDS.
- CI/CD Pipeline Production: Commits to Main will be automatically deployed to AWS EC2 Container.
- CI/CD Pipeline Development: Commits (Non-Main) and PRs will be automatically tested in a custom testing environment.
- Comprehensive Integration testing suite to test endpoints.
- AWS Secrets Manager to manage production environment key-value pairs.
- AWS Route 53 to provide authoritative NS and a hosted DNS Zone for www.virusscanapi.lat

## AWS Cloud Architecture
![image](https://github.com/Dr-Lazarus/VirusScanGateway/assets/99006087/d7b5929b-671c-4b49-b208-c00988ec255d)


## CI - Dev
- Any commits to non-main branches will trigger a workflow which runs 3 integration tests.


## How to develop
To install VirusScanGateway, follow these steps:

1. Clone the repository:

```bash
git clone https://github.com/yourusername/VirusScanGateway.git
cd VirusScanGateway
```

2. Start your PostGreSQL Instance on Local
3. Modify .env.dev file to comment CI/CD and uncomment local
```plaintext
# DEV ENVIRONMENT (LOCAL) 
 DATABASE_URL=postgres://cloudsine:password@host.docker.internal:5432/VirusScanGatewayDB?sslmode=disable
 VIRUSTOTAL_API_KEY=Your-Test-API-Key
 ENVIRONMENT=DEV

# DEV ENVIRONMENT (CI/CD)
DATABASE_URL=postgresql://postgres:password@db:5432/postgres?sslmode=disable
VIRUSTOTAL_API_KEY=Your-Test-API-Key
ENVIRONMENT=DEV
```
5. Run the container:
```bash
docker-compose -f docker-compose-dev-local.yml build
docker-compose -f docker-compose-dev-local.yml up
```
5. Server started on localhost:8080

## Prerequisites

Before you begin, ensure you have met the following requirements:

- You have installed Docker and Docker Compose.
- You have a VirusTotal API key. If you don't have one, you can register for an API key at [VirusTotal](https://www.virustotal.com/gui/join-us).


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
