#  SpotSync API

> Smart Parking & EV Charging Reservation System

A centralized backend API for busy airports and malls to manage parking zones, specifically handling the high-demand reservation of limited EV charging spots.

---

## Live URL

```
https://spotsync-api.onrender.com
```

---

## Features

- JWT-based Authentication (Register & Login)
- Role-based Access Control (driver & admin)
- Parking Zone Management (CRUD)
- EV Spot Reservation with Concurrency Protection
- Race Condition Prevention using DB Transaction + FOR UPDATE row lock
- Dynamic available_spots calculation

---

## Tech Stack

| Technology | Purpose |
|---|---|
| Go 1.22 | Backend language |
| Echo v4 | HTTP web framework |
| GORM | ORM for database operations |
| PostgreSQL (NeonDB) | Relational database |
| JWT (golang-jwt/v5) | Authentication tokens |
| bcrypt | Password hashing |
| validator/v10 | Request validation |

---

## 🏛️ Architecture

Feature-Based Domain-Driven Design (DDD)

```
backend/
├── cmd/
│   └── main.go                    # Application entry point
├── internal/
│   ├── auth/                      # JWT token generation & validation
│   ├── config/                    # Database connection & config
│   ├── domain/
│   │   ├── user/                  # User auth feature
│   │   │   ├── dto/               # Request/Response DTOs
│   │   │   ├── model.go           # GORM model
│   │   │   ├── repository.go      # Database operations
│   │   │   ├── service.go         # Business logic
│   │   │   └── handler.go         # HTTP handlers
│   │   ├── zone/                  # Parking zone feature
│   │   │   ├── dto/
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── handler.go
│   │   └── reservation/           # Reservation feature
│   │       ├── dto/
│   │       ├── model.go
│   │       ├── repository.go
│   │       ├── service.go
│   │       └── handler.go
│   ├── httpresponse/              # Standard response helpers
│   ├── middlewares/               # JWT middleware
│   └── server/                    # Echo server setup & routes
├── .env.example
├── go.mod
└── go.sum
```

### Layer Interaction

```
Request → Handler → Service → Repository → Database
                ↑
           JWT Middleware
```

- **Handler** — Bind & validate DTOs, extract JWT claims, return JSON
- **Service** — Business logic, password hashing, JWT generation
- **Repository** — All GORM operations, transactions, row locks
- **DTO** — Request/Response structs (never expose GORM models directly)

---

## Local Setup

### Prerequisites
- Go 1.22+
- PostgreSQL (NeonDB / Supabase)

### Steps

**1. Clone the repo**
```bash
git clone https://github.com/HST159075/spotSync.git
cd spotSync/backend
```

**2. Install dependencies**
```bash
go mod tidy
```

**3. Setup environment variables**
```bash
cp .env.example .env
```

Edit `.env`:
```env
DB_HOST=your-neondb-host
DB_USER=your-username
DB_PASSWORD=your-password
DB_NAME=neondb
DB_PORT=5432
JWT_SECRET=your-secret-key
PORT=8080
```

**4. Run the server**
```bash
go run cmd/main.go
```

Server starts at `http://localhost:8080`
live url: https://spotsync-1q5j.onrender.com

---

## API Endpoints

### Auth
| Method | Endpoint | Access | Description |
|---|---|---|---|
| POST | `/api/v1/auth/register` | Public | Register new user |
| POST | `/api/v1/auth/login` | Public | Login & get JWT token |

### Parking Zones
| Method | Endpoint | Access | Description |
|---|---|---|---|
| GET | `/api/v1/zones` | Public | Get all zones with available spots |
| GET | `/api/v1/zones/:id` | Public | Get single zone |
| POST | `/api/v1/zones` | Admin | Create new zone |
| PUT | `/api/v1/zones/:id` | Admin | Update zone |
| DELETE | `/api/v1/zones/:id` | Admin | Delete zone |

### Reservations
| Method | Endpoint | Access | Description |
|---|---|---|---|
| POST | `/api/v1/reservations` | Auth | Reserve a spot |
| GET | `/api/v1/reservations/my-reservations` | Auth | Get my reservations |
| DELETE | `/api/v1/reservations/:id` | Auth | Cancel reservation |
| GET | `/api/v1/reservations` | Admin | Get all reservations |

---

## Authentication

All protected endpoints require:
```
Authorization: Bearer <token>
```

---

## Concurrency Solution

The reservation system uses **DB Transaction + FOR UPDATE row lock** to prevent race conditions:

```go
tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID)
```

This ensures that when two drivers try to book the last EV spot simultaneously, only one succeeds.

---

##  Environment Variables

| Variable | Description |
|---|---|
| `DB_HOST` | PostgreSQL host |
| `DB_USER` | Database username |
| `DB_PASSWORD` | Database password |
| `DB_NAME` | Database name |
| `DB_PORT` | Database port (5432) |
| `JWT_SECRET` | Secret key for JWT signing |
| `PORT` | Server port (default: 8080) |
