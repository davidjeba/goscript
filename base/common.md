# Common Rules

GoScript is a guided language for AI-era systems. It is not a framework clone and it is not an unguided language where every agent invents a new project shape.

## Required Decisions

Every AI coder must decide these before generating a project:

1. Mode: `cs` or `sw`
2. Project type: `website`, `app`, or `erp`
3. Which `base/config/*.md` files apply
4. Whether runtime autonomous agents are required

## Shared Project Shape

Use this baseline tree unless a mode or project type explicitly extends it:

```text
project/
  manifest.json
  base/
  agents/
  app/
    modules/
    pages/
    components/
    services/
    routes/
    assets/
  core/
  tests/
  docs/
  deploy/
```

## Global Rules

- Keep naming predictable and reusable across projects
- Prefer modules over one-off feature folders
- Keep `core/` outside normal application execution boundaries
- Treat `base/` as build-time guidance and `agents/` as runtime autonomous roles
- Keep the folder tree as uniform as possible across `website`, `app`, and `erp`

