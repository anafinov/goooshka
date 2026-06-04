# Student Certificate Service

A microservices-based student certificate processing system with asynchronous job handling.

## Technology Stack

- **Backend & Worker:** Go 1.25, Clean Architecture
- **Database:** PostgreSQL 15 (with auto-migrations)
- **Message Broker:** RabbitMQ
- **Frontend:** Vanilla JavaScript, HTML, CSS (Glassmorphism UI)
- **Monitoring:** Prometheus, Grafana
- **Infrastructure:** Docker, Docker Compose, Nginx

## Architecture Overview

The system is composed of several independent components:

1. **Backend API (`/backend`)**: The primary HTTP REST API. Implemented using Clean Architecture principles (Handler -> Service -> Repository). Responsibilities include:
   - User registration and JWT-based authentication.
   - Processing certificate requests and persisting them in PostgreSQL with a `pending` status.
   - Publishing request events to the `cert_requests` RabbitMQ queue.
   - Retrieving user-specific certificate statuses.
2. **Worker (`/worker`)**: An asynchronous background service. It consumes messages from RabbitMQ, simulates a long-running certificate generation process (5 seconds), and updates the request status to `ready` in the database.
3. **Frontend (`/frontend`)**: A lightweight Single Page Application (SPA) that interacts with the REST API. Served via Nginx.
4. **Load Tester (`/loadtester`)**: A Go script that continuously generates API requests (registration/login) to simulate active user load for monitoring purposes.

## Quick Start

**Prerequisites:** Docker and Docker Compose must be installed on your system.

1. Navigate to the project directory:
   ```bash
   cd /Users/anafinov/.gemini/antigravity/scratch/student-cert-service
   ```

2. Build and start all services in detached mode:
   ```bash
   docker-compose up -d --build
   ```
   *Note: Initial startup may take a few minutes as Docker pulls the required base images and builds the Go binaries.*

3. To stop and remove the containers, run:
   ```bash
   docker-compose down
   ```

## Services and Ports

Once the containers are running, the following interfaces will be available:

- **Frontend Application:** http://localhost (Port 80)
- **Backend API:** http://localhost:8080
- **RabbitMQ Management UI:** http://localhost:15672 (Credentials: `guest` / `guest`)
- **Grafana Dashboards:** http://localhost:3000 (Credentials: `admin` / `admin`)
- **Prometheus UI:** http://localhost:9090

## API Reference

| Method | Endpoint | Description | Authentication |
| :--- | :--- | :--- | :--- |
| `POST` | `/register` | Register a new student account | Not Required |
| `POST` | `/login` | Authenticate and retrieve a JWT | Not Required |
| `GET` | `/requests` | Retrieve all certificate requests for the current user | Required |
| `POST` | `/request` | Submit a new certificate request | Required |
| `GET` | `/metrics` | Export Prometheus metrics | Not Required |

## Key Features

- **Graceful Shutdown**: The backend server ensures all active HTTP connections are processed before terminating.
- **Security**: Passwords are securely hashed using `bcrypt`. Cross-Origin Resource Sharing (CORS) is explicitly configured. All protected endpoints require a `Bearer JWT` header.
- **Performance**: Docker images are optimized for size using `alpine` Linux and multi-stage builds. Go binaries are compiled with strict minimal dependencies.
