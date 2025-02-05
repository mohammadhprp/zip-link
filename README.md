# Zip Link 🥃

> Shorten and share URLs effortlessly with Zip Link. Fast, simple, and reliable—make long links a thing of the past!

## Prerequisites

Before getting started, ensure you have the following installed:

- **Docker**: [Install Docker](https://www.docker.com/get-started)

## Technologies Used

- **Go**: Backend service for high performance and scalability.
- **Fiber**: High-performance web framework for Go.
- **MongoDB**: NoSQL database for storing URL information and analytics.
- **Redis**: In-memory data store used for caching and session management.

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/mohammadhprp/zip-link.git
    cd zip-link
    ```

2. Build and run the Docker container:

    ```bash
    docker compose up -d --build
    ```

## Usage

Once the application is running, access the following endpoints for service usage:

- Health Check: `GET http://localhost:3000/api/up`
- Short URL Creation: `POST http://localhost:3000/api/urls`
- Access Shortened URL: `GET http://localhost:3000/:code`
