"use client";
import Link from "next/link";
import Image from "next/image";
import { motion } from "framer-motion";

export default function Hero() {
  return (
    <section className="relative bg-gradient-to-br from-green-50 via-white to-green-100">
      <div className="max-w-7xl mx-auto px-6 md:px-12 lg:px-16 flex flex-col-reverse md:flex-row items-center min-h-screen gap-12">
        
        {/* Left: Text */}
        <motion.div
          initial={{ opacity: 0, x: -50 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.8 }}
          className="flex-1 text-center md:text-left"
        >
          <h1 className="text-4xl md:text-6xl font-extrabold text-gray-900 leading-tight">
            Achieve Your <span className="text-green-600">Fitness Goals</span>
          </h1>
          <p className="mt-6 text-gray-700 text-lg md:text-xl leading-relaxed max-w-xl mx-auto md:mx-0">
            Track workouts, follow personalized fitness & diet plans, and stay
            motivated with <span className="font-semibold text-green-600">FitTrack+</span>.
          </p>
          <div className="mt-8 flex justify-center md:justify-start gap-4">
            <Link
              href="/auth/register"
              className="px-8 py-4 bg-green-600 text-white text-lg font-semibold rounded-xl shadow-md hover:bg-green-700 transition transform hover:scale-105"
            >
              Get Started
            </Link>
            <Link
              href="/auth/login"
              className="px-8 py-4 border-2 border-green-600 text-green-600 text-lg font-semibold rounded-xl hover:bg-green-50 transition transform hover:scale-105"
            >
              Login
            </Link>
          </div>
        </motion.div>

        {/* Right: Big Image */}
        <motion.div
          initial={{ opacity: 0, x: 50 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.8 }}
          className="flex-1 flex justify-center md:justify-end relative"
        >
          <div className="relative w-full max-w-md md:max-w-lg lg:max-w-xl">
            <Image
              src="/assets/image/gym.png"
              alt="Fitness Training"
              width={600}
              height={500}
              className="rounded-3xl object-cover shadow-2xl border border-gray-200"
              priority
            />
            {/* Decorative Glow Effect */}
            <div className="absolute inset-0 rounded-3xl bg-green-400/20 blur-3xl -z-10"></div>
          </div>
        </motion.div>
      </div>
    </section>
  );
}
