"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

export default function RegisterPage() {
  const router = useRouter();
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
      first_name: String(form.get("first_name")),
      last_name: String(form.get("last_name")),
      phone: String(form.get("phone") || ""),
      role: String(form.get("role") || "member"),
    };

    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/auth/register`, {
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
        setErr(data?.error || "Registration failed");
        return;
      }

      // After registering, go to login
      router.push("/auth/login");
    } catch (error) {
      console.error(error);
      setErr("Something went wrong. Please try again.");
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
            <h2 className="text-3xl font-bold text-white tracking-tight">Join <span className="text-primary">FitTrack+</span> today</h2>
            <p className="mt-3 text-white/80 max-w-md">Personalized plans, easy booking, and clear progress — get started in minutes.</p>
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
            <h1 className="text-2xl font-semibold tracking-tight">Create your account</h1>
            <p className="text-sm text-muted-foreground">Start your journey with a free account</p>
          </div>
          {err && <p className="text-sm text-red-600">{err}</p>}

          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            <input className="w-full border border-border px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary" name="first_name" placeholder="First name" required />
            <input className="w-full border border-border px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary" name="last_name" placeholder="Last name" required />
          </div>
          <input className="w-full border border-border px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary" name="email" type="email" placeholder="Email" required />
          <input className="w-full border border-border px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary" name="password" type="password" placeholder="Password" required />
          <input className="w-full border border-border px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary" name="phone" placeholder="Phone (optional)" />
          <select className="w-full border border-border px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary" name="role" defaultValue="member">
            <option value="member">Member</option>
            <option value="trainer">Trainer</option>
            <option value="physio">Physio</option>
            <option value="admin">Admin</option>
          </select>

          <button disabled={loading} className="w-full h-11 rounded-full bg-primary text-primary-foreground hover:bg-primary/90 transition-colors">
            {loading ? "Creating..." : "Sign up"}
          </button>

          <p className="text-sm text-muted-foreground">
            Already have an account? <a href="/auth/login" className="text-primary hover:underline">Log in</a>
          </p>
        </form>
      </div>
    </div>
  );
}
