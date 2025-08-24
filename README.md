########## Run with Docker Compose
go mod tidy
docker compose up --build

########## ENDPOINTS : 
API on http://localhost:8080
- POST /orders
  - body: {"product_id":1,"quantity":1,"buyer_id":"user-123"}
- GET /orders/:id
- POST /jobs/settlement
  - body: {"from":"2025-08-01","to":"2025-08-30"}
- GET /jobs/:id
- POST /jobs/:id/cancel
- GET /downloads/<jobId>.csv (served from /tmp/settlements)