# 🐳 Docker Guide - Library Management System

## Quick Start

### 1. Build and Start Everything
```bash
# Build images and start containers in foreground
docker-compose up --build

# Or run in background (detached mode)
docker-compose up -d --build
```

### 2. Check Status
```bash
# View running containers
docker-compose ps

# View logs (all services)
docker-compose logs

# View logs for specific service
docker-compose logs app
docker-compose logs db

# Follow logs in real-time
docker-compose logs -f app
```

### 3. Test the API
```bash
# Health check
curl http://localhost:8080/health

# Create a book
curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "1984",
    "author": "George Orwell",
    "quantity": 5
  }'

# Get all books
curl http://localhost:8080/api/v1/books

# Get book by ID
curl http://localhost:8080/api/v1/books/1

# Update book
curl -X PUT http://localhost:8080/api/v1/books/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "1984 (Updated)",
    "author": "George Orwell",
    "quantity": 10,
    "available": 8
  }'

# Delete book
curl -X DELETE http://localhost:8080/api/v1/books/1
```

### 4. Stop and Clean Up
```bash
# Stop containers (keeps data)
docker-compose stop

# Stop and remove containers (keeps volumes/data)
docker-compose down

# Stop, remove containers AND volumes (deletes all data!)
docker-compose down -v
```

---

## Common Commands

### Container Management
```bash
# Start existing containers
docker-compose start

# Stop containers
docker-compose stop

# Restart containers
docker-compose restart

# Restart specific service
docker-compose restart app
```

### Rebuilding
```bash
# Rebuild after code changes
docker-compose up --build

# Rebuild specific service
docker-compose build app
docker-compose up -d app

# Force rebuild (no cache)
docker-compose build --no-cache
```

### Logs and Debugging
```bash
# View all logs
docker-compose logs

# Last 100 lines
docker-compose logs --tail=100

# Follow logs
docker-compose logs -f

# Logs for specific service
docker-compose logs app
docker-compose logs db
```

### Execute Commands in Container
```bash
# Open shell in app container
docker-compose exec app sh

# Open shell in database container
docker-compose exec db sh

# Run psql in database
docker-compose exec db psql -U admin -d library_db

# Check database tables
docker-compose exec db psql -U admin -d library_db -c "\dt"
```

---

## Troubleshooting

### Problem: Port Already in Use
```bash
# Error: Bind for 0.0.0.0:8080 failed: port is already allocated

# Solution 1: Stop the process using the port
lsof -ti:8080 | xargs kill -9

# Solution 2: Change port in docker-compose.yml
ports:
  - "8081:8080"  # Use 8081 on host instead
```

### Problem: Database Connection Failed
```bash
# Check if database is healthy
docker-compose ps

# View database logs
docker-compose logs db

# Restart database
docker-compose restart db
```

### Problem: Code Changes Not Reflected
```bash
# Rebuild the image
docker-compose up --build

# Or force rebuild without cache
docker-compose build --no-cache app
docker-compose up -d
```

### Problem: Permission Denied
```bash
# Make sure Docker daemon is running
docker ps

# On Linux, add user to docker group
sudo usermod -aG docker $USER
# Then log out and log back in
```

---

## Development Workflow

### 1. Start Development Environment
```bash
docker-compose up -d
docker-compose logs -f app
```

### 2. Make Code Changes
- Edit your Go files
- Save changes

### 3. Rebuild and Test
```bash
docker-compose up --build -d
curl http://localhost:8080/health
```

### 4. View Logs
```bash
docker-compose logs -f app
```

### 5. Stop When Done
```bash
docker-compose down
```

---

## Production Considerations

### 1. Use Environment Variables
Don't hardcode secrets in docker-compose.yml. Use `.env` file or pass via command line:

```bash
DB_PASSWORD=secure_password docker-compose up -d
```

### 2. Enable SSL for Database
Update docker-compose.yml:
```yaml
environment:
  DB_SSLMODE: require  # or verify-full
```

### 3. Use Docker Secrets (Docker Swarm)
For production deployments with Docker Swarm.

### 4. Resource Limits
Add resource constraints:
```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
```

---

## What Each File Does

| File | Purpose |
|------|---------|
| **Dockerfile** | Recipe for building Go app image (multi-stage build) |
| **docker-compose.yml** | Orchestrates app + database containers |
| **.dockerignore** | Excludes files from Docker build (faster builds) |

---

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                    Host Machine                      │
│                                                      │
│  ┌────────────────────────────────────────────────┐ │
│  │         Docker Compose                         │ │
│  │                                                │ │
│  │  ┌──────────────────────────────────────────┐ │ │
│  │  │      library-network (bridge)            │ │ │
│  │  │                                          │ │ │
│  │  │  ┌─────────────┐    ┌─────────────┐    │ │ │
│  │  │  │     app     │    │     db      │    │ │ │
│  │  │  │ (Go API)    │───▶│ (PostgreSQL)│    │ │ │
│  │  │  │             │    │             │    │ │ │
│  │  │  │ Port: 8080  │    │ Port: 5432  │    │ │ │
│  │  │  └─────────────┘    └──────┬──────┘    │ │ │
│  │  │                             │           │ │ │
│  │  └─────────────────────────────┼───────────┘ │ │
│  │                                │             │ │
│  │                        ┌───────▼────────┐   │ │
│  │                        │ postgres_data  │   │ │
│  │                        │   (volume)     │   │ │
│  │                        └────────────────┘   │ │
│  └────────────────────────────────────────────┘ │
│                                                  │
│  Access: http://localhost:8080                   │
└──────────────────────────────────────────────────┘
```

