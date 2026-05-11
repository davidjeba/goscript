# App Type

Use this type for product applications, dashboards, productivity tools, and software that behaves more like an interactive product than a content website.

## Shared Focus

- Reusable product modules
- State-heavy interfaces
- Service-oriented app flows
- Strong internal navigation

## Typical Folders

```text
project/
  app/
    modules/
    pages/
    components/
    services/
    routes/
```

## Notes

- `app` may run in either `cs` or `sw`
- Prefer explicit service contracts between modules
- Keep product actions and UI flows easy for AI coders to extend later

