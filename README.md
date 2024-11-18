# Clicker Application

## Introduction
The Clicker application is designed to track and manage click statistics for various banners. It utilizes a robust Go backend with PostgreSQL to store and retrieve data efficiently.

## Prerequisites
Before you begin, ensure you have the following installed:
- Docker
- Docker Compose
- Go (version 1.22.3 as specified in `go.mod`)
- Protobuf compiler (for generating Go code from `.proto` files)

## Installation

### Cloning the Repository
To get started, clone the repository to your local machine:

bash
git clone https://your-repository-url.com
cd clicker

### Setting Up the Environment
Copy the `.env.example` file to create your own environment variables file:

bash
cp .env.example .env

Make sure to adjust the `.env` file with your specific configurations.

### Building the Application
You can build the application using Docker Compose:

bash
docker-compose up --build


## Usage

### Starting the Application
To start all services, including the PostgreSQL database and the Clicker application:

bash
make start


### Accessing the Application
The application can be accessed at:
- REST API: `http://localhost:8080`
- gRPC services: Running on port `50051`

### Stopping the Application
To stop the application and remove containers:

bash
make down


## Development

### Migrations
To apply database migrations:

bash
make migrate


### Seeding the Database
To seed the database with initial data:

bash
make seed


### Generating Protobuf Files
To generate Go code from `.proto` files:

bash
make proto
