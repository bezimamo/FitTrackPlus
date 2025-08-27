import { NextResponse, NextRequest } from "next/server";

export async function GET(req: NextRequest) {
  const cookieName = process.env.COOKIE_NAME || "ft_token";
  const api = process.env.NEXT_PUBLIC_API_BASE_URL!;

  // Read token from cookie
  const token = req.cookies.get(cookieName)?.value;

  if (!token) {
    return NextResponse.json({ error: "Unauthorized" }, { status: 401 });
  }

  // Call backend API
  const res = await fetch(`${api}/users/profile`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${token}`,
    },
  });

  const data = await res.json();

  if (!res.ok) {
    return NextResponse.json({ error: data?.error || "Failed to load profile" }, { status: res.status });
  }

  return NextResponse.json(data);
}
