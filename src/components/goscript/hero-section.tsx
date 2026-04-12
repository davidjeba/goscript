"use client";

import { motion } from "framer-motion";
import {
  ArrowDown,
  Github,
  Sparkles,
  Gauge,
  Shield,
  Code2,
} from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

export function HeroSection() {
  return (
    <section className="relative min-h-screen flex items-center justify-center overflow-hidden">
      {/* Animated gradient background */}
      <div className="absolute inset-0 -z-10">
        <div className="absolute inset-0 bg-gradient-to-br from-background via-background to-background" />
        <div className="hero-gradient absolute inset-0 opacity-30" />
        {/* Grid pattern */}
        <div
          className="absolute inset-0 opacity-[0.03]"
          style={{
            backgroundImage:
              "linear-gradient(rgba(255,255,255,0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.1) 1px, transparent 1px)",
            backgroundSize: "60px 60px",
          }}
        />
        {/* Floating orbs */}
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-emerald-500/10 rounded-full blur-3xl animate-pulse" />
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-amber-500/10 rounded-full blur-3xl animate-pulse [animation-delay:1s]" />
      </div>

      <div className="container mx-auto px-4 py-24 lg:py-32">
        <div className="flex flex-col items-center text-center max-w-5xl mx-auto space-y-8">
          {/* Version badge */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
          >
            <Badge
              variant="outline"
              className="px-4 py-1.5 text-sm border-emerald-500/30 bg-emerald-500/5 text-emerald-400 hover:bg-emerald-500/10 transition-colors gap-2"
            >
              <Sparkles className="w-3.5 h-3.5" />
              Version 2.0 — Complete Rewrite
            </Badge>
          </motion.div>

          {/* Title */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.1 }}
            className="space-y-4"
          >
            <h1 className="text-5xl sm:text-6xl lg:text-8xl font-bold tracking-tight">
              <span className="text-foreground">Go</span>
              <span className="text-emerald-500">Script</span>{" "}
              <span className="text-muted-foreground/60">2.0</span>
            </h1>
            <p className="text-xl sm:text-2xl lg:text-3xl text-muted-foreground font-light max-w-3xl mx-auto leading-relaxed">
              A{" "}
              <span className="text-foreground font-medium">
                production-ready
              </span>{" "}
              Go web framework that brings Next.js-level DX with{" "}
              <span className="text-emerald-400 font-medium">
                compiled performance
              </span>
            </p>
          </motion.div>

          {/* Feature highlights */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            className="flex flex-wrap justify-center gap-4 text-sm"
          >
            {[
              { icon: Gauge, label: "Compiled binary, instant cold starts" },
              { icon: Code2, label: "File-system routing & App Router" },
              { icon: Shield, label: "Built-in middleware & security" },
              { icon: Sparkles, label: "Streaming SSR & HMR" },
            ].map(({ icon: Icon, label }) => (
              <div
                key={label}
                className="flex items-center gap-2 text-muted-foreground"
              >
                <Icon className="w-4 h-4 text-emerald-500" />
                <span>{label}</span>
              </div>
            ))}
          </motion.div>

          {/* CTA buttons */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.3 }}
            className="flex flex-col sm:flex-row gap-3"
          >
            <Button
              size="lg"
              className="gap-2 bg-emerald-600 hover:bg-emerald-700 text-white px-6"
              onClick={() => {
                document
                  .getElementById("improvements")
                  ?.scrollIntoView({ behavior: "smooth" });
              }}
            >
              Explore Improvements
              <ArrowDown className="w-4 h-4" />
            </Button>
            <Button
              size="lg"
              variant="outline"
              className="gap-2 px-6"
              onClick={() => {
                document
                  .getElementById("comparison")
                  ?.scrollIntoView({ behavior: "smooth" });
              }}
            >
              <Github className="w-4 h-4" />
              View Comparison
            </Button>
          </motion.div>

          {/* Quick stats */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.4 }}
            className="grid grid-cols-2 sm:grid-cols-4 gap-6 pt-8"
          >
            {[
              { value: "10", label: "New Features" },
              { value: "0ms", label: "Cold Start" },
              { value: "~5KB", label: "JS Bundle" },
              { value: "8+", label: "Middleware" },
            ].map((stat) => (
              <div key={stat.label} className="text-center">
                <div className="text-2xl sm:text-3xl font-bold text-emerald-500">
                  {stat.value}
                </div>
                <div className="text-xs text-muted-foreground mt-1">
                  {stat.label}
                </div>
              </div>
            ))}
          </motion.div>
        </div>
      </div>

      {/* Scroll indicator */}
      <motion.div
        className="absolute bottom-8 left-1/2 -translate-x-1/2"
        animate={{ y: [0, 8, 0] }}
        transition={{ duration: 2, repeat: Infinity }}
      >
        <ArrowDown className="w-5 h-5 text-muted-foreground/40" />
      </motion.div>
    </section>
  );
}
