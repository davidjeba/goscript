"use client";

import { useRef } from "react";
import { motion, useInView } from "framer-motion";
import { Trophy, Equal } from "lucide-react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  comparisonData,
  goscriptWins,
  nextjsWins,
  ties,
} from "@/lib/goscript-improvements";

export function ComparisonTable() {
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
        <CardHeader className="pb-4">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <CardTitle className="text-xl">
              Feature Comparison
            </CardTitle>
            <div className="flex items-center gap-3 flex-wrap">
              <Badge className="bg-emerald-500/15 text-emerald-400 border-emerald-500/30">
                <Trophy className="w-3 h-3 mr-1" />
                GoScript 2.0: {goscriptWins} wins
              </Badge>
              <Badge className="bg-amber-500/15 text-amber-400 border-amber-500/30">
                <Trophy className="w-3 h-3 mr-1" />
                Next.js: {nextjsWins} wins
              </Badge>
              <Badge variant="outline" className="text-muted-foreground">
                <Equal className="w-3 h-3 mr-1" />
                Tied: {ties}
              </Badge>
            </div>
          </div>
        </CardHeader>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="hover:bg-transparent border-border/50">
                  <TableHead className="w-[200px] text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                    Feature
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-emerald-400 uppercase tracking-wider">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2 w-2 rounded-full bg-emerald-500" />
                      GoScript 2.0
                    </div>
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-amber-400 uppercase tracking-wider">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2 w-2 rounded-full bg-amber-500" />
                      Next.js 16
                    </div>
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2 w-2 rounded-full bg-muted-foreground/50" />
                      Original goscript
                    </div>
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {comparisonData.map((row) => (
                  <TableRow
                    key={row.feature}
                    className="border-border/30 hover:bg-muted/30 transition-colors"
                  >
                    <TableCell className="font-medium text-sm py-3">
                      {row.feature}
                    </TableCell>
                    <TableCell
                      className={`text-sm py-3 ${
                        row.winner === "goscript2"
                          ? "text-emerald-400 font-medium"
                          : "text-muted-foreground"
                      }`}
                    >
                      <div className="flex items-center gap-1.5">
                        {row.winner === "goscript2" && (
                          <Trophy className="w-3.5 h-3.5 text-emerald-500 shrink-0" />
                        )}
                        {row.goscript2}
                      </div>
                    </TableCell>
                    <TableCell
                      className={`text-sm py-3 ${
                        row.winner === "nextjs"
                          ? "text-amber-400 font-medium"
                          : "text-muted-foreground"
                      }`}
                    >
                      <div className="flex items-center gap-1.5">
                        {row.winner === "nextjs" && (
                          <Trophy className="w-3.5 h-3.5 text-amber-500 shrink-0" />
                        )}
                        {row.nextjs}
                      </div>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground/60 py-3">
                      {row.original}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </motion.div>
  );
}
