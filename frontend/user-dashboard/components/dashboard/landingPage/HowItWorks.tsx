"use client";

export default function HowItWorks() {
  const steps = [
    "Sign up and complete your profile",
    "Get a personalized fitness plan",
    "Log your daily progress and track improvements",
    "Book sessions with trainers and physiotherapists",
  ];

  return (
    <section id="how" className="py-20 bg-green-50">
      <div className="max-w-4xl mx-auto px-8 text-center">
        <h2 className="text-4xl font-bold text-gray-800 mb-10">How It Works</h2>
        <ol className="space-y-8 text-left md:text-center md:space-y-0 md:flex md:justify-between md:gap-6">
          {steps.map((step, idx) => (
            <li key={idx} className="bg-white p-6 rounded-2xl shadow hover:shadow-lg transition flex-1">
              <span className="text-green-600 font-bold text-xl">{idx + 1}.</span>
              <p className="mt-2 text-gray-700">{step}</p>
            </li>
          ))}
        </ol>
      </div>
    </section>
  );
}
