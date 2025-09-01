"use client";

import Link from "next/link";

export default function Navbar() {
  return (
    <nav className="fixed top-0 left-0 w-full bg-white shadow-md px-6 py-4 flex justify-between items-center z-50">
      {/* Brand */}
      <div className="text-2xl font-bold text-gray-800">FitTrack+</div>

      {/* Navigation Links */}
      <div className="hidden md:flex space-x-6">
        <Link
          href="/auth/login"
          className="text-gray-700 hover:text-green-600 font-semibold transition-colors"
        >
          Login
        </Link>
        <Link
          href="/auth/register"
          className="text-gray-700 hover:text-green-600 font-semibold transition-colors"
        >
          Register
        </Link>
      </div>
    </nav>
  );
}
