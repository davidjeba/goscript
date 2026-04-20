# Runtime Agents

`agents/` is reserved for autonomous roles that exist inside a GoScript application at runtime.

This folder is not the same as `base/`. `base/` guides AI coders while they build a project. `agents/` defines agentic employees that the built application can hire, configure, and remove.

## Example Layout

```text
agents/
  ceo/
    brain.md
    skills.md
    manifest.json
  cfo/
    brain.md
    skills.md
    manifest.json
```

## Rules

- Use one folder per autonomous role
- Put long-term thinking and mission in `brain.md`
- Put task behavior and role capabilities in `skills.md`
- Put runtime identity, permissions, and lifecycle settings in `manifest.json`

## Templates

Starter templates live under `agents/templates/employee/`.

