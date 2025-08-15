import { NextResponse } from "next/server";

export async function POST(req: Request) {
  const body = await req.json();
  const api = process.env.NEXT_PUBLIC_API_BASE_URL!;
  const cookieName = process.env.COOKIE_NAME || "ft_token";

  const res = await fetch(`${api}/auth/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  const data = await res.json();

  if (!res.ok) {
    return NextResponse.json(
      { error: data?.error || "Registration failed" },
      { status: res.status }
    );
  }

  const resp = NextResponse.json({
    user: data.user,
    expiresAt: data.expires_at,
  });

  // auto-login after register
  resp.cookies.set(cookieName, data.token, {
    httpOnly: true,
    sameSite: "lax",
    secure: false,
    path: "/",
    maxAge: 60 * 60 * 24,
  });

  return resp;
}
