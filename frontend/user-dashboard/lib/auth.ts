// lib/auth.ts
export function isLoggedIn(): boolean {
  if (typeof window === "undefined") return false;
  return document.cookie.split(";").some((c) => c.trim().startsWith("ft_token="));
}
