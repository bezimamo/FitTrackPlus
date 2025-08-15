"use client";

import { useState, useEffect } from "react";
import Image from "next/image";
import { FaDumbbell, FaAppleAlt, FaHeartbeat } from "react-icons/fa";

type Plan = {
  id: number;
  type: "workout" | "diet" | "physio";
  title: string;
  description: string;
  image: string;
};

export default function PlansPage() {
  const [plans, setPlans] = useState<Plan[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const mockPlans: Plan[] = [
      {
        id: 1,
        type: "workout",
        title: "Full Body Strength",
        description: "3 sets of squats, push-ups, and pull-ups every day.",
        image: "/assets/image/push-up.png",
      },
      {
        id: 2,
        type: "diet",
        title: "High Protein Diet",
        description: "Include eggs, chicken, beans, and milk in your meals.",
        image: "/assets/image/diet.png",
      },
      {
        id: 3,
        type: "physio",
        title: "Back Pain Relief",
        description: "Daily stretches and posture exercises for 15 minutes.",
        image: "/assets/image/daily-stretches.png",
      },
    ];

    setTimeout(() => {
      setPlans(mockPlans);
      setLoading(false);
    }, 500);
  }, []);

  const getPlanIcon = (type: string) => {
    switch (type) {
      case "workout":
        return <FaDumbbell />;
      case "diet":
        return <FaAppleAlt />;
      case "physio":
        return <FaHeartbeat />;
      default:
        return null;
    }
  };

  if (loading) return <p className="p-8 text-center">Loading plans...</p>;

  return (
    <div className="p-6 md:p-8 space-y-12">
      {/* Hero Section */}
      <div className="relative w-full h-72 md:h-96 rounded-lg overflow-hidden flex items-center justify-center bg-gradient-to-r from-green-400 via-yellow-300 to-red-400">
        <Image
          src="/assets/image/gym.png"
          alt="Gym Hero"
          fill
          className="object-cover opacity-30"
          priority
        />
        <div className="absolute text-center px-4 md:px-0">
          <h1 className="text-2xl md:text-4xl font-extrabold text-white drop-shadow-lg mb-3">
            Your Personalized Plans
          </h1>
          <p className="text-white text-base md:text-lg mb-6 drop-shadow-md">
            Achieve your fitness goals with tailored workout, diet, and physiotherapy plans.
          </p>
          <button className="px-6 py-3 bg-white text-green-600 font-semibold rounded-lg shadow hover:bg-gray-100 transition">
            Explore Plans
          </button>
        </div>
      </div>

      {/* Plans Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
        {plans.map((plan) => (
          <div
            key={plan.id}
            className="bg-white rounded-lg shadow-md hover:shadow-xl hover:scale-105 transition-transform duration-300 overflow-hidden"
          >
            <div className="relative w-full h-60">
              <Image src={plan.image} alt={plan.title} fill className="object-cover" />
            </div>
            <div className="p-4">
              <div className="flex items-center gap-2 text-blue-600 text-xl mb-2">
                {getPlanIcon(plan.type)}
                <h2 className="text-xl font-semibold">{plan.title}</h2>
              </div>
              <p className="text-gray-700 mt-2">{plan.description}</p>
              <span className="inline-block mt-3 px-2 py-1 text-sm bg-blue-100 text-blue-700 rounded">
                {plan.type.toUpperCase()}
              </span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
