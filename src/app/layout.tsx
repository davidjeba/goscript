import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { ThemeProvider } from "@/components/theme-provider";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "GoScript 2.0 — Full-Stack Go Web Framework",
  description:
    "GoScript 2.0 is a production-ready Go web framework bringing Next.js-level developer experience with compiled Go performance. Streaming SSR, file-system routing, server/client components, and more.",
  keywords: [
    "GoScript",
    "Go",
    "web framework",
    "SSR",
    "SSG",
    "streaming",
    "server components",
    "full-stack",
  ],
  authors: [{ name: "GoScript Team" }],
  openGraph: {
    title: "GoScript 2.0",
    description: "Full-Stack Go Web Framework with Next.js-level DX",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning className="dark">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased bg-background text-foreground`}
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="dark"
          enableSystem
          disableTransitionOnChange
        >
          {children}
        </ThemeProvider>
      </body>
    </html>
  );
}
