"use client";

import Link from "next/link";

export default function PlansPreview() {
  const plans = [
    { name: "Weight Loss Plan", goal: "Lose Weight", duration: "4 Weeks" },
    { name: "Muscle Gain Plan", goal: "Gain Muscle", duration: "6 Weeks" },
    { name: "Flexibility Plan", goal: "Flexibility", duration: "8 Weeks" },
  ];

  return (
    <section id="plans" className="py-20 bg-white">
      <div className="max-w-7xl mx-auto px-8 text-center">
        <h2 className="text-4xl font-bold text-gray-800 mb-10">Our Plans</h2>
        <div className="grid md:grid-cols-3 gap-8">
          {plans.map((plan, idx) => (
            <div key={idx} className="bg-green-50 p-6 rounded-2xl shadow hover:shadow-lg transition">
              <h3 className="text-xl font-semibold text-gray-800 mb-2">{plan.name}</h3>
              <p className="text-gray-600">{plan.goal}</p>
              <p className="text-gray-600">{plan.duration}</p>
              <Link href="/auth/signup" className="mt-4 inline-block px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition">
                Join Plan
              </Link>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
