import { useCallback, useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

type Status = "loading" | "ok" | "error";

const statusText: Record<Status, string> = {
  loading: "Checking API…",
  ok: "API is healthy",
  error: "API is unreachable",
};

const statusDot: Record<Status, string> = {
  loading: "bg-muted-foreground animate-pulse",
  ok: "bg-emerald-500",
  error: "bg-destructive",
};

export default function App() {
  const [status, setStatus] = useState<Status>("loading");

  const check = useCallback(async () => {
    setStatus("loading");

    try {
      const response = await fetch("/api/health");
      setStatus(response.ok ? "ok" : "error");
    } catch {
      setStatus("error");
    }
  }, []);

  useEffect(() => {
    void check();
  }, [check]);

  return (
    <main className="flex min-h-screen items-center justify-center p-6">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>goilerplate</CardTitle>
          <CardDescription>Chi API with a React frontend, built by Vite+.</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          <div className="flex items-center gap-2 text-sm" aria-live="polite">
            <span className={`size-2 rounded-full ${statusDot[status]}`} aria-hidden="true" />
            <span>{statusText[status]}</span>
          </div>
          <Button variant="outline" onClick={check} disabled={status === "loading"}>
            Check again
          </Button>
          <a
            className="text-muted-foreground text-xs hover:underline"
            href="/api/health"
            target="_blank"
            rel="noreferrer"
          >
            /api/health
          </a>
        </CardContent>
      </Card>
    </main>
  );
}
