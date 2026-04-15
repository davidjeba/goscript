# GoScript Task Manager

A full-featured test webapp built with [GoScript 2.0](https://github.com/davidjeba/goscript) — the Go-native web framework.

## What it tests

| Feature | How |
|---|---|
| **AppRouter** | 4 routes: `/`, `/projects`, `/projects/:project` (dynamic), `/about` |
| **Reactive gs-* attributes** | Add task (gs-post), toggle done (gs-get/outerHTML), delete (gs-get/outerHTML) |
| **Middleware pipeline** | Gzip + CORS + Security Headers + Logging + Recovery + Request ID |
| **Metadata / SEO** | `NewMetadata()` builder on every page with title, description, theme-color |
| **Server Components** | `shell()` layout wrapper composing page bodies |
| **HTML Fragment API** | POST returns a `<tr>` row; toggle/delete return partial HTML |
| **JSON API** | `/api/health` returns JSON with framework version and task stats |
| **Embedded runtime.js** | `gs.RuntimeHandler()` serving the 42KB reactive runtime |

## Running

```bash
git clone https://github.com/davidjeba/goscript.git
cd goscript-taskapp

# Edit go.mod — uncomment the replace directive pointing to your clone:
# replace github.com/davidjeba/goscript => ../goscript

go mod tidy
go run main.go
# → http://localhost:8080
```

## Routes

```
GET  /                      Task board (reactive add/toggle/delete)
GET  /projects              Project list (dynamic cards)
GET  /projects/:project     Project detail (dynamic segment)
GET  /about                 Framework feature showcase

POST /api/tasks             Add task → returns <tr> fragment
GET  /api/tasks/:id/toggle  Toggle done → returns updated <tr>
GET  /api/tasks/:id/delete  Delete task → returns empty (outerHTML removes row)
GET  /api/stats             Stats bar fragment
GET  /api/health            JSON health check
GET  /__goscript/runtime.js Embedded reactive runtime
```

## Binary size

~7.6MB static binary. Zero runtime dependencies.
