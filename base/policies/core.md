# Core Policy

`core/` is outside the normal application runtime boundary.

## Rules

- Modules cannot import, browse, mount, or directly execute inside `core/`
- `core/` is only reachable through approved admin entrypoints
- Admin access should use verified domains such as `admin.domain.com` or approved subdomains
- CNAME verification helps prove domain control, but domain verification alone is not the trust anchor
- Treat `core/` as a protected environment for config, secrets, authority data, and privileged DB access

