# Backend Developer Project Documentation

## Overview
I was tasked with building REST API in Golang using Gin library such that CRUD actions support dual database (PostgreSQL and MongoDB) with JWT Authentication.

## Tech Stack
- Language: Golang
- Library: Gin
- Databases: PostgreSQL (GORM), MongoDB
- Auth: JWT, Bcrypt
- Environment: .env support using godotenv
- UI for testing endpoints: HTML, CSS, Javascript

## Key Design Decisions

### 1. Why Repo Factory Pattern?
Following the project specs strictly would have meant choosing between databases in each route handler, leading to:
- Data inconsistencies (user creates in Postgres, searches in Mongo, gets null)
- Database logic scattered across handlers
- Maintenance nightmare

**Solution**: Built a RepoFactory that abstracts database selection away and implements dual-write with logical atomicity. Since actual atomicity is impossible across different database systems, I implemented rollback logic to maintain consistency.

### 2. Authentication: Bearer + HTTP-only Cookies
The spec required Bearer tokens, but I have security concerns with this approach - JavaScript on the frontend should never be able to touch auth tokens (XSS vulnerability).

**Solution**: Middleware supports both:
- HTTP-only cookies (my preference - more secure)
- Bearer tokens in headers (spec compliance)

### 3. MongoDB Setup
Had issues starting a MongoDB Atlas instance, so I used Docker instead. Pulled MongoDB and PostgreSQL images for local development.

## Features

### Dual Database Architecture: Implemented dual write with logical atomicity for data consistency.
- Writes: Data is written to PostgreSQL and MongoDB simultaneously such that if one fails, both fails
- Reads: Data can be read from either database by passing a "storage" key in the request body which can be either "sql" for PostgreSQL or "doc" for MongoDB
- Atomicity: If either database write fails, the operation rolls back on both databases

**Repository Pattern -> Factory Pattern -> Generic Interfaces**

### Clean Architecture
- Generic Repository[T] interface so other modules can stay dumb about implementation
- RepoFactory[T] manages writes with rollback logic across the two databases
- RepoFactory[T] implements PostgresRepo[T] and MongoRepo[T]

### Auth System
- JWT token generation and validation
- Tokens via HTTP-only cookies and Bearer headers
- Middleware protection on all api routes except pages routes and auth routes

### API Endpoints
- Public:
    POST /api/register - Create account with password validation
    POST /api/login - Authenticate and receive JWT
- Protected (JWT required):
    GET /api/customers?storage=sql|doc - List customers
    GET /api/customers/:id?storage=sql|doc - Get single customer
    PUT /api/customers/:id - Update customer
    DELETE /api/customers/:id - Soft delete
    PATCH /api/customers/:id - Hard delete
    GET /api/profile - Get current user profile
    POST /api/logout - Invalidate session

### Data Models
- Customer: ID (UUID), Name, Email (unique), Phone, Password (hashed), Timestamps, Soft delete
- Session: ID, CustomerID, Email, Token, TTL, Soft delete

## Getting Started

1. Clone the repository from GitHub
2. Copy `.env.example` to `.env` and populate:
   - Only `POSTGRES_URL` needed if not using Docker
   - Skip `MONGOEXPRESS_USER` and `MONGOEXPRESS_PASS` if not using Docker
3. Run `make start` (requires Docker Engine)
   - Or `make run` if not using Docker
   
The make command installs dependencies and starts the server.

## Environment Variables
- POSTGRES_USER
- POSTGRES_PASSWORD
- POSTGRES_DB
- POSTGRES_URL
- MONGODB_URI
- MONGO_DB_NAME
- MONGOEXPRESS_USER
- MONGOEXPRESS_PASS
- PORT
- JWT_SECRET
- GIN_MODE (production/development)
- DOMAIN_URL

## Project Structure
```
main.go              # Server initialization & routing
/db                  # Database connections & migrations
/handlers            # HTTP request handlers
/middleware          # JWT authentication
/models              # Data structures
/repo                # Repository pattern implementation
/utils               # JWT, password hashing, helpers
/static              # UI For testing routes
docker-compose.yaml  # Docker setup for PostgreSQL and MongoDB
```

## Result
Fully functional API meeting all requirements + bonus features (repository pattern, clean architecture, dual write logical atomicity)