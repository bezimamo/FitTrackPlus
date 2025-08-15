import { NextResponse } from "next/server";

export async function POST() {
  // Clear the cookie by setting it expired
  const res = NextResponse.redirect(new URL("/auth/login", process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000"));
  res.cookies.set({
    name: "ft_token",
    value: "",
    path: "/",
    expires: new Date(0),
  });
  return res;
}
