# CS Mode

`cs` means client-server mode.

Use this mode when the application has one primary server that owns the backend contract for the system.

## Characteristics

- One main server owns the application state and service contract
- Clients consume pages, APIs, and sessions from that server
- Deployment stays centralized and simple
- Best for standard websites, internal tools, dashboards, and smaller ERP deployments

## Folder Extensions

`cs` keeps the shared tree and may add these folders when needed:

```text
project/
  app/
    api/
    controllers/
    views/
```

## Rules

- Do not add swarm-only folders unless the project later moves to `sw`
- Keep routing and service ownership centralized
- Use `manifest.mode = "cs"`

