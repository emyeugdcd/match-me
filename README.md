<div align="center">
  <h1> Match-Me: The Developer Connection App</h1>
  <p>
    <strong>A full-stack social networking platform engineered to connect coders based on custom matching algorithms.</strong>
  </p>
  <p>
    <img alt="Go" src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white" />
    <img alt="React" src="https://img.shields.io/badge/react-%2320232a.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB" />
    <img alt="Docker" src="https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white" />
    <img alt="SQLite" src="https://img.shields.io/badge/sqlite-%2307405e.svg?style=for-the-badge&logo=sqlite&logoColor=white" />
  </p>
</div>

<br />

## Overview

Dubbed "Tinder for Coders", **Match-Me** is a robust, full-stack application designed to connect developers. It leverages a custom-built recommendation engine to match users based on their technical skills, interests, and geographical proximity. The application features real-time WebSocket communication, secure authentication, and a containerized deployment setup.

### Core Features

- **Intelligent Matchmaking:** An advanced, weighted recommendation algorithm that queries user profiles against location (Haversine formula), interests, stack, and hobbies, generating dynamic match scores.
- **Real-time Messaging:** WebSockets power instant messaging and live online presence tracking between connected users.
- **Secure Authentication:** JWT-based session management seamlessly tracks authenticated sessions, combined with robust password hashing.
- **Dynamic Profiles:** Complete user lifecycle management including registration, profile customisation, image uploads, and customizable preferences.
- **Streamlined Workflow:** Manage incoming, outgoing, and mutual connections efficiently.

## Tech Stack & Architecture

- **Backend:** Go (Golang), WebSockets
- **Frontend:** React.js
- **Database:** Relational persistence (SQLite/PostgreSQL compatible)
- **Infrastructure:** Docker & Docker Compose for isolated, consistent environments

## The Recommendation Engine

At the core of and separating Match-Me from standard CRUD apps is its `recommendationService.go`. The algorithm evaluates compatibility via:
1. **Exclusion Filtering:** Ensures matched or dismissed profiles are excluded.
2. **Proximity Fencing:** Computes geographic distance natively via the Haversine formula against a defined radius threshold.
3. **Scoring Logic:** Weights shared attributes dynamically (e.g., Interests natively score 3x more than Food Preferences). 
4. **Data Normalization:** Maps the raw score calculations perfectly onto a localized 0-100% threshold metric.

## Getting Started

The project is heavily encapsulated using Docker to ensure dependency-free installations.

### Prerequisites
- Docker & Docker Compose

### Running the Environment

**Run Everything (Docker + Local):**
Spins up required databases and triggers both client and server gracefully.
```bash
make run
```

**Run Components Individually:**
```bash
make run-server # Starts the Go backend
make run-client # Starts the React frontend
```

### Seeding Test Data
Quickly populate your instance with fake users for testing:
```bash
curl -X POST http://localhost:8080/dev/seed-users -H "Content-Type: application/json" -d '{"count": 50}'
```
