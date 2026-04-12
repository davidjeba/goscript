"use client";

import { useState } from "react";
import { Copy, Check } from "lucide-react";
import { cn } from "@/lib/utils";

interface CodeBlockProps {
  code: string;
  language?: string;
  title?: string;
  className?: string;
}

export function CodeBlock({
  code,
  language = "go",
  title,
  className,
}: CodeBlockProps) {
  const [copied, setCopied] = useState(false);
  const [expanded, setExpanded] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  // Simple Go syntax highlighting
  const highlighted = highlightGo(code);

  const isLong = code.split("\n").length > 20;
  const displayCode = expanded ? code : isLong ? code.split("\n").slice(0, 20).join("\n") : code;

  return (
    <div
      className={cn(
        "relative rounded-lg border border-border overflow-hidden bg-[#1e1e1e]",
        className
      )}
    >
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-2 bg-[#252526] border-b border-[#333]">
        <div className="flex items-center gap-3">
          <div className="flex gap-1.5">
            <div className="w-3 h-3 rounded-full bg-[#ff5f57]" />
            <div className="w-3 h-3 rounded-full bg-[#ffbd2e]" />
            <div className="w-3 h-3 rounded-full bg-[#28c840]" />
          </div>
          {title && (
            <span className="text-xs text-[#858585] font-mono">{title}</span>
          )}
          {!title && (
            <span className="text-xs text-[#858585] font-mono">{language}</span>
          )}
        </div>
        <button
          onClick={handleCopy}
          className="flex items-center gap-1.5 px-2 py-1 rounded text-xs text-[#858585] hover:text-[#cccccc] hover:bg-[#333] transition-colors"
        >
          {copied ? (
            <>
              <Check className="w-3.5 h-3.5 text-emerald-400" />
              <span className="text-emerald-400">Copied</span>
            </>
          ) : (
            <>
              <Copy className="w-3.5 h-3.5" />
              <span>Copy</span>
            </>
          )}
        </button>
      </div>

      {/* Code */}
      <div
        className={cn(
          "overflow-x-auto overflow-y-auto font-mono text-[13px] leading-6 p-4 custom-scrollbar",
          !expanded && isLong && "max-h-96"
        )}
      >
        <pre className="text-[#d4d4d4]">
          <code dangerouslySetInnerHTML={{ __html: highlightGo(displayCode) }} />
        </pre>
      </div>

      {/* Expand button */}
      {isLong && (
        <button
          onClick={() => setExpanded(!expanded)}
          className="w-full py-2 text-xs text-center text-[#858585] hover:text-[#cccccc] bg-[#252526] border-t border-[#333] transition-colors"
        >
          {expanded ? "Show less" : `Show more (${code.split("\n").length} lines)`}
        </button>
      )}
    </div>
  );
}

function highlightGo(code: string): string {
  let result = escapeHtml(code);

  // Comments (single-line)
  result = result.replace(
    /(\/\/[^\n]*)/g,
    '<span class="text-[#6a9955]">$1</span>'
  );

  // Comments (block)
  result = result.replace(
    /(\/\*[\s\S]*?\*\/)/g,
    '<span class="text-[#6a9955]">$1</span>'
  );

  // Strings (double quotes)
  result = result.replace(
    /(&quot;((?:[^&]|&(?!quot;))*?)&quot;)/g,
    '<span class="text-[#ce9178]">$1</span>'
  );

  // Strings (backticks/raw strings)
  result = result.replace(
    /(`[^`]*`)/g,
    '<span class="text-[#ce9178]">$1</span>'
  );

  // Keywords
  const keywords = [
    "package", "import", "func", "type", "struct", "interface", "map",
    "chan", "go", "defer", "return", "if", "else", "for", "range",
    "switch", "case", "default", "break", "continue", "select",
    "var", "const", "nil", "true", "false", "make", "new", "append",
    "len", "cap", "copy", "delete", "close", "panic", "recover",
  ];

  keywords.forEach((kw) => {
    const regex = new RegExp(`\\b(${kw})\\b`, "g");
    result = result.replace(
      regex,
      '<span class="text-[#569cd6]">$1</span>'
    );
  });

  // Types
  const types = [
    "string", "int", "int8", "int16", "int32", "int64",
    "uint", "uint8", "uint16", "uint32", "uint64",
    "float32", "float64", "bool", "byte", "rune", "error",
    "http.ResponseWriter", "http.Request", "http.Handler",
    "http.Flusher", "context.Context", "io.Writer", "io.Reader",
    "sync.RWMutex", "sync.Mutex", "sync.WaitGroup",
    "time.Duration", "time.Time", "os.FileInfo",
  ];

  types.forEach((t) => {
    const escaped = t.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    const regex = new RegExp(`\\b(${escaped})\\b`, "g");
    result = result.replace(
      regex,
      '<span class="text-[#4ec9b0]">$1</span>'
    );
  });

  // Numbers
  result = result.replace(
    /\b(\d+\.?\d*)\b/g,
    '<span class="text-[#b5cea8]">$1</span>'
  );

  // Function calls
  result = result.replace(
    /\b([a-zA-Z_][a-zA-Z0-9_]*)\s*\(/g,
    '<span class="text-[#dcdcaa]">$1</span>('
  );

  return result;
}

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}
