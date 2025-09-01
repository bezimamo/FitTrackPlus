"use client";

import Link from "next/link";

export default function DashboardNavbar() {
  return (
    <nav className="fixed top-0 left-0 w-full bg-white shadow-md px-6 py-4 flex justify-between items-center z-50">
      {/* Brand */}
      <div className="text-2xl font-bold text-gray-800">FitTrack+ Dashboard</div>

      {/* Navigation Links */}
      <div className="hidden md:flex space-x-6">
        <Link
          href="/dashboard/profile"
          className="text-gray-700 hover:text-blue-600 font-semibold transition-colors"
        >
          Profile
        </Link>
        <Link
          href="/dashboard/membership"
          className="text-gray-700 hover:text-blue-600 font-semibold transition-colors"
        >
          Membership
        </Link>
        <Link
          href="/dashboard/plans"
          className="text-gray-700 hover:text-blue-600 font-semibold transition-colors"
        >
          Plans
        </Link>
        <Link
          href="/auth/logout"
          className="text-red-600 font-semibold hover:text-red-800 transition-colors"
        >
          Logout
        </Link>
      </div>
    </nav>
  );
}
