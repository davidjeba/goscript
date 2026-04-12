"use client";

import { useRef } from "react";
import { motion, useInView } from "framer-motion";
import {
  Route,
  Stream,
  Server,
  Webhook,
  Layers,
  HardDrive,
  ShieldAlert,
  Search,
  Zap,
  Terminal,
} from "lucide-react";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { CodeBlock } from "./code-block";
import {
  type GoImprovement,
  categoryLabels,
  categoryColors,
} from "@/lib/goscript-improvements";

const iconMap: Record<string, React.ComponentType<{ className?: string }>> = {
  Route,
  Stream,
  Server,
  Webhook,
  Layers,
  HardDrive,
  ShieldAlert,
  Search,
  Zap,
  Terminal,
};

interface ImprovementCardProps {
  improvement: GoImprovement;
  index: number;
}

export function ImprovementCard({ improvement, index }: ImprovementCardProps) {
  const ref = useRef<HTMLDivElement>(null);
  const isInView = useInView(ref, { once: true, margin: "-80px" });

  const IconComponent = iconMap[improvement.icon] || Terminal;

  return (
    <motion.div
      ref={ref}
      initial={{ opacity: 0, y: 40 }}
      animate={isInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 40 }}
      transition={{ duration: 0.5, delay: 0.1 * (index % 3) }}
    >
      <Card className="border-border/50 bg-card/80 backdrop-blur-sm overflow-hidden hover:border-border transition-colors">
        <CardHeader className="pb-4">
          <div className="flex items-start justify-between gap-4">
            <div className="flex items-start gap-3">
              <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                <IconComponent className="h-5 w-5" />
              </div>
              <div className="space-y-1">
                <div className="flex items-center gap-2 flex-wrap">
                  <CardTitle className="text-lg">{improvement.title}</CardTitle>
                  <Badge
                    variant="outline"
                    className={`text-[10px] px-2 py-0 ${categoryColors[improvement.category]}`}
                  >
                    {categoryLabels[improvement.category]}
                  </Badge>
                </div>
                <CardDescription className="text-sm">
                  {improvement.subtitle}
                </CardDescription>
              </div>
            </div>
            <span className="text-xs text-muted-foreground font-mono shrink-0">
              #{String(index + 1).padStart(2, "0")}
            </span>
          </div>
        </CardHeader>

        <CardContent className="space-y-4">
          <Tabs defaultValue="problem" className="w-full">
            <TabsList className="w-full grid grid-cols-3 h-9">
              <TabsTrigger value="problem" className="text-xs">
                Problem
              </TabsTrigger>
              <TabsTrigger value="solution" className="text-xs">
                Solution
              </TabsTrigger>
              <TabsTrigger value="code" className="text-xs">
                Code
              </TabsTrigger>
            </TabsList>

            <TabsContent value="problem" className="mt-3">
              <div className="rounded-lg border border-destructive/20 bg-destructive/5 p-4">
                <div className="flex items-center gap-2 mb-2">
                  <div className="h-2 w-2 rounded-full bg-destructive" />
                  <span className="text-xs font-medium text-destructive uppercase tracking-wide">
                    Before
                  </span>
                </div>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  {improvement.problem}
                </p>
              </div>
            </TabsContent>

            <TabsContent value="solution" className="mt-3">
              <div className="rounded-lg border border-emerald-500/20 bg-emerald-500/5 p-4">
                <div className="flex items-center gap-2 mb-2">
                  <div className="h-2 w-2 rounded-full bg-emerald-500" />
                  <span className="text-xs font-medium text-emerald-500 uppercase tracking-wide">
                    After
                  </span>
                </div>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  {improvement.solution}
                </p>
              </div>
            </TabsContent>

            <TabsContent value="code" className="mt-3">
              <CodeBlock
                code={improvement.code}
                title={improvement.id.replace(/-/g, "_") + ".go"}
              />
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </motion.div>
  );
}
