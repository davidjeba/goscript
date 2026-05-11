# ERP Type

Use this type for ERP systems, admin-heavy business software, operations platforms, and modular enterprise applications.

## Shared Focus

- Strong module boundaries
- Predictable business workflows
- Clear role-based surfaces
- Protected `core/` environment

## Typical Folders

```text
project/
  app/
    modules/
      task/
      pos/
      pay/
      hr/
      crm/
    services/
    routes/
  core/
  deploy/
```

## Notes

- `erp` may run in either `cs` or `sw`
- `sw` is often the stronger long-term fit for modular enterprise systems
- Never let a module directly own `core/`

