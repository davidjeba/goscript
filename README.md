# GoScript

GoScript is a Go-native language and web application runtime that replaces JavaScript for teams that want one language, one toolchain, and one deployment story.

It follows a familiar app-structure that feels approachable to people used to modern web frameworks, but the point is different: GoScript is not a Next.js replacement. It is a language-first alternative to JavaScript for Go developers who want to own the full stack in Go.

## Why it matters

JavaScript made the web programmable, but it also created a long tail of package churn, runtime drift, and extra mental context for teams that already live in Go.

GoScript changes that tradeoff:

- Build web apps in Go instead of splitting the product across Go and JavaScript.
- Keep the same language from API to UI, which lowers context switching and improves ownership.
- Compile to a predictable Go-native deliverable, instead of shipping a second language ecosystem to production.
- Stay aligned with the Go community’s bias for simplicity, explicitness, and operational clarity.

## GoScript vs JavaScript

| Area | GoScript | JavaScript |
|---|---|---|
| Language model | Go-native, language-first | Browser-native, ecosystem-first |
| Team workflow | One language across backend and UI | Often split across Go plus JS/TS |
| Runtime feel | Explicit, compiled, predictable | Dynamic, flexible, highly distributed |
| Deployment story | Go-shaped and operationally simple | Node/browser ecosystem and its dependency graph |
| Productivity tradeoff | Less context switching for Go teams | Massive ecosystem, but more seams to manage |
| Long-term ownership | Easier for Go teams to keep inside one stack | Best when the team wants JavaScript's native world |

## Why the Go community should care

GoScript is compelling because it gives Go developers a native path to the web without forcing them to abandon the language that already powers their services, tools, and infrastructure.

That makes it feel like a native blessing for the Go ecosystem:

- Go stays the source of truth, from business logic to interface logic.
- Teams can move faster without asking everyone to become fluent in a second language.
- The app model stays familiar, but the ownership model becomes more Go-like.
- Web development becomes a first-class Go workload instead of a parallel JavaScript world.

## The short version

JavaScript will still remain the language of the browser. GoScript's value is different: it lets the Go community build modern web apps without leaving Go, and that is exactly why it deserves attention.

If you want a quick place to see the shape of the project, start with the example apps in `examples/`.
