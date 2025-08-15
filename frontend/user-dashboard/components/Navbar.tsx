"use client";
import Link from "next/link";

export default function Navbar() {
  return (
    <nav className="bg-white shadow-md px-8 py-4 flex justify-between items-center rounded-b-lg">
      {/* Left side: Brand / Dashboard title */}
      <div className="text-2xl font-bold text-gray-800">FitTrack+</div>

      {/* Right side: Navigation links */}
      <div className="flex space-x-6">
        <Link
          href="/dashboard/profile"
          className="text-gray-700 hover:text-green-600 font-semibold transition-colors"
        >
          Profile
        </Link>
        <Link
          href="/dashboard/membership"
          className="text-gray-700 hover:text-green-600 font-semibold transition-colors"
        >
          Membership
        </Link>
        <Link
          href="/dashboard/about"
          className="text-gray-700 hover:text-green-600 font-semibold transition-colors"
        >
          About
        </Link>
      </div>
    </nav>
  );
}
