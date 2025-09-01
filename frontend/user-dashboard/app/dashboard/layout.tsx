// app/dashboard/layout.tsx
"use client";

import DashboardNavbar from "@/components/DashboardNavbar";

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen flex flex-col">
      {/* Navbar */}
      <DashboardNavbar />

      {/* Page Content */}
      <main className="flex-1 mt-16 px-6">{children}</main>
    </div>
  );
}
