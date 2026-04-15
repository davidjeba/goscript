// GoScript Task Manager — tests:
//   AppRouter (static + dynamic routes), reactive gs-* attributes,
//   middleware pipeline, metadata/SEO builder, server components,
//   API fragment responses, JSON API, embedded runtime.js
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	gs "github.com/davidjeba/goscript/pkg/goscript"
)

// ─── Domain model ────────────────────────────────────────────────────────────

type Priority string

const (
	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

type Task struct {
	ID        int
	Title     string
	Project   string
	Priority  Priority
	Done      bool
	CreatedAt time.Time
}

type Store struct {
	mu      sync.RWMutex
	tasks   []Task
	counter int
}

func NewStore() *Store {
	s := &Store{}
	seeds := []struct {
		title    string
		project  string
		priority Priority
	}{
		{"Build GoScript AppRouter", "Core", PriorityHigh},
		{"Implement reactive gs-* attributes", "Core", PriorityHigh},
		{"Write middleware pipeline", "Infrastructure", PriorityMedium},
		{"Design component model", "Core", PriorityMedium},
		{"Create SSG + ISR engine", "Core", PriorityHigh},
		{"Build gopm CLI tool", "Tooling", PriorityLow},
		{"Write SEO metadata builder", "Features", PriorityLow},
		{"Document framework API", "Docs", PriorityMedium},
	}
	for _, seed := range seeds {
		s.counter++
		s.tasks = append(s.tasks, Task{
			ID: s.counter, Title: seed.title,
			Project: seed.project, Priority: seed.priority,
			CreatedAt: time.Now(),
		})
	}
	return s
}

func (s *Store) All() []Task {
	s.mu.RLock(); defer s.mu.RUnlock()
	out := make([]Task, len(s.tasks)); copy(out, s.tasks); return out
}

func (s *Store) ByProject(p string) []Task {
	s.mu.RLock(); defer s.mu.RUnlock()
	var out []Task
	for _, t := range s.tasks {
		if strings.EqualFold(t.Project, p) { out = append(out, t) }
	}
	return out
}

func (s *Store) Add(title, project string, priority Priority) Task {
	s.mu.Lock(); defer s.mu.Unlock()
	s.counter++
	t := Task{ID: s.counter, Title: title, Project: project,
		Priority: priority, CreatedAt: time.Now()}
	s.tasks = append(s.tasks, t)
	return t
}

func (s *Store) Toggle(id int) (Task, bool) {
	s.mu.Lock(); defer s.mu.Unlock()
	for i, t := range s.tasks {
		if t.ID == id { s.tasks[i].Done = !s.tasks[i].Done; return s.tasks[i], true }
	}
	return Task{}, false
}

func (s *Store) Delete(id int) {
	s.mu.Lock(); defer s.mu.Unlock()
	for i, t := range s.tasks {
		if t.ID == id { s.tasks = append(s.tasks[:i], s.tasks[i+1:]...); return }
	}
}

func (s *Store) Stats() (total, done int) {
	s.mu.RLock(); defer s.mu.RUnlock()
	total = len(s.tasks)
	for _, t := range s.tasks { if t.Done { done++ } }
	return
}

func (s *Store) Projects() []string {
	s.mu.RLock(); defer s.mu.RUnlock()
	seen := map[string]bool{}
	var out []string
	for _, t := range s.tasks {
		if !seen[t.Project] { seen[t.Project] = true; out = append(out, t.Project) }
	}
	return out
}

// ─── HTML rendering helpers ───────────────────────────────────────────────────

var pColor = map[Priority]string{PriorityHigh: "#ef4444", PriorityMedium: "#f59e0b", PriorityLow: "#22c55e"}
var pBg    = map[Priority]string{PriorityHigh: "#fef2f2", PriorityMedium: "#fffbeb", PriorityLow: "#f0fdf4"}

func taskRow(t Task) string {
	doneClass, doneStyle, icon := "", "", "○"
	if t.Done { doneClass = " done"; doneStyle = "text-decoration:line-through;opacity:0.45;"; icon = "✓" }
	c, bg := pColor[t.Priority], pBg[t.Priority]
	return fmt.Sprintf(`<tr id="task-%d" class="task-row%s">
  <td><button class="check-btn" gs-trigger="click" gs-get="/api/tasks/%d/toggle" gs-target="#task-%d" gs-swap="outerHTML">%s</button></td>
  <td style="%s">%s</td>
  <td><span class="proj-tag">%s</span></td>
  <td><span class="pri-badge" style="color:%s;background:%s;border:1px solid %s">%s</span></td>
  <td><button class="del-btn" gs-trigger="click" gs-get="/api/tasks/%d/delete" gs-target="#task-%d" gs-swap="outerHTML">✕</button></td>
</tr>`, t.ID, doneClass, t.ID, t.ID, icon, doneStyle, t.Title, t.Project, c, bg, c, string(t.Priority), t.ID, t.ID)
}

func taskRows(tasks []Task) string {
	if len(tasks) == 0 {
		return `<tr><td colspan="5" class="empty">No tasks — add one above!</td></tr>`
	}
	var b strings.Builder
	for _, t := range tasks { b.WriteString(taskRow(t)) }
	return b.String()
}

func statsFragment(s *Store) string {
	total, done := s.Stats()
	pending := total - done
	pct := 0
	if total > 0 { pct = done * 100 / total }
	return fmt.Sprintf(`<div id="stats">
  <div class="sc"><span class="sn">%d</span><span class="sl">Total</span></div>
  <div class="sc done-c"><span class="sn">%d</span><span class="sl">Done</span></div>
  <div class="sc pend-c"><span class="sn">%d</span><span class="sl">Pending</span></div>
  <div class="sc pct-c"><span class="sn">%d%%</span><span class="sl">Complete</span></div>
</div>
<div id="prog-wrap"><div id="prog-bar" style="width:%d%%"></div></div>`, total, done, pending, pct, pct)
}

// ─── CSS ─────────────────────────────────────────────────────────────────────

const css = `
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
body{font-family:'Inter',system-ui,sans-serif;background:#0f172a;color:#e2e8f0;min-height:100vh}
a{color:#38bdf8;text-decoration:none}a:hover{text-decoration:underline}

nav{background:#1e293b;border-bottom:1px solid #334155;padding:0 2rem;display:flex;align-items:center;gap:1.5rem;height:56px}
.brand{font-size:1.05rem;font-weight:800;color:#38bdf8}
.nl{color:#94a3b8;font-size:.875rem;padding:.25rem .65rem;border-radius:6px;transition:all .15s}
.nl:hover,.nl.active{color:#f1f5f9;background:#334155;text-decoration:none}
.ver{margin-left:auto;font-size:.7rem;color:#475569;background:#0f172a;padding:.15rem .55rem;border-radius:999px;border:1px solid #334155}

.container{max-width:960px;margin:0 auto;padding:2rem 1.5rem}
h1{font-size:1.6rem;font-weight:700;margin-bottom:.2rem}
.sub{color:#64748b;font-size:.875rem;margin-bottom:1.75rem}

#stats{display:flex;gap:.75rem;margin-bottom:.6rem;flex-wrap:wrap}
.sc{flex:1;min-width:90px;background:#1e293b;border:1px solid #334155;border-radius:10px;padding:.85rem;text-align:center}
.sn{display:block;font-size:1.6rem;font-weight:800;color:#38bdf8}
.sl{font-size:.7rem;color:#64748b;text-transform:uppercase;letter-spacing:.05em}
.done-c .sn{color:#22c55e}.pend-c .sn{color:#f59e0b}.pct-c .sn{color:#a78bfa}
#prog-wrap{height:5px;background:#1e293b;border-radius:999px;margin-bottom:1.75rem;overflow:hidden}
#prog-bar{height:100%;background:linear-gradient(90deg,#38bdf8,#a78bfa);border-radius:999px;transition:width .4s ease}

.add-form{background:#1e293b;border:1px solid #334155;border-radius:10px;padding:1.1rem;margin-bottom:1.25rem;display:flex;gap:.65rem;flex-wrap:wrap;align-items:flex-end}
.field{display:flex;flex-direction:column;gap:.25rem;flex:1;min-width:130px}
.field label{font-size:.7rem;color:#64748b;text-transform:uppercase;letter-spacing:.05em}
.field input,.field select{background:#0f172a;border:1px solid #334155;border-radius:7px;color:#e2e8f0;padding:.5rem .7rem;font-size:.875rem;outline:none;transition:border-color .15s}
.field input:focus,.field select:focus{border-color:#38bdf8}
.field select option{background:#1e293b}
.btn-add{background:#0ea5e9;color:#fff;border:none;border-radius:7px;padding:.5rem 1.1rem;font-size:.875rem;font-weight:600;cursor:pointer;height:36px;transition:background .15s;white-space:nowrap}
.btn-add:hover{background:#38bdf8}

.tbl-wrap{background:#1e293b;border:1px solid #334155;border-radius:10px;overflow:hidden}
table{width:100%;border-collapse:collapse}
thead tr{background:#0f172a}
th{text-align:left;padding:.65rem 1rem;font-size:.7rem;color:#475569;text-transform:uppercase;letter-spacing:.05em;font-weight:600}
.task-row td{padding:.7rem 1rem;border-top:1px solid #1e3a5f18;font-size:.875rem}
.task-row:hover{background:#0f172a44}
.task-row.done{opacity:.6}
.check-btn{background:none;border:2px solid #334155;border-radius:50%;width:26px;height:26px;cursor:pointer;color:#94a3b8;font-size:.8rem;display:inline-flex;align-items:center;justify-content:center;transition:all .15s}
.check-btn:hover{border-color:#22c55e;color:#22c55e}
.del-btn{background:none;border:none;color:#475569;cursor:pointer;padding:.2rem .45rem;border-radius:5px;transition:all .15s;font-size:.8rem}
.del-btn:hover{background:#fef2f2;color:#ef4444}
.proj-tag{font-size:.72rem;background:#0f172a;border:1px solid #334155;border-radius:5px;padding:.12rem .45rem;color:#94a3b8}
.pri-badge{font-size:.68rem;padding:.12rem .45rem;border-radius:999px;font-weight:600;text-transform:capitalize}
.empty{padding:2.5rem;text-align:center;color:#475569}

.proj-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:.85rem;margin-top:1.25rem}
.proj-card{background:#1e293b;border:1px solid #334155;border-radius:10px;padding:1.1rem;transition:border-color .2s;display:block;color:inherit;text-decoration:none}
.proj-card:hover{border-color:#38bdf8;text-decoration:none}
.proj-card h3{font-size:.95rem;margin-bottom:.35rem}
.proj-card .cnt{color:#64748b;font-size:.8rem}

.feat-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(240px,1fr));gap:.85rem;margin-top:1.25rem}
.feat-card{background:#1e293b;border:1px solid #334155;border-radius:10px;padding:1.1rem}
.feat-card .ico{font-size:1.35rem;margin-bottom:.4rem}
.feat-card h3{font-size:.9rem;margin-bottom:.3rem;color:#38bdf8}
.feat-card p{font-size:.8rem;color:#64748b;line-height:1.55}

.gs-loading{opacity:.5;transition:opacity .2s}
`

// ─── Page shell ───────────────────────────────────────────────────────────────

func shell(pageTitle, activeNav, body string) string {
	meta := gs.NewMetadata().
		SetTitle(pageTitle+" — GoScript Task Manager").
		SetDescription("A full-featured test webapp built with GoScript 2.0, the Go-native web framework.").
		SetThemeColor("#0ea5e9").
		Build()

	activeIf := func(page string) string {
		if activeNav == page { return " active" }
		return ""
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  %s
  <script src="/__goscript/runtime.js"></script>
  <style>%s</style>
</head>
<body>
<nav>
  <span class="brand">⚡ GoScript</span>
  <a href="/" class="nl%s">Tasks</a>
  <a href="/projects" class="nl%s">Projects</a>
  <a href="/about" class="nl%s">About</a>
  <span class="ver">v2.0.0</span>
</nav>
<div class="container">%s</div>
</body>
</html>`, meta.Render(), css,
		activeIf("tasks"), activeIf("projects"), activeIf("about"),
		body)
}

// ─── Page bodies ──────────────────────────────────────────────────────────────

func bodyHome(s *Store) string {
	return fmt.Sprintf(`
<h1>Task Board</h1>
<p class="sub">GoScript 2.0 reactive demo — gs-* attributes, middleware pipeline, dynamic routes</p>

<div id="stats-section">%s</div>

<form class="add-form"
  gs-trigger="submit"
  gs-post="/api/tasks"
  gs-target="#tbody"
  gs-swap="beforeend"
  onsubmit="this.reset();return false;">
  <div class="field">
    <label>Title</label>
    <input type="text" name="title" placeholder="What needs doing?" required>
  </div>
  <div class="field" style="max-width:150px">
    <label>Project</label>
    <input type="text" name="project" value="General">
  </div>
  <div class="field" style="max-width:120px">
    <label>Priority</label>
    <select name="priority">
      <option value="high">🔴 High</option>
      <option value="medium" selected>🟡 Medium</option>
      <option value="low">🟢 Low</option>
    </select>
  </div>
  <button type="submit" class="btn-add">+ Add</button>
</form>

<div class="tbl-wrap">
  <table>
    <thead><tr>
      <th style="width:40px"></th>
      <th>Task</th>
      <th>Project</th>
      <th>Priority</th>
      <th style="width:40px"></th>
    </tr></thead>
    <tbody id="tbody">%s</tbody>
  </table>
</div>`, statsFragment(s), taskRows(s.All()))
}

func bodyProjects(s *Store) string {
	var cards strings.Builder
	for _, p := range s.Projects() {
		tasks := s.ByProject(p)
		done := 0
		for _, t := range tasks { if t.Done { done++ } }
		cards.WriteString(fmt.Sprintf(
			`<a href="/projects/%s" class="proj-card"><h3>%s</h3><p class="cnt">%d tasks · %d done</p></a>`,
			p, p, len(tasks), done))
	}
	return fmt.Sprintf(`<h1>Projects</h1><p class="sub">Click a project to filter tasks</p>
<div class="proj-grid">%s</div>`, cards.String())
}

func bodyProjectDetail(s *Store, project string) string {
	rows := taskRows(s.ByProject(project))
	return fmt.Sprintf(`<h1>%s</h1>
<p class="sub"><a href="/projects">← All Projects</a></p>
<div class="tbl-wrap">
  <table>
    <thead><tr>
      <th style="width:40px"></th><th>Task</th><th>Project</th><th>Priority</th><th style="width:40px"></th>
    </tr></thead>
    <tbody id="tbody">%s</tbody>
  </table>
</div>`, project, rows)
}

func bodyAbout() string {
	feats := [][3]string{
		{"🗺️", "App Router", "Next.js-style routing with dynamic :param and *catch-all segments, layouts, and route groups — pure Go."},
		{"⚡", "Reactive gs-* Attributes", "HTMX-inspired zero-JS reactivity: gs-trigger, gs-get, gs-post, gs-target, gs-swap on any HTML element."},
		{"🧩", "Server Components", "Composable server-rendered HTML components with prop validation, lifecycle hooks, and children support."},
		{"🔒", "Middleware Pipeline", "Composable pipeline: Gzip, CORS, Security Headers, Rate Limiting, Sessions, Logging, Recovery, Request ID."},
		{"📦", "Single Binary", "Entire stack compiles to one ~8MB static binary. Zero runtime dependencies — no Node.js, no npm."},
		{"🔍", "Metadata / SEO", "Fluent builder for Open Graph, Twitter Cards, JSON-LD structured data, robots, and canonical tags."},
		{"📝", "Forms + Validation", "Declarative server-side forms with built-in validation rules, CSRF protection, and file upload support."},
		{"🚀", "SSR + SSG + ISR", "Streaming SSR with Suspense boundaries, static generation, and incremental static regeneration built-in."},
	}
	var cards strings.Builder
	for _, f := range feats {
		cards.WriteString(fmt.Sprintf(
			`<div class="feat-card"><div class="ico">%s</div><h3>%s</h3><p>%s</p></div>`,
			f[0], f[1], f[2]))
	}
	return fmt.Sprintf(`<h1>About GoScript 2.0</h1>
<p class="sub">A Go-native web framework challenging Next.js — built by <strong>davidjeba</strong></p>
<div class="feat-grid">%s</div>`, cards.String())
}

// ─── main ────────────────────────────────────────────────────────────────────

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" { port = p }

	store := NewStore()

	// Middleware pipeline
	pipeline := gs.NewPipeline().
		Use(gs.RequestIDMiddleware()).
		Use(gs.GzipMiddleware()).
		Use(gs.CORSMiddleware(gs.DefaultCORSConfig())).
		Use(gs.SecurityHeadersMiddleware()).
		Use(gs.LoggingMiddleware(log.Printf)).
		Use(gs.RecoveryMiddleware(log.Printf))

	// App Router — page routes
	router := gs.NewAppRouter("/")

	router.RegisterRoute("/", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, shell("Tasks", "tasks", bodyHome(store)))
	}, []string{"GET"})

	router.RegisterRoute("/projects", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, shell("Projects", "projects", bodyProjects(store)))
	}, []string{"GET"})

	router.RegisterRoute("/projects/:project", func(w http.ResponseWriter, r *http.Request, p map[string]string) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, shell(p["project"], "projects", bodyProjectDetail(store, p["project"])))
	}, []string{"GET"})

	router.RegisterRoute("/about", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, shell("About", "about", bodyAbout()))
	}, []string{"GET"})

	// Main mux
	mux := http.NewServeMux()

	// Runtime.js
	mux.Handle("/__goscript/runtime.js", gs.RuntimeHandler())

	// API: add task (POST)
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.NotFound(w, r); return }
		r.ParseForm()
		title := strings.TrimSpace(r.FormValue("title"))
		project := strings.TrimSpace(r.FormValue("project"))
		priority := Priority(r.FormValue("priority"))
		if title == "" { http.Error(w, "title required", 400); return }
		if project == "" { project = "General" }
		if priority != PriorityHigh && priority != PriorityMedium && priority != PriorityLow {
			priority = PriorityMedium
		}
		t := store.Add(title, project, priority)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, taskRow(t))
	})

	// API: toggle / delete (GET with action suffix)
	mux.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 4 { http.NotFound(w, r); return }
		id, err := strconv.Atoi(parts[2])
		if err != nil { http.Error(w, "bad id", 400); return }
		w.Header().Set("Content-Type", "text/html")
		switch parts[3] {
		case "toggle":
			if t, ok := store.Toggle(id); ok { fmt.Fprint(w, taskRow(t)) }
		case "delete":
			store.Delete(id)
			fmt.Fprint(w, "") // empty = outerHTML removes the row
		default:
			http.NotFound(w, r)
		}
	})

	// API: stats fragment
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, statsFragment(store))
	})

	// API: health (JSON)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		total, done := store.Stats()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "healthy", "framework": "GoScript 2.0",
			"version": gs.Version,
			"tasks":   map[string]int{"total": total, "done": done},
			"ts":      time.Now().Format(time.RFC3339),
		})
	})

	// All page routes through pipeline → AppRouter
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pipeline.Execute(w, r, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			router.ServeHTTP(w, r)
		}))
	}))

	fmt.Printf(`
  ╔══════════════════════════════════════╗
  ║  ⚡ GoScript Task Manager            ║
  ║  http://localhost:%-18s║
  ╠══════════════════════════════════════╣
  ║  Page routes (AppRouter + pipeline)  ║
  ║   GET  /                             ║
  ║   GET  /projects                     ║
  ║   GET  /projects/:project            ║
  ║   GET  /about                        ║
  ╠══════════════════════════════════════╣
  ║  API routes (fragment + JSON)        ║
  ║   POST /api/tasks                    ║
  ║   GET  /api/tasks/:id/toggle         ║
  ║   GET  /api/tasks/:id/delete         ║
  ║   GET  /api/stats                    ║
  ║   GET  /api/health                   ║
  ╚══════════════════════════════════════╝
`+" ", port)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
