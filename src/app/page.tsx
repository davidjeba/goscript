import Link from "next/link";
import { ArrowRight, Code2, Layers3, ShieldCheck, Sparkles, Terminal, Workflow } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

const comparisonRows = [
  {
    label: "What it is",
    goscript: "A Go-native language and runtime for web apps.",
    javascript: "The browser's native language and its enormous ecosystem.",
  },
  {
    label: "Ownership",
    goscript: "Frontend and backend can stay inside one Go mental model.",
    javascript: "UI work often introduces a second language and toolchain.",
  },
  {
    label: "Predictability",
    goscript: "Compiled, explicit, and easier to reason about operationally.",
    javascript: "Dynamic, flexible, and powerful, but often more fragmented.",
  },
  {
    label: "Deployment",
    goscript: "A Go-shaped delivery story that fits Go teams naturally.",
    javascript: "A runtime-plus-package-graph story that many teams already know.",
  },
  {
    label: "Best fit",
    goscript: "Teams that already trust Go and want the web to feel native.",
    javascript: "Teams that need direct access to the browser ecosystem.",
  },
];

const reasons = [
  {
    icon: Workflow,
    title: "One language, less context switching",
    body:
      "GoScript lets teams keep product logic, interface logic, and backend logic in the same language, which lowers friction and keeps ownership clear.",
  },
  {
    icon: ShieldCheck,
    title: "A more Go-shaped operational story",
    body:
      "Go developers already value predictability, explicitness, and simple deployment. GoScript extends those values to the web layer instead of outsourcing them to JavaScript.",
  },
  {
    icon: Layers3,
    title: "Familiar structure, different purpose",
    body:
      "It can follow the app-structure modern web developers recognize, but the real point is language ownership, not copying another framework.",
  },
  {
    icon: Sparkles,
    title: "A native blessing for the Go community",
    body:
      "GoScript gives Go teams a credible path to full-stack web development without leaving the language they already use for systems, services, and tooling.",
  },
];

export default function Home() {
  return (
    <main className="relative min-h-screen overflow-hidden bg-background">
      <div className="absolute inset-x-0 top-0 -z-10 h-[420px] bg-[radial-gradient(circle_at_top,rgba(16,185,129,0.18),transparent_58%)]" />
      <div className="absolute inset-x-0 top-40 -z-10 h-[320px] bg-[radial-gradient(circle_at_top_right,rgba(56,189,248,0.14),transparent_60%)]" />
      <div className="absolute inset-x-0 bottom-0 -z-10 h-[240px] bg-[radial-gradient(circle_at_bottom,rgba(245,158,11,0.10),transparent_60%)]" />

      <section className="container mx-auto flex max-w-6xl flex-col gap-10 px-4 py-16 sm:px-6 lg:px-8 lg:py-24">
        <div className="max-w-4xl space-y-6">
          <Badge
            variant="outline"
            className="border-emerald-500/30 bg-emerald-500/5 px-3 py-1 text-emerald-300"
          >
            Language-first, not framework-first
          </Badge>

          <div className="space-y-4">
            <h1 className="max-w-3xl text-4xl font-semibold tracking-tight sm:text-5xl lg:text-7xl">
              GoScript is the Go-native alternative to JavaScript.
            </h1>
            <p className="max-w-3xl text-lg leading-8 text-muted-foreground sm:text-xl">
              It follows a familiar app structure so the learning curve feels approachable, but its real purpose is different:
              GoScript lets Go developers build modern web experiences without handing the language layer to JavaScript.
            </p>
          </div>

          <div className="flex flex-wrap gap-3">
            <Button asChild className="bg-emerald-600 text-white hover:bg-emerald-700">
              <Link href="#comparison">
                Compare with JavaScript
                <ArrowRight className="size-4" />
              </Link>
            </Button>
            <Button asChild variant="outline">
              <Link href="#why">
                Why Go should care
              </Link>
            </Button>
          </div>
        </div>

        <div className="grid gap-4 md:grid-cols-3">
          <Card className="border-border/60 bg-card/70">
            <CardContent className="p-5">
              <p className="text-sm text-muted-foreground">Language ownership</p>
              <p className="mt-2 text-3xl font-semibold">1</p>
              <p className="mt-2 text-sm text-muted-foreground">
                One language from backend to UI instead of a split Go and JS stack.
              </p>
            </CardContent>
          </Card>
          <Card className="border-border/60 bg-card/70">
            <CardContent className="p-5">
              <p className="text-sm text-muted-foreground">Context switches avoided</p>
              <p className="mt-2 text-3xl font-semibold">Fewer</p>
              <p className="mt-2 text-sm text-muted-foreground">
                Teams stay in the Go mental model while building the web surface.
              </p>
            </CardContent>
          </Card>
          <Card className="border-border/60 bg-card/70">
            <CardContent className="p-5">
              <p className="text-sm text-muted-foreground">Deployment feel</p>
              <p className="mt-2 text-3xl font-semibold">Go-shaped</p>
              <p className="mt-2 text-sm text-muted-foreground">
                The operational story stays close to what Go teams already trust.
              </p>
            </CardContent>
          </Card>
        </div>
      </section>

      <section className="container mx-auto max-w-6xl px-4 pb-4 sm:px-6 lg:px-8">
        <Card className="border-emerald-500/20 bg-emerald-500/5">
          <CardContent className="flex flex-col gap-4 p-6 lg:flex-row lg:items-center lg:justify-between">
            <div className="space-y-2">
              <p className="text-sm font-medium uppercase tracking-[0.2em] text-emerald-300/80">
                Important framing
              </p>
              <p className="text-base leading-7 text-muted-foreground">
                GoScript is not a Next.js replacement. Next.js is a structure people already know.
                GoScript is the language-side answer for Go teams that want the web to feel native.
              </p>
            </div>
            <div className="flex items-center gap-3 text-sm text-muted-foreground">
              <Terminal className="size-4 text-emerald-300" />
              GoScript follows the structure; JavaScript is what it replaces.
            </div>
          </CardContent>
        </Card>
      </section>

      <section id="comparison" className="container mx-auto max-w-6xl px-4 py-16 sm:px-6 lg:px-8">
        <div className="mb-8 space-y-3">
          <Badge variant="outline">GoScript vs JavaScript</Badge>
          <h2 className="text-3xl font-semibold tracking-tight sm:text-4xl">
            A comparison the Go community can read without marketing fog.
          </h2>
          <p className="max-w-3xl text-muted-foreground">
            The goal is not to attack JavaScript. The goal is to show why GoScript deserves to exist as a native
            option for Go developers who want the web without a language split.
          </p>
        </div>

        <Card className="border-border/60 bg-card/70">
          <CardContent className="overflow-x-auto p-0">
            <div className="min-w-[760px]">
              <div className="grid grid-cols-[1.1fr_1.2fr_1.2fr] border-b border-border/60 bg-muted/30 px-6 py-4 text-sm font-medium text-muted-foreground">
                <div>Area</div>
                <div>GoScript</div>
                <div>JavaScript</div>
              </div>
              {comparisonRows.map((row) => (
                <div
                  key={row.label}
                  className="grid grid-cols-[1.1fr_1.2fr_1.2fr] gap-4 border-b border-border/60 px-6 py-5 last:border-b-0"
                >
                  <div className="font-medium">{row.label}</div>
                  <div className="text-sm leading-7 text-muted-foreground">{row.goscript}</div>
                  <div className="text-sm leading-7 text-muted-foreground">{row.javascript}</div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </section>

      <section id="why" className="container mx-auto max-w-6xl px-4 pb-16 sm:px-6 lg:px-8">
        <div className="mb-8 space-y-3">
          <Badge variant="outline">Why it matters</Badge>
          <h2 className="text-3xl font-semibold tracking-tight sm:text-4xl">
            Why the Go community should consider GoScript a native blessing.
          </h2>
          <p className="max-w-3xl text-muted-foreground">
            GoScript is compelling because it lets Go teams keep the benefits they already value while finally
            treating the web as a first-class Go workload.
          </p>
        </div>

        <div className="grid gap-6 lg:grid-cols-2">
          {reasons.map((reason) => {
            const Icon = reason.icon;
            return (
              <Card key={reason.title} className="border-border/60 bg-card/70">
                <CardHeader className="space-y-4">
                  <div className="flex size-12 items-center justify-center rounded-2xl bg-emerald-500/10 text-emerald-300">
                    <Icon className="size-6" />
                  </div>
                  <div className="space-y-2">
                    <CardTitle>{reason.title}</CardTitle>
                    <CardDescription className="text-base leading-7">
                      {reason.body}
                    </CardDescription>
                  </div>
                </CardHeader>
              </Card>
            );
          })}
        </div>
      </section>

      <Separator className="container mx-auto max-w-6xl" />

      <section className="container mx-auto max-w-6xl px-4 py-16 sm:px-6 lg:px-8">
        <Card className="border-border/60 bg-card/70">
          <CardContent className="space-y-4 p-6 sm:p-8">
            <div className="flex flex-wrap items-center gap-3">
              <Code2 className="size-5 text-emerald-300" />
              <p className="text-sm font-medium uppercase tracking-[0.2em] text-muted-foreground">
                The short version
              </p>
            </div>
            <p className="max-w-4xl text-lg leading-8 text-muted-foreground">
              JavaScript will keep being the language of the browser. GoScript is valuable for a different reason:
              it gives the Go community a credible, native way to build modern web apps without abandoning Go.
              That is why it can feel like a blessing, not just another framework idea.
            </p>
          </CardContent>
        </Card>
      </section>
    </main>
  );
}
