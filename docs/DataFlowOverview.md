# Data Flow Overview

How will your microservices interact with our tool?

## Step 1️⃣ - SDK/Webhook Collects Data

- Users install an SDK (Golang) or send HTTP requests to our API.
- Each call tracks service interactions (e.g., AuthService calls UserService).
- Example:

```json
{
  "service_1": "AuthService",
  "service_2": "UserService",
  "method": "CALL"
}
```

- Data Sent To → Our Backend API

## Step 2️⃣ - Backend API Stores Data

- The API validates and stores dependencies in the database.
- API example:

```
POST /api/track
Content-Type: application/json
{
  "service_1": "AuthService",
  "service_2": "UserService",
  "method": "CALL"
}
```

- Our backend assigns timestamps for tracking historical changes.
- Stored in → PostgreSQL

## Step 3️⃣ - Frontend Fetches Dependency Data

- The frontend requests stored dependencies via an API call:

```
GET /api/dependencies
```

- The backend returns a list of microservice interactions in JSON format:

```json
[
  {
    "service_1": "AuthService",
    "service_2": "UserService",
    "method": "CALL",
    "timestamp": "2024-03-30T12:00:00Z"
  },
  {
    "service_1": "UserService",
    "service_2": "PaymentService",
    "method": "CALL",
    "timestamp": "2024-03-30T12:05:00Z"
  }
]
```

## Step 4️⃣ - Visualization using D3.js Graph

- The frontend renders microservice dependencies as a graph.
- Example:
  - AuthService → UserService → PaymentService
  - Shows as a directed graph with arrows indicating the call flow.
