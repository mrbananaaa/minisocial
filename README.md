A collaboration project between me and GPT. This project is one of my journey of learning go(lang) and this time we’ll focus on the architecture. I’m learning something called Modular Monolith but this project is combination about a lot of buzzword bullshit like Clean Architecture, Hexagonal, DDD, CQRS, and any other shit that obviously I don’t like until it solve the problem. So the point is I’m trying to learn the process of making an app with Modular Monolith approach with the hope that it’s gonna be more easier to swap or split one of the modules or domain to microservice in the future when the app start to grow.

The app are simply Medium or Dev.to clone, so the user can create account and start writing, or just come to read an articles.

## Tech Stack

- Go 1.26
- go-chi/chi (router)
- Postgresql
- SQLC (type safe query)
- pressly/goose (migration tools)
- The ungodly slog
- NATS Jetstream

## Rule Books

- Modules own business capabilities
- Consumer own interfaces
- No abstraction until it hurts
- Domain knows nothing about infrastructure
- Workflows coordinate modules
- Events represent business facts

## Current Modules

- User
- Post

## Folder Structure

```
cmd/
internal/
├─ post/
│  ├─ application/
│  ├─ domain/
│  ├─ infra/
│  ├─ repositories/
│  ├─ api/
│  ├─ module.go
├─ user/
│  ├─ application/
│  ├─ domain/
│  ├─ infra/
│  ├─ repositories/
│  ├─ api/
│  ├─ module.go
├─ workflows/
│  ├─ create_post.go
│  ├─ delete_user.go
│  ├─ ...
├─ platform/
│  ├─ validation/
│  ├─ logger/
│  ├─ config/
│  ├─ db/
│  ├─ ...
```

## Planning

### Phase I - Foundation

- Episode 1 - Bootstrapping the project
  - Folder structure
  - cmd/server
  - internal
  - platform
  - Docker
  - PostgreSQL
  - Goose
  - SQLC
  - NATS
  - Configuration
  - Logger
  - main.go
  - Composition Root
- Episode 2 - Building the User Module
  - module.go
  - repositories
  - application
  - domain
  - infra
  - api
- Episode 3 - Building the Post Module
  - Create Post
  - Edit Post
  - Archive Post

### Phase II - Connecting Modules

- Episode 4 - Cross Module Workflow
  - User creates a post
  - workflows/create_post.go
- Episode 5 - Transaction
  - Transaction Manager
  - Unit of Work
  - SQLC integration
  - Repository participation
  - Cross-module transaction
- Episode 6 - Domain Events
  - Not NATS yet, pure domain events
- Episode 7 - Event Dispatcher
  - Application → Dispatcher
  - Dispatcher → Handlers
- Episode 8 - Refactor Dispatcher
  - Replace Dispatcher with NATS JetStream
  - Expect zero business changes

### Phase III - Growing the System

- Episode 9 - Notification Module
  - Consumes PostCreated
- Episode 10 - CQRS-lite
  - Separate write model from read model
- Episode 11 - Outbox Pattern
  - The reason it exists
  - Why publishing directly to NATS inside a transaction is dangerous
- Episode 12 - Background Workers
  - NATS Consumers
  - Retries
  - Dead Letters

### Phase IV - Real Production

- Episode 13 - Authentication
  - JWT
  - Sessions
  - Refresh Tokens
- Episode 14 - Authorization
  - Policies
  - Ownership
- Episode 15 - Observability
  - Logging
  - Tracing
  - Metrics
  - Request IDs
- Episode 16 - Split one module into a microservice
