# Task Management Microservice

This repository contains a simple Task Management System implemented as a microservice in Go, using Chi router and PostgreSQL database. 
The system allows users to create, read, update, and delete tasks with features like pagination and filtering by status.

## Problem Breakdown and Design Decisions

### Architecture Overview

The service follows a clean architecture approach with clear separation of concerns:

1. **API Layer**: Handles HTTP requests and responses using Chi router
2. **Service Layer**: Contains business logic and validation
3. **Repository Layer**: Manages data access and persistence with PostgreSQL
4. **Domain Layer**: Defines core entities and interfaces

### Design Decisions

- **PostgreSQL Storage**: Persistent data storage using PostgreSQL database
- **Chi Router**: Lightweight, idiomatic router for Go HTTP services
- **Dockerized**: Complete Docker and Docker Compose setup for easy deployment
- **REST API**: The service exposes a RESTful API for CRUD operations on tasks
- **Pagination**: Implemented with `limit` and `offset` parameters
- **Filtering**: Allows filtering tasks by status (Pending, InProgress, Completed)
- **Graceful Shutdown**: Properly handles shutdown signals for clean termination

## Microservices Concepts Demonstrated

### Single Responsibility Principle
Each component has a clear, single responsibility:
- **Repository**: Data access and persistence
- **Service**: Business logic and validation
- **API Handlers**: Request handling and response formatting
- **Config**: Application configuration
- **Main**: Application bootstrapping and orchestration

### API Design
- Clear, consistent RESTful endpoints following standard conventions:
    - `GET /api/v1/task/`: List tasks with pagination and filtering
    - `POST /api/v1/task/`: Create a new task
    - `GET /api/v1/task/{id}`: Get a specific task
    - `PUT /api/v1/task/{id}`: Update a task
    - `DELETE /api/v1/task/{id}`: Delete a task

### Scalability
The service can be scaled horizontally in several ways:

1. **Stateless Design**: The service is stateless, making it easy to run multiple instances behind a load balancer.
2. **Container Orchestration**: The service is containerized and can be deployed in a Kubernetes cluster for automated scaling.
3. **Database Scaling**: PostgreSQL can be scaled using read replicas or sharding for higher throughput.

## Installation and Setup

### Prerequisites
- Docker and Docker Compose

### Getting Started

1. Clone the repository:
```bash
git clone https://github.com/mrpawan-gupta/TaskEase
cd TaskEase
```

2. Start the services using Docker Compose:
```bash
docker-compose up -d
```

The service will be available at [http://localhost:8080](http://localhost:8080)

3. To stop the services:
```bash
docker-compose down
```

4. To view logs:
```bash
docker-compose logs -f
```

## API Documentation

### Create a Task
```
POST /api/v1/task/
Request:
{
  "title": "Complete microservice assignment",
  "description": "Implement a task management system using Go",
}
```

### Get a Task
```
GET /api/v1/task/{id}
```

### Update a Task
```
PUT /api/v1/task/{id}
```

### Delete a Task
```
DELETE /api/v1/task/{id}
```

### List Tasks (with pagination and filtering)
```
GET /api/v1/task/?limit=10&offset=0
```

### Health Check
```
GET /health
```