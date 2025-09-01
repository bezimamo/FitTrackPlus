"use client";

import Link from "next/link";
import Image from "next/image";
import { motion } from "framer-motion";

export default function Navbar() {
  return (
    <nav className="ml-64 p-6 md:p-8"> {/* Pushes away from sidebar */}
      <div className="bg-gradient-to-r from-white via-gray-50 to-gray-100 shadow-xl px-6 md:px-12 py-8 flex flex-col md:flex-row items-center justify-between rounded-2xl gap-8">
        
        {/* Left side: Welcome text */}
        <motion.div
          initial={{ opacity: 0, x: -30 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6 }}
          className="flex flex-col max-w-xl text-center md:text-left"
        >
          <h1 className="text-3xl md:text-4xl font-extrabold text-gray-900">
            Welcome to Your Dashboard
          </h1>
          <p className="text-gray-600 text-base md:text-lg mt-3 leading-relaxed">
            Stay on track with your fitness journey.{" "}
            <Link
              href="/auth/login"
              className="underline text-blue-600 font-medium hover:text-blue-800 transition"
            >
              Login
            </Link>{" "}
            or{" "}
            <Link
              href="/auth/register"
              className="underline text-blue-600 font-medium hover:text-blue-800 transition"
            >
              Register
            </Link>{" "}
            to get started.
          </p>
        </motion.div>

        {/* Right side: Hero Image */}
        <motion.div
          initial={{ opacity: 0, x: 30 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6 }}
          className="flex-shrink-0"
        >
          <Image
            src="/assets/image/gym.png"
            alt="Gym Hero"
            width={450}
            height={250}
            className="rounded-2xl object-cover shadow-lg border border-gray-200"
            priority
          />
        </motion.div>
      </div>
    </nav>
  );
}
