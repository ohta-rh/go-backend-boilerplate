# Go Backend Boilerplate

A Go backend application boilerplate adopting Clean Architecture.

## Architecture

This project is designed according to the principles of [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

![Clean Architecture](https://blog.cleancoder.com/uncle-bob/images/2012-08-13-the-clean-architecture/CleanArchitecture.jpg)

### Project Architecture Diagram

```mermaid
graph TD
    subgraph "External"
        Web["Web / API"]
        DB["Database"]
    end
    
    subgraph "Interface Adapters"
        H["Handlers<br/>(internal/handler)"]
        R["Repositories<br/>(internal/repository)"]
    end
    
    subgraph "Application Core"
        UC["Use Cases<br/>(internal/usecase)"]
        D["Domain<br/>(internal/domain)"]
    end
    
    subgraph "Infrastructure"
        I["Infrastructure<br/>(internal/infrastructure)"]
    end
    
    Web --> H
    H --> UC
    UC --> D
    R --> D
    R --> I
    I --> DB
    
    classDef core fill:#f9f,stroke:#333,stroke-width:2px;
    classDef adapter fill:#bbf,stroke:#333,stroke-width:1px;
    classDef external fill:#bfb,stroke:#333,stroke-width:1px;
    classDef infra fill:#fbb,stroke:#333,stroke-width:1px;
    
    class D,UC core;
    class H,R adapter;
    class Web,DB external;
    class I infra;
```

### Layer Structure

- **Domain Layer** (`internal/domain/`): Contains business entities and repository interfaces. This layer represents the core business concepts and rules. It does not depend on any other layer.
- **Schema Layer** (`internal/schema/`): Contains request and response schemas for API communication. This layer is responsible for data transfer objects (DTOs) and conversion between domain entities and API schemas.
- **Use Case Layer** (`internal/usecase/`): Implements the business logic of the application. Depends only on the domain layer.
- **Interface Adapter Layer**:
  - **Handler** (`internal/handler/`): Processes HTTP requests, converts between schemas and domain entities, and calls use cases.
  - **Repository** (`internal/repository/`): Provides implementations of domain repository interfaces.
- **Infrastructure Layer** (`internal/infrastructure/`): Responsible for integration with external services such as database connections.

## Project Structure

```
.
├── .github/                # GitHub specific files
│   └── workflows/          # GitHub Actions workflows
│       └── ci.yaml         # Continuous Integration workflow
├── cmd/                    # Application entry points
│   ├── api/                # API server
│   └── entgen/             # Ent schema generation tool
├── ent/                    # Ent framework related code
│   └── schema/             # Database schema definitions
├── internal/               # Internal packages
│   ├── domain/             # Domain layer (business entities)
│   │   └── user/           # User domain entities
│   ├── handler/            # HTTP handlers
│   ├── infrastructure/     # Infrastructure layer
│   ├── schema/             # Request/Response schemas
│   │   └── database/       # Database connections
│   ├── repository/         # Repository implementations
│   └── usecase/            # Use cases (application logic)
├── .air.toml               # Air configuration file (hot reload)
├── .clinerules             # Cline editor dependency rules
├── .cursorrules            # Cursor editor dependency rules
├── .golangci.yml           # GolangCI-Lint configuration
├── compose.yml             # Docker Compose configuration
├── Dockerfile              # Docker image build configuration
├── go.mod                  # Go module definition
└── Makefile                # Build commands
```

## Technology Stack

- [Go](https://golang.org/) - Programming language
- [Ent](https://entgo.io/) - Entity framework (ORM)
- [Docker](https://www.docker.com/) - Containerization
- [Air](https://github.com/cosmtrek/air) - Hot reload development tool
- [GolangCI-Lint](https://golangci-lint.run/) - Go linters aggregator
- [GitHub Actions](https://github.com/features/actions) - CI/CD platform

## Setup

### Prerequisites

- Go 1.20 or higher
- Docker & Docker Compose
- Make

### Development Environment Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/easy-go-backend.git
   cd easy-go-backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Start Docker containers:
   ```bash
   make up
   ```

4. Start development server:
   ```bash
   make dev
   ```

## Development Guidelines

### Development Workflow

```mermaid
graph TD
    Start[Start Development] --> Feature[Create Feature Branch]
    Feature --> Implement[Implement Feature]
    Implement --> Test[Write Tests]
    Test --> Lint[Run Linting]
    Lint --> Fix[Fix Issues]
    Fix --> PR[Create Pull Request]
    PR --> Review[Code Review]
    Review --> Merge[Merge to Main]
    
    classDef start fill:#bfb,stroke:#333,stroke-width:1px;
    classDef process fill:#bbf,stroke:#333,stroke-width:1px;
    classDef end fill:#fbb,stroke:#333,stroke-width:1px;
    
    class Start start;
    class Implement,Test,Lint,Fix,PR,Review process;
    class Merge end;
```

### Adding a New Entity

```mermaid
flowchart TD
    A[1. Create Ent Schema] --> B[2. Generate Ent Code]
    B --> C[3. Add Domain Entity]
    C --> D[4. Create API Schemas]
    D --> E[5. Add Repository Implementation]
    E --> F[6. Add Use Case]
    F --> G[7. Add Handler]
    G --> H[8. Add Route]
    
    classDef schema fill:#f9f,stroke:#333,stroke-width:1px;
    classDef domain fill:#bbf,stroke:#333,stroke-width:1px;
    classDef impl fill:#bfb,stroke:#333,stroke-width:1px;
    classDef api fill:#fbb,stroke:#333,stroke-width:1px;
    
    class A,B schema;
    class C domain;
    class D,G,H api;
    class E,F impl;
```

The steps in detail:

1. Create a new schema file in the `ent/schema/` directory
2. Run `cmd/entgen/entgen.go` to generate Ent code
3. Add a new domain entity and repository interface in `internal/domain/`
4. Create request/response schemas in `internal/schema/`
5. Add repository implementation in `internal/repository/`
6. Add use case in `internal/usecase/`
7. Add handler in `internal/handler/` (with conversion between schemas and domain entities)
8. Add route in `internal/handler/router.go`

### Testing

The project includes automated tests for all components except the `ent` directory, which is excluded from test runs. Each layer of the architecture has its own tests:

```mermaid
graph TD
    subgraph "Test Coverage"
        Domain["Domain Layer Tests"]
        UseCase["Use Case Layer Tests"]
        Repository["Repository Layer Tests"]
        Handler["Handler Layer Tests"]
        Schema["Schema Layer Tests"]
        Infra["Infrastructure Layer Tests"]
    end
    
    Domain --> |Validates| DomainEntities["Business Entities<br/>Business Rules"]
    UseCase --> |Validates| BusinessLogic["Application Logic<br/>Use Case Flows"]
    Repository --> |Validates| DataAccess["Data Access<br/>Storage Operations"]
    Handler --> |Validates| API["API Endpoints<br/>Request Handling"]
    Schema --> |Validates| DataTransformation["DTO Conversion<br/>Validation"]
    Infra --> |Validates| ExternalServices["External Integrations<br/>Database Connections"]
    
    classDef domain fill:#f9f,stroke:#333,stroke-width:1px;
    classDef usecase fill:#bbf,stroke:#333,stroke-width:1px;
    classDef repo fill:#bfb,stroke:#333,stroke-width:1px;
    classDef handler fill:#fbb,stroke:#333,stroke-width:1px;
    classDef schema fill:#fbf,stroke:#333,stroke-width:1px;
    classDef infra fill:#bff,stroke:#333,stroke-width:1px;
    
    class Domain,DomainEntities domain;
    class UseCase,BusinessLogic usecase;
    class Repository,DataAccess repo;
    class Handler,API handler;
    class Schema,DataTransformation schema;
    class Infra,ExternalServices infra;
```

- **Domain Layer**: Tests for business entities and their behavior
- **Use Case Layer**: Tests for application logic and business rules
- **Repository Layer**: Tests for data access implementations
- **Handler Layer**: Tests for HTTP request handling
- **Schema Layer**: Tests for data transformation between API and domain models
- **Infrastructure Layer**: Tests for external service integrations

Run tests with:

```bash
make test       # Run tests (excluding ent directory)
make test.cover # Run tests with coverage (excluding ent directory)
make test.pkg   # Run tests for a specific package
```

### Linting

The project uses GolangCI-Lint for code quality checks and static analysis. This helps maintain code quality by detecting:

- Unused imports and variables
- Potential nil pointer dereferences
- Code style issues
- Possible bugs and anti-patterns
- And many other code quality issues

Run linting with:

```bash
make lint       # Run linters and automatically fix issues where possible
make api.lint.check  # Run linters without auto-fixing (useful for CI)
```

When adding new code, always ensure it passes both tests and linting checks before committing.

### CI/CD

The project includes GitHub Actions workflows for continuous integration. The CI pipeline runs tests and linting on each pull request and push to the main branch.

```mermaid
flowchart LR
    PR[Pull Request] --> Build[Build]
    Build --> Lint[Lint Code]
    Lint --> Test[Run Tests]
    Test --> Coverage[Check Coverage]
    Coverage --> Report[Generate Report]
    
    Push[Push to Main] --> BuildProd[Build Production]
    BuildProd --> TestProd[Test Production]
    TestProd --> Deploy[Deploy]
    
    classDef trigger fill:#f9f,stroke:#333,stroke-width:1px;
    classDef process fill:#bbf,stroke:#333,stroke-width:1px;
    classDef deploy fill:#bfb,stroke:#333,stroke-width:1px;
    
    class PR,Push trigger;
    class Build,Lint,Test,Coverage,Report,BuildProd,TestProd process;
    class Deploy deploy;
```

### Building

```bash
make build
```

## Dependency Rules

This project uses both `.cursorrules` and `.clinerules` files to strictly enforce Clean Architecture dependency rules. These files prevent inner layers from depending on outer layers.

### Dependency Flow Diagram

```mermaid
flowchart TD
    subgraph "Dependency Direction"
        direction LR
        Outer["Outer Layers"] --> Inner["Inner Layers"]
    end
    
    Domain["Domain Layer<br/>(internal/domain)"]
    UseCase["Use Case Layer<br/>(internal/usecase)"]
    Handler["Handler Layer<br/>(internal/handler)"]
    Repository["Repository Layer<br/>(internal/repository)"]
    Infrastructure["Infrastructure Layer<br/>(internal/infrastructure)"]
    Schema["Schema Layer<br/>(internal/schema)"]
    
    Handler --> UseCase
    Handler --> Domain
    Handler --> Schema
    UseCase --> Domain
    Repository --> Domain
    Repository --> Infrastructure
    Infrastructure --> Domain
    Schema --> Domain
    
    classDef core fill:#f9f,stroke:#333,stroke-width:2px;
    classDef outer fill:#bbf,stroke:#333,stroke-width:1px;
    
    class Domain core;
    class UseCase,Handler,Repository,Infrastructure,Schema outer;
```

The dependency rules are:

- Domain layer should not depend on any other layer.
- Use case layer should only depend on domain layer.
- Repository layer should only depend on domain and infrastructure layers.
- Handler layer should only depend on domain and use case layers.
- Infrastructure layer should only depend on domain layer.

## License

MIT