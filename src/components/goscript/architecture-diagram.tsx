"use client";

import { useRef } from "react";
import { motion, useInView } from "framer-motion";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface ArchNode {
  label: string;
  sublabel?: string;
  color: string;
}

interface ArchLayer {
  title: string;
  nodes: ArchNode[];
  description: string;
}

const layers: ArchLayer[] = [
  {
    title: "Request Layer",
    description: "HTTP server, WebSocket HMR, static file serving",
    nodes: [
      { label: "HTTP Server", color: "emerald" },
      { label: "HMR WebSocket", color: "amber" },
      { label: "Static Files", color: "cyan" },
    ],
  },
  {
    title: "Pipeline Layer",
    description: "Middleware chain with built-in handlers",
    nodes: [
      { label: "CORS", color: "purple" },
      { label: "Gzip", color: "rose" },
      { label: "Security", color: "emerald" },
      { label: "Rate Limit", color: "amber" },
      { label: "Session", color: "cyan" },
    ],
  },
  {
    title: "Routing Layer",
    description: "File-system auto-discovery with trie matching",
    nodes: [
      { label: "App Router", color: "emerald" },
      { label: "API Router", color: "amber" },
      { label: "Route Groups", color: "purple" },
    ],
  },
  {
    title: "Rendering Layer",
    description: "SSR, SSG, ISR, and Streaming with Suspense boundaries",
    nodes: [
      { label: "Streaming SSR", color: "emerald" },
      { label: "SSG / ISR", color: "rose" },
      { label: "Suspense", color: "amber" },
      { label: "Error Bounds", color: "cyan" },
    ],
  },
  {
    title: "Component Layer",
    description: "Server and client components with zero-bundle optimization",
    nodes: [
      { label: "Server Comp", color: "emerald" },
      { label: "Client Comp", color: "amber" },
      { label: "Layouts", color: "purple" },
      { label: "Metadata", color: "rose" },
    ],
  },
];

const colorClasses: Record<string, string> = {
  emerald: "bg-emerald-500/15 text-emerald-400 border-emerald-500/30",
  amber: "bg-amber-500/15 text-amber-400 border-amber-500/30",
  purple: "bg-purple-500/15 text-purple-400 border-purple-500/30",
  rose: "bg-rose-500/15 text-rose-400 border-rose-500/30",
  cyan: "bg-cyan-500/15 text-cyan-400 border-cyan-500/30",
};

export function ArchitectureDiagram() {
  const ref = useRef<HTMLDivElement>(null);
  const isInView = useInView(ref, { once: true, margin: "-80px" });

  return (
    <motion.div
      ref={ref}
      initial={{ opacity: 0, y: 40 }}
      animate={isInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 40 }}
      transition={{ duration: 0.6 }}
    >
      <Card className="border-border/50 bg-card/80 backdrop-blur-sm overflow-hidden">
        <CardHeader>
          <CardTitle className="text-xl">Architecture Overview</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-3">
            {layers.map((layer, layerIdx) => (
              <motion.div
                key={layer.title}
                initial={{ opacity: 0, x: -20 }}
                animate={isInView ? { opacity: 1, x: 0 } : {}}
                transition={{ duration: 0.4, delay: layerIdx * 0.1 }}
                className="relative"
              >
                {/* Connection arrow */}
                {layerIdx > 0 && (
                  <div className="flex justify-center py-1">
                    <div className="w-px h-4 bg-gradient-to-b from-border to-border/30" />
                  </div>
                )}

                <div className="rounded-lg border border-border/50 bg-muted/20 p-4">
                  <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 mb-3">
                    <div>
                      <h4 className="text-sm font-semibold text-foreground">
                        {layer.title}
                      </h4>
                      <p className="text-xs text-muted-foreground mt-0.5">
                        {layer.description}
                      </p>
                    </div>
                    <Badge
                      variant="outline"
                      className="text-[10px] px-2 py-0 w-fit"
                    >
                      Layer {layerIdx + 1}
                    </Badge>
                  </div>
                  <div className="flex flex-wrap gap-2">
                    {layer.nodes.map((node) => (
                      <div
                        key={node.label}
                        className={`inline-flex items-center rounded-md border px-3 py-1.5 text-xs font-medium ${colorClasses[node.color]}`}
                      >
                        {node.label}
                      </div>
                    ))}
                  </div>
                </div>
              </motion.div>
            ))}
          </div>

          {/* Data flow note */}
          <div className="rounded-lg border border-emerald-500/20 bg-emerald-500/5 p-3 mt-4">
            <p className="text-xs text-emerald-400 text-center">
              <span className="font-medium">Data Flow:</span> Request → Pipeline → Routing → Rendering → Component Tree → HTML Response
            </p>
          </div>
        </CardContent>
      </Card>
    </motion.div>
  );
}
