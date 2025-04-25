# Arago

**Arago** is a microservice-based application designed to manage and track ad clicks. It includes two main services:

- `adserver`: Handles ad creation and click tracking.
- `tracker`: Receives and stores click events in MongoDB.

## Features

- gRPC communication between services
- MongoDB for data storage
- Redis-compatible cache with DragonflyDB
- Docker-based development environment

## Getting Started

### Prerequisites

- Docker & Docker Compose
- `make` (optional, for convenience)

### Setup

```bash
# Clone the repo
git clone https://github.com/Bellzebuth/arago.git
cd arago

# Build and start all services
docker-compose up --build
```

Or using Make:

```bash
make proto   # Generate protobuf files
make build   # Build Docker images
make up      # Start containers
```

### gRPC Testing

You can use `grpcurl` to test endpoints.

### Create an Ad

```bash
grpcurl -plaintext -d '{
  "title": "Nike Air Max",
  "description": "You will love it",
  "url": "https://nike.air.max.com"
}' localhost:50051 ad.AdService/CreateAd
```

### Get an Ad

```bash
grpcurl -plaintext -d '{
  "id": "REPLACE_WITH_AD_ID"
}' localhost:50051 ad.AdService/GetAd
```

### Track a Click

```bash
grpcurl -plaintext -d '{
  "adId": "REPLACE_WITH_AD_ID"
}' localhost:50051 ad.AdService/ClickAd
```

## Environment

All environment variables are set in `.env`. Example:

```bash
MONGO_URI=mongodb://mongo:27017
DB_NAME=arago
AD_COLLECTION=ads
TRACKER_COLLECTION=tracker
ADSERVER_PORT=50051
TRACKER_PORT=50052
DRAGONFLY_PORT=6379
```

Update ports, Mongo URI, and collection name as needed.

## License

MIT
