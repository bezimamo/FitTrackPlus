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
    <div className="min-h-screen grid place-items-center p-6">
      <form onSubmit={onSubmit} className="w-full max-w-md space-y-4 p-6 rounded-2xl shadow">
        <h1 className="text-2xl font-semibold">Create your account</h1>
        {err && <p className="text-sm text-red-600">{err}</p>}

        <input className="w-full border p-2 rounded" name="first_name" placeholder="First name" required />
        <input className="w-full border p-2 rounded" name="last_name" placeholder="Last name" required />
        <input className="w-full border p-2 rounded" name="email" type="email" placeholder="Email" required />
        <input className="w-full border p-2 rounded" name="password" type="password" placeholder="Password" required />
        <input className="w-full border p-2 rounded" name="phone" placeholder="Phone (optional)" />
        <select className="w-full border p-2 rounded" name="role" defaultValue="member">
          <option value="member">Member</option>
          <option value="trainer">Trainer</option>
          <option value="physio">Physio</option>
          <option value="admin">Admin</option>
        </select>

        <button disabled={loading} className="w-full p-2 rounded bg-black text-white">
          {loading ? "Creating..." : "Sign up"}
        </button>

        <p className="text-sm">
          Already have an account? <a href="/auth/login" className="underline">Log in</a>
        </p>
      </form>
    </div>
  );
}
