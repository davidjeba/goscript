# Base Guidance

`base/` is the build-time guidance layer for AI coders working in GoScript.

It is not a runtime framework folder. It defines how AI should shape projects, choose `cs` or `sw`, apply project-type rules, and keep the ecosystem consistent while the language is still young.

## Load Order

1. `base/common.md`
2. `base/modes/cs.md` or `base/modes/sw.md`
3. `base/types/website.md`, `base/types/app.md`, or `base/types/erp.md`
4. Selected files under `base/config/`
5. Selected files under `base/policies/`
6. Project-local overrides

## Purpose

- Keep project structure about 90 percent uniform across the ecosystem
- Let AI coders follow explicit generation rules instead of inventing ad hoc layouts
- Separate build-time guidance from runtime autonomous agents
- Keep GoScript language-driven while still being a guided language

## Layout

```text
base/
  common.md
  modes/
    cs.md
    sw.md
  types/
    website.md
    app.md
    erp.md
  config/
    ui.md
    generative.md
    security.md
    data.md
    erp.md
  policies/
    core.md
    swarm-strict.md
```

