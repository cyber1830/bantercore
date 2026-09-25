# Bantercore

Bantercore is a small, deployable Go discussion app: a responsive web UI backed by a JSON API with JWT-protected writes, request IDs, caching, and event publishing seams.

## Run locally

```powershell
go run .
```

Open <http://localhost:8080>.

Run the tests with:

```powershell
go test ./...
```

## Deploy

This repository includes a production Dockerfile and `render.yaml`. In Render, create a new Blueprint from the repository; Render will build the container and expose `/health` as its health check.

The current demo uses in-memory storage so it is intentionally easy to review. For production persistence, replace `MemoryStore` with PostgreSQL and `NoopCache` with Redis while keeping the service interfaces unchanged.