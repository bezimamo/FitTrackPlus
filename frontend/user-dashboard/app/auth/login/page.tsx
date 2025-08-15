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

      router.push(next);
    } catch (error) {
      console.error(error);
      setErr("Could not connect to server.");
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen grid place-items-center p-6">
      <form onSubmit={onSubmit} className="w-full max-w-md space-y-4 p-6 rounded-2xl shadow">
        <h1 className="text-2xl font-semibold">Welcome back</h1>
        {err && <p className="text-sm text-red-600">{err}</p>}

        <input className="w-full border p-2 rounded" name="email" type="email" placeholder="Email" required />
        <input className="w-full border p-2 rounded" name="password" type="password" placeholder="Password" required />

        <button disabled={loading} className="w-full p-2 rounded bg-black text-white">
          {loading ? "Signing in..." : "Sign in"}
        </button>

        <p className="text-sm">
          New here? <a href="/auth/register" className="underline">Create an account</a>
        </p>
      </form>
    </div>
  );
}
