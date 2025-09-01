"use client";
import Link from "next/link";

export default function CTA() {
  return (
    <section className="py-20 bg-green-600 text-white text-center">
      <h2 className="text-4xl font-bold mb-6">Ready to Get Started?</h2>
      <p className="mb-6">Sign up today and start your fitness journey with personalized plans!</p>
      <Link href="/auth/signup" className="px-8 py-4 bg-white text-green-600 font-semibold rounded-lg hover:bg-gray-100 transition">
        Sign Up Now
      </Link>
    </section>
  );
}
