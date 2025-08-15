// src/components/Sidebar.tsx
"use client";

import Link from "next/link";
import { FaHome, FaDumbbell, FaChartLine, FaCalendarAlt, FaUser } from "react-icons/fa";

export default function Sidebar() {
  return (
    <aside className="w-64 bg-gradient-to-b from-blue-600 to-blue-400 text-white flex flex-col p-6 min-h-screen">
      {/* Sidebar Header */}
      <h2 className="text-2xl font-bold mb-8">Dashboard</h2>

      {/* Navigation Links */}
      <nav className="flex flex-col gap-4">
        <Link
          href="/dashboard/overview"
          className="flex items-center gap-3 py-2 px-4 rounded hover:bg-white hover:text-blue-600 transition"
        >
          <FaHome /> Overview
        </Link>
        <Link
          href="/dashboard/plans"
          className="flex items-center gap-3 py-2 px-4 rounded hover:bg-white hover:text-blue-600 transition"
        >
          <FaDumbbell /> Plans
        </Link>
        <Link
          href="/dashboard/progress"
          className="flex items-center gap-3 py-2 px-4 rounded hover:bg-white hover:text-blue-600 transition"
        >
          <FaChartLine /> Progress
        </Link>
        <Link
          href="/dashboard/bookings"
          className="flex items-center gap-3 py-2 px-4 rounded hover:bg-white hover:text-blue-600 transition"
        >
          <FaCalendarAlt /> Bookings
        </Link>
        <Link
          href="/dashboard/profile"
          className="flex items-center gap-3 py-2 px-4 rounded hover:bg-white hover:text-blue-600 transition"
        >
          <FaUser /> Profile
        </Link>
      </nav>
    </aside>
  );
}
