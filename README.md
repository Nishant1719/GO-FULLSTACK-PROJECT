## ADR-001: Layered Web Application Architecture with BFF and Domain Services
- Status : Nil
- Date : 21-01-2026

# Context
- The system requires:
    - A clean separation between UI, routing/authentication, and business logic
    - Secure, stateless authentication
    - High performance backend services
    - Strong data consistency
    - Ability to scale services independently
    - Long-term maintainability
    - To meet these requirements, we evaluated multiple backend architecture patterns and selected a layered architecture with a Backend-for-Frontend (BFF).

# Decision
- We adopt the following architecture:
```
React.js → Node.js (BFF) → Go (DDD) → PostgreSQL
```
Each layer has strict responsibilities and communicates only with the adjacent layer.

# Layer Responsibilities

1. React.js (Presentation Layer)

- Purpose: User Interface and User Interaction
- Responsibilities:
    - UI rendering and state management
    - Client-side routing
    - Basic input validation
    - API consumption
- Constraints:
    - Must call only Node.js APIs
    - Must not contain business rules
    - Must not handle authentication logic

2. Node.js (Backend-for-Frontend Layer)

- Purpose: System boundary and request orchestration
- Responsibilities:
    - Acts as the single entry point for frontend
    - JWT authentication and authorization
    - Request validation and transformation
    - API routing to Go services
    - Response shaping for frontend needs
- Key Design Choice:
    - JWT validation occurs only in Node.js
    - Go services trust identity forwarded by Node.js
- Constraints:
    - Must not implement business rules
    - Must not access the database

3. Go (Domain / Business Logic Layer)

- Purpose: Core business logic and data integrity
- Responsibilities:
    - Enforce business rules
    - Validate domain data
    - Handle transactions
    - Interact with PostgreSQL
    - Maintain domain consistency
- Trust Model
    - Assumes requests are authenticated
    - Consumes validated user context from Node.js
- Constraints
    - Must not perform JWT validation
    - Must not contain frontend-facing HTTP or UI-specific logic.

4. PostgreSQL (Data Layer)

- Purpose: Persistent data storage
- Responsibilities:
    - Data persistence
    - Referential integrity
    - Constraints and indexing
    - Transaction support
- Constraints
    - Accessed only by Go services
    - No direct access from Node.js or React.js

# Authentication & Authorization Strategy

- ## Authentication

- Authentication is handled at the Node.js (BFF) layer
- The BFF is responsible for:
    - Verifying the identity of incoming requests
    - Ensuring requests originate from authenticated clients
    - Rejecting unauthenticated or malformed requests
- Authentication details are abstracted away from downstream services
- The specific authentication mechanism (e.g., token-based, session-based, or external identity provider) is an implementation detail of the BFF and is not coupled to the domain layer.

- ## Authorization
- Request-level authorization (e.g., role or access checks) is enforced at the BFF layer.
- Domain-level authorization (business rule–driven access constraints) is enforced within Go services where required.
- Authorization decisions are based on validated identity context propagated from the BFF

- ## Identity Propagation

- Once a request is authenticated, the BFF forwards a validated identity context to the Go services
- Go services trust this context and do not perform authentication themselves
- Identity context is used only for domain decisions, not for security boundary enforcement

- ## Rationale

- Authentication is a system boundary concern, not a domain concern
- Keeping authentication logic in the BFF:
- Preserves clean separation of responsibilities
- Prevents security logic from leaking into business code
- Reduces duplication across domain services
- Go services remain focused on business correctness and data integrity

# Consequences

- ## Positive Outcomes
- Clear separation of concerns
- Reduced attack surface
- Simplified business logic
- Independent scaling
- Easier testing and debugging
- ## Trade-offs
- Additional network hop (Node → Go)
- Increased operational complexity
- Requires strict API contracts

# Final Decision Statement
