# E-Commerce Microservices

This is a microservices-based e-commerce application consisting of the following services:
- **Inventory Service**: Manages products and categories
- **Order Service**: Handles order processing and management
- **API Gateway**: Routes requests to appropriate services

## Prerequisites

- Docker
- Docker Compose

## Running with Docker

1. Clone the repository:
```bash
git clone https://github.com/your-username/e-commerce.git
cd e-commerce
```

2. Start the services:
```bash
docker-compose up -d
```

3. Check the status of the services:
```bash
docker-compose ps
```

4. Access the API Gateway:
```
http://localhost:8080
```

## Service Endpoints

### API Gateway
- Base URL: `http://localhost:8080`

### Inventory Service (direct access)
- Base URL: `http://localhost:8081`
- Health Check: `GET /health`
- API Routes:
  - Products: `GET/POST /api/v1/products`
  - Categories: `GET/POST /api/v1/categories`

### Order Service (direct access)
- Base URL: `http://localhost:8082`
- Health Check: `GET /health`
- API Routes:
  - Orders: `GET/POST /api/v1/orders`

## Database Setup

This project is configured to use PostgreSQL. When running with Docker, the database service will start automatically, but **you will need to create the databases manually** as follows:

1. Connect to the PostgreSQL container:
```bash
docker exec -it e-commerce-postgres bash
```

2. Connect to PostgreSQL:
```bash
psql -U postgres
```

3. Create the inventory database:
```sql
CREATE DATABASE inventory_db;
```

4. Create the order database:
```sql
CREATE DATABASE order_db;
```

5. Exit PostgreSQL and the container:
```
\q
exit
```

## Stopping the Application

To stop all services:
```bash
docker-compose down
```

To stop all services and remove volumes (including database data):
```bash
docker-compose down -v
``` 