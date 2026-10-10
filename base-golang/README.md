# Base Golang

> Rest API **Fiber v3 + PostgreSQL**

```text
Backend
├── Config / env
├── Routing
├── Middleware
│   ├── Logging
│   ├── Recovery
│   ├── CORS
│   └── Rate Limit
├── No Auth now
├── Validation
├── Error Handling
├── Context + Timeout
├── Database
│   ├── Connection Pool
│   ├── Transaction
│   └── Migration
├── Health Check
├── Graceful Shutdown
├── Structured Logging + Request ID
├── Testing
└── Docker Compose for PostgreSQL
```

## Architecture

package by feature and each module has a separate function:

```text
HTTP → Handler → Service → Repository → PostgreSQL
                    ↓
                  Domain
```

```text
cmd/api/                    entry point of API
cmd/migrate/                run migration separately
internal/app/               wiring, routes, startup/shutdown
internal/modules/note/
  delivery/http/            handlers and DTOs
  service/                  use cases and repository interface
  domain/                   model and validation
  infrastructure/           PostgreSQL repository
internal/platform/          HTTP server, middleware, logger, database
internal/shared/apperror/   application errors
migrations/                 embedded SQL and migration runner
```

Add new feature by using note module as an example, then wiring in `internal/app/app.go` and adding routes in `routes.go`
