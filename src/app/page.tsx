"use client";

import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import { Menu, X, Go } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { ThemeToggle } from "@/components/theme-toggle";
import { HeroSection } from "@/components/goscript/hero-section";
import { ImprovementCard } from "@/components/goscript/improvement-card";
import { ComparisonTable } from "@/components/goscript/comparison-table";
import { ArchitectureDiagram } from "@/components/goscript/architecture-diagram";
import { goscriptImprovements, categoryLabels, type ImprovementCategory } from "@/lib/goscript-improvements";

const navSections = [
  { id: "hero", label: "Home" },
  { id: "improvements", label: "Improvements" },
  { id: "architecture", label: "Architecture" },
  { id: "comparison", label: "Comparison" },
];

const categories: (ImprovementCategory | "all")[] = [
  "all",
  "routing",
  "rendering",
  "api",
  "performance",
  "dx",
];

export default function Home() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [activeFilter, setActiveFilter] = useState<ImprovementCategory | "all">("all");
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 20);
    };
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  const filteredImprovements =
    activeFilter === "all"
      ? goscriptImprovements
      : goscriptImprovements.filter((i) => i.category === activeFilter);

  return (
    <div className="min-h-screen flex flex-col">
      {/* Sticky Navigation */}
      <motion.header
        className={`sticky top-0 z-50 w-full transition-all duration-300 ${
          scrolled
            ? "bg-background/80 backdrop-blur-xl border-b border-border/50 shadow-sm"
            : "bg-transparent"
        }`}
        initial={{ y: -100 }}
        animate={{ y: 0 }}
        transition={{ duration: 0.3 }}
      >
        <nav className="container mx-auto flex h-14 items-center justify-between px-4">
          <div className="flex items-center gap-6">
            {/* Logo */}
            <a
              href="#hero"
              className="flex items-center gap-2 font-bold text-lg hover:opacity-80 transition-opacity"
            >
              <Go className="h-5 w-5 text-emerald-500" />
              <span>
                Go<span className="text-emerald-500">Script</span>{" "}
                <span className="text-muted-foreground/60 text-sm font-normal">2.0</span>
              </span>
            </a>

            {/* Desktop nav links */}
            <div className="hidden md:flex items-center gap-4">
              {navSections.map((section) => (
                <a
                  key={section.id}
                  href={`#${section.id}`}
                  className="nav-link text-sm text-muted-foreground hover:text-foreground transition-colors"
                >
                  {section.label}
                </a>
              ))}
            </div>
          </div>

          <div className="flex items-center gap-2">
            <ThemeToggle />
            {/* Mobile menu button */}
            <Button
              variant="ghost"
              size="icon"
              className="md:hidden"
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            >
              {mobileMenuOpen ? (
                <X className="h-4 w-4" />
              ) : (
                <Menu className="h-4 w-4" />
              )}
            </Button>
          </div>
        </nav>

        {/* Mobile menu */}
        {mobileMenuOpen && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: "auto" }}
            exit={{ opacity: 0, height: 0 }}
            className="md:hidden border-t border-border/50 bg-background/95 backdrop-blur-xl"
          >
            <div className="container mx-auto px-4 py-4 flex flex-col gap-3">
              {navSections.map((section) => (
                <a
                  key={section.id}
                  href={`#${section.id}`}
                  className="text-sm text-muted-foreground hover:text-foreground transition-colors py-2"
                  onClick={() => setMobileMenuOpen(false)}
                >
                  {section.label}
                </a>
              ))}
            </div>
          </motion.div>
        )}
      </motion.header>

      {/* Main Content */}
      <main className="flex-1">
        {/* Hero Section */}
        <section id="hero">
          <HeroSection />
        </section>

        {/* Improvements Section */}
        <section id="improvements" className="py-16 lg:py-24">
          <div className="container mx-auto px-4">
            <motion.div
              className="text-center mb-12 space-y-4"
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5 }}
            >
              <h2 className="text-3xl sm:text-4xl font-bold">
                10 Key{" "}
                <span className="text-emerald-500">Improvements</span>
              </h2>
              <p className="text-muted-foreground max-w-2xl mx-auto text-lg">
                Every major feature added to GoScript 2.0 to match and surpass
                Next.js capabilities — all in idiomatic Go.
              </p>
            </motion.div>

            {/* Category Filter */}
            <motion.div
              className="flex flex-wrap justify-center gap-2 mb-10"
              initial={{ opacity: 0, y: 10 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.4, delay: 0.1 }}
            >
              {categories.map((cat) => {
                const isActive = activeFilter === cat;
                const label = cat === "all" ? "All" : categoryLabels[cat];
                const count =
                  cat === "all"
                    ? goscriptImprovements.length
                    : goscriptImprovements.filter((i) => i.category === cat).length;

                return (
                  <button
                    key={cat}
                    onClick={() => setActiveFilter(cat)}
                    className={`inline-flex items-center gap-1.5 rounded-full px-4 py-1.5 text-sm font-medium border transition-all duration-200 ${
                      isActive
                        ? "bg-emerald-500/15 text-emerald-400 border-emerald-500/30"
                        : "bg-muted/50 text-muted-foreground border-border/50 hover:bg-muted hover:text-foreground"
                    }`}
                  >
                    {label}
                    <span
                      className={`text-xs px-1.5 py-0.5 rounded-full ${
                        isActive
                          ? "bg-emerald-500/20 text-emerald-400"
                          : "bg-muted text-muted-foreground"
                      }`}
                    >
                      {count}
                    </span>
                  </button>
                );
              })}
            </motion.div>

            {/* Improvement Cards Grid */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {filteredImprovements.map((improvement, index) => (
                <ImprovementCard
                  key={improvement.id}
                  improvement={improvement}
                  index={index}
                />
              ))}
            </div>
          </div>
        </section>

        <Separator className="container mx-auto max-w-5xl" />

        {/* Architecture Section */}
        <section id="architecture" className="py-16 lg:py-24">
          <div className="container mx-auto px-4 max-w-4xl">
            <motion.div
              className="text-center mb-12 space-y-4"
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5 }}
            >
              <h2 className="text-3xl sm:text-4xl font-bold">
                System{" "}
                <span className="text-emerald-500">Architecture</span>
              </h2>
              <p className="text-muted-foreground max-w-2xl mx-auto text-lg">
                A layered architecture designed for performance, extensibility,
                and developer happiness.
              </p>
            </motion.div>

            <ArchitectureDiagram />
          </div>
        </section>

        <Separator className="container mx-auto max-w-5xl" />

        {/* Comparison Section */}
        <section id="comparison" className="py-16 lg:py-24">
          <div className="container mx-auto px-4 max-w-5xl">
            <motion.div
              className="text-center mb-12 space-y-4"
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5 }}
            >
              <h2 className="text-3xl sm:text-4xl font-bold">
                Feature{" "}
                <span className="text-amber-500">Comparison</span>
              </h2>
              <p className="text-muted-foreground max-w-2xl mx-auto text-lg">
                How GoScript 2.0 stacks up against Next.js 16 and the original
                goscript framework across 20 key dimensions.
              </p>
            </motion.div>

            <ComparisonTable />
          </div>
        </section>
      </main>

      {/* Sticky Footer */}
      <footer className="border-t border-border/50 bg-card/50 backdrop-blur-sm mt-auto">
        <div className="container mx-auto px-4 py-8">
          <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
            <div className="flex items-center gap-2">
              <Go className="h-4 w-4 text-emerald-500" />
              <span className="text-sm font-medium">
                Go<span className="text-emerald-500">Script</span>{" "}
                <span className="text-muted-foreground">2.0</span>
              </span>
            </div>

            <p className="text-xs text-muted-foreground text-center">
              A showcase of framework improvements — designed to challenge and
              surpass Next.js with Go&apos;s performance.
            </p>

            <div className="flex items-center gap-4 text-xs text-muted-foreground">
              <span>Built with Next.js 16</span>
              <Separator orientation="vertical" className="h-3" />
              <span>Go + TypeScript</span>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
