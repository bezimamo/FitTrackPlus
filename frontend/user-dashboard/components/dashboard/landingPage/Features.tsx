"use client";

export default function Features() {
  const features = [
    {
      title: "Personalized Plans",
      description: "Receive workout, diet, and physiotherapy plans tailored to your fitness goals.",
    },
    {
      title: "Track Your Progress",
      description: "Monitor weight, BMI, and workout completion with visual progress charts.",
    },
    {
      title: "Book Sessions Easily",
      description: "Check trainer availability and book, reschedule, or cancel training sessions with ease.",
    },
    {
      title: "Membership Management",
      description: "View your membership status and payment history directly on the platform.",
    },
  ];

  return (
    <section id="features" className="py-20 bg-white">
      <div className="max-w-7xl mx-auto px-8 md:px-16 text-center">
        <h2 className="text-4xl font-bold text-gray-800 mb-10">Features</h2>
        <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
          {features.map((f, idx) => (
            <div key={idx} className="bg-green-50 p-6 rounded-2xl shadow hover:shadow-lg transition">
              <h3 className="text-xl font-semibold text-gray-800 mb-3">{f.title}</h3>
              <p className="text-gray-600">{f.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
