# MiniSocial

A collaboration project between me and GPT. This project is one of my journey of learning go(lang) and this time we’ll focus on the architecture. I’m learning something called Modular Monolith but this project is a combination about a lot of buzzword bullshit like Clean Architecture, Hexagonal, DDD, CQRS, and any other shit that obviously I don’t like until it solve the problem. So the point is I’m trying to learn the process of making an app with Modular Monolith approach with the hope that it’s gonna be more easier to swap or split one of the modules or domain to microservice in the future when the app start to grow.

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

## Architecture Conventions

### Modules

- Each business module owns its own domain, application, repositories, infrastructure, and transport.
- module.go is responsible only for constructing the module's internal dependencies.
- Modules expose only their application service through Service().
- Modules do not know about other modules.

### Workflows

- Workflows orchestrate multiple modules to implement a business use case.
- Each workflow exposes a single Execute(ctx, input) entry point.
- Workflows depend only on module application services, never repositories or infrastructure.
- Workflows are composed only in the application composition root.

### HTTP / Transport

- HTTP handlers are transport adapters only.
- Handlers may call either a module application service or a workflow depending on the use case.
- The transport layer is allowed to import workflow packages directly when invoking a workflow.
- Request/response mapping belongs in handlers, not in workflows or application services.

### Composition Root

- The composition root owns dependency injection.
- The composition root is the only place allowed to know every module and workflow.
- Routing is centralized inside internal/app/router.go.
- Public HTTP endpoints are organized by resources (/users, /posts, etc.), not by internal module structure.

### Interfaces

- The consumer owns interfaces.
- Prefer concrete types by default.
- Introduce interfaces only when they provide a real benefit (multiple implementations, testing seams, or dependency inversion).
- Do not create interfaces solely to avoid importing another package.

### Dependency Direction

```
HTTP Router
      ↓
HTTP Handlers
      ↓
Workflow or Application Service
      ↓
Application Service
      ↓
Repository
      ↓
Infrastructure
```

Dependencies should always point inward toward the business logic.

### Guiding Principle

When making architectural decisions, prioritize:

1. Clear ownership of responsibilities.
2. Module independence.
3. Simplicity over premature abstraction.
4. Composition over coupling.
5. Refactor only when the current design begins to hurt.

### Architectural Philosophy

Architectural patterns (DDD, Hexagonal Architecture, Clean Architecture, CQRS, etc.) are tools—not goals.

MiniSocial adopts ideas from these patterns only when they solve an existing problem. Every abstraction should justify its existence through simplicity, maintainability, or testability rather than theoretical purity.

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
│  ├─ create_post/
│  ├─ delete_user/
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
