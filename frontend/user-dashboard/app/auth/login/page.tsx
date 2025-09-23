"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { useState } from "react";

export default function LoginPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const next = searchParams.get("next") || "/dashboard";
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setErr(null);
    setLoading(true);

    const form = new FormData(e.currentTarget);
    const payload = {
      email: String(form.get("email")),
      password: String(form.get("password")),
    };

    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      const raw = await res.text();
      let data: any = {};
      try {
        data = raw ? JSON.parse(raw) : {};
      } catch {}

      setLoading(false);

      if (!res.ok) {
        setErr(data?.error || "Login failed");
        return;
      }

      // Save token in cookie (for example, using document.cookie)
      document.cookie = `ft_token=${data.token}; path=/`;

      // 👇 NEW LOGIC
      if (!data.user?.hasProfile) {
        router.push("/profile"); // redirect if profile not complete
      } else {
        router.push(next || "/dashboard"); // otherwise to dashboard
      }

    } catch (error) {
      console.error(error);
      setErr("Could not connect to server.");
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen grid md:grid-cols-2 bg-background">
      {/* Left: Visual / Animation */}
      <div className="relative hidden md:block">
        <div className="absolute inset-0 bg-[linear-gradient(to_bottom_right,rgba(2,44,34,0.95),rgba(15,23,42,0.95))]" />
        <div className="absolute inset-0 p-10 flex flex-col justify-end">
          <div className="mb-auto">
            <h2 className="text-3xl font-bold text-white tracking-tight">Welcome back to <span className="text-primary">FitTrack+</span></h2>
            <p className="mt-3 text-white/80 max-w-md">Sign in to book sessions faster, follow your plan, and track progress elegantly.</p>
          </div>
          {/* Animated blobs */}
          <div className="pointer-events-none absolute -top-10 -left-10 h-48 w-48 rounded-full bg-primary/30 blur-3xl animate-pulse" />
          <div className="pointer-events-none absolute bottom-20 right-10 h-40 w-40 rounded-full bg-emerald-400/20 blur-3xl animate-pulse" />
          <div className="pointer-events-none absolute bottom-6 left-6 h-24 w-24 rounded-full bg-white/10 blur-2xl" />
        </div>
      </div>

      {/* Right: Form */}
      <div className="flex items-center justify-center p-6 md:p-10">
        <form onSubmit={onSubmit} className="w-full max-w-md space-y-5 p-6 md:p-8 rounded-2xl border border-primary/10 bg-background/80 backdrop-blur-sm shadow-sm">
          <div>
            <h1 className="text-2xl font-semibold tracking-tight">Sign in</h1>
            <p className="text-sm text-muted-foreground">Use your email and password to continue</p>
          </div>
          {err && <p className="text-sm text-red-600">{err}</p>}

          <div className="space-y-3">
            <input className="w-full border border-border px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary" name="email" type="email" placeholder="Email" required />
            <input className="w-full border border-border px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary" name="password" type="password" placeholder="Password" required />
          </div>

          <button disabled={loading} className="w-full h-11 rounded-full bg-primary text-primary-foreground hover:bg-primary/90 transition-colors">
            {loading ? "Signing in..." : "Sign in"}
          </button>

          <p className="text-sm text-muted-foreground">
            New here? <a href="/auth/register" className="text-primary hover:underline">Create an account</a>
          </p>
        </form>
      </div>
    </div>
  );
}
