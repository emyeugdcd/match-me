<div align="center">
  <h1>Match-Me: The Developer Connection Platform</h1>
  <p>
    <strong>A modern, full-stack developer matchmaking application engineered with Go, React 19, TypeScript, WebSockets, and containerized CI/CD automation.</strong>
  </p>
  <p>
    <img alt="Go" src="https://img.shields.io/badge/Go-1.23-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white" />
    <img alt="TypeScript" src="https://img.shields.io/badge/TypeScript-5.9-%233178C6.svg?style=for-the-badge&logo=typescript&logoColor=white" />
    <img alt="React" src="https://img.shields.io/badge/React-19-%2320232a.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB" />
    <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-18-%234169E1.svg?style=for-the-badge&logo=postgresql&logoColor=white" />
    <img alt="Docker" src="https://img.shields.io/badge/Docker-Compose-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white" />
    <img alt="GitHub Actions" src="https://img.shields.io/badge/GitHub_Actions-CI-%232671E5.svg?style=for-the-badge&logo=githubactions&logoColor=white" />
  </p>
</div>

<br />

## Overview

**Match-Me** is a full-stack social networking platform designed to connect software engineers based on technical compatibility, shared interests, and geographic proximity. 

The application combines a high-performance **Go** backend, a strictly-typed **React 19 + TypeScript** frontend, persistent **PostgreSQL** storage, and **real-time bi-directional WebSockets** for live presence and chat.

---

## Key Features

- **Intelligent Recommendation Engine**: Custom multi-attribute matchmaking evaluating technical interests, hobbies, music, and food preferences with distance proximity filtering using the Haversine formula.
- **Real-Time Communication**: Bi-directional WebSockets powering instant messaging, typing indicators, and dynamic user presence tracking.
- **Secure Authentication**: JWT-based session authorization, bcrypt password hashing, and protected route guards.
- **Type-Safe Full-Stack Architecture**: Clean separation of concerns with Go services on the backend and comprehensive TypeScript interfaces across the frontend.
- **Automated CI/CD**: End-to-end GitHub Actions pipeline enforcing Go unit tests, TypeScript type checking, ESLint quality gates, and container image builds.

---

## Tech Stack & Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Client Layer                           │
│  React 19 • TypeScript • Vite • Strict Type Safety          │
└──────────────────────────────┬──────────────────────────────┘
                               │ HTTP / REST & WebSockets
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                      Backend Layer                          │
│  Go (Golang 1.23) • Gin Web Framework • Gorilla WebSockets  │
└──────────────────────────────┬──────────────────────────────┘
                               │ SQL Persistence
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                      Database Layer                         │
│  PostgreSQL 18 • Relational Persistence • Connection Pool   │
└─────────────────────────────────────────────────────────────┘
```

* **Backend**: Go 1.23, Gin framework, Gorilla WebSocket, JWT authentication.
* **Frontend**: React 19, TypeScript 5.9, Vite, React Router 7.
* **Database**: PostgreSQL 18 with relational schemas and indexing.
* **DevOps & Infrastructure**: Docker multi-stage builds, Docker Compose orchestration, and GitHub Actions CI.

---

## The Recommendation Engine

At the core of Match-Me is `recommendationService.go`, which implements a 4-phase matchmaking pipeline:

1. **Exclusion Filtering**: Automatically prunes existing connections, pending requests, and dismissed profiles from consideration.
2. **Geographic Proximity**: Calculates real-world distances using the **Haversine formula** against configurable search radii:
   $$\Delta\sigma = 2 \arcsin\left(\sqrt{\sin^2\left(\frac{\Delta\phi}{2}\right) + \cos(\phi_1)\cos(\phi_2)\sin^2\left(\frac{\Delta\lambda}{2}\right)}\right)$$
3. **Multi-Attribute Scoring Logic**: Evaluates compatibility via weighted attribute overlaps:
   * **Interests / Tech Stack**: Weight $3\times$
   * **Hobbies**: Weight $2\times$
   * **Music & Food**: Weight $1\times$
4. **Normalized Match Score**: Maps total compatibility onto a normalized 0–100% threshold metric.

---

## CI/CD Pipeline (GitHub Actions)

The project includes an automated GitHub Actions pipeline (`.github/workflows/ci.yml`) triggering on pushes and pull requests:

| Stage | Tooling | Quality Gates |
|---|---|---|
| **Backend CI** | Go 1.23 | Dependency verification (`go mod verify`), static analysis (`go vet`), unit tests with race detection (`go test -race ./...`), binary compilation. |
| **Frontend CI** | Node.js 20 | Zero-warning ESLint validation (`npm run lint`), TypeScript compile check (`tsc -b`), and production bundle build (`vite build`). |
| **Container Validation** | Docker Buildx | Validates `docker compose config` and executes multi-stage container builds for both client and server images. |

---

## Getting Started

### Prerequisites
- Docker & Docker Compose
- (Optional for local development) Go 1.23+ and Node.js 20+

### Quick Start (Docker Compose)

Spin up the database, backend server, and frontend client with a single command:

```bash
docker compose up --build -d
```

- **Frontend Client**: `http://localhost:5173`
- **Backend API**: `http://localhost:8080`
- **Database**: PostgreSQL on `localhost:5432`

### Running Locally for Development

```bash
# 1. Start the PostgreSQL database
docker compose up -d db

# 2. Run the Go backend
make run-server

# 3. Run the React TypeScript frontend (in a separate terminal)
make run-client
```

### Running Automated Tests

```bash
# Backend unit tests
cd server && go test -v ./...

# Frontend type check & linting
cd client && npm run lint && npm run build
```

### Seeding Test Profiles

Quickly populate your local instance with mock developers for matchmaking:

```bash
curl -X POST http://localhost:8080/dev/seed-users \
  -H "Content-Type: application/json" \
  -d '{"count": 50}'
```

---

## License & Attribution

Developed collaboratively as a full-stack showcase project demonstrating software engineering, real-time networking, and DevOps best practices.
