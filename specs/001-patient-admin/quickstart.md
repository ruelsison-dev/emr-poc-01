# Quickstart — Patient Administration Sub-System (developer)

**Goal**: Get services and tests running locally for development and integration testing.

## Prerequisites
- Docker & Docker Compose
- Node.js 20.x (LTS) and npm/yarn
- TypeScript (dev dependency)
- Python 3.11 (optional if running legacy FastAPI services)
- Terraform (for infra plan checks) and `conftest` (for OPA checks)

## Local services (recommended)
- PostgreSQL (docker)
- Redis (docker)
- Localstack (optional) for S3 and SQS emulation in dev

## Example steps
1. Start dependencies:
   - docker run -d --name pg -e POSTGRES_PASSWORD=pass -p 5432:5432 postgres:15
   - docker run -d --name redis -p 6379:6379 redis:7
2. Create `.env` from `.env.example` and set DB, Redis, and IdP settings
3. Start services:
   - Backend (Express API):
     - cd backend/express-api
     - npm install
     - npm run dev   # starts the TypeScript dev server with ts-node-dev or nodemon
   - (Optional) Legacy Python services remain available:
     - cd services/registry-service
     - pip install -r requirements.txt
     - uvicorn registry.main:app --reload --port 8001
4. Run tests:
   - Backend (Node):
     - cd backend/express-api
     - npm test        # runs jest + supertest
   - (Optional) Python tests:
     - pytest tests/unit
     - pytest tests/integration (requires local infra running)
5. IaC checks:
   - terraform init && terraform validate (inside infra/)
   - terraform plan -out plan.out
   - conftest test plan.out
6. Load tests:
   - k6 run load_test_scripts/booking_test.js

## Notes
- Use Localstack to reduce reliance on real AWS services for development, but verify end-to-end in staging with real AWS services before release (especially for PHI handling).
- Follow the Consent and Notification policy when testing outbound messages: use mock delivery channels where possible.

## Notes
- Use Localstack to reduce reliance on real AWS services for development, but verify end-to-end in staging with real AWS services before release (especially for PHI handling).
- Follow the Consent and Notification policy when testing outbound messages: use mock delivery channels where possible.

