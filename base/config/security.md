# Security Config

Use this file when the project has trust boundaries, external integrations, or sensitive workflows.

## Rules

- Keep `core/` out of bounds for normal app runtime
- Prefer explicit capability grants over unrestricted access
- Treat module-to-service communication as a contract, not a shortcut
- Use stronger review for payment, identity, audit, and admin surfaces

