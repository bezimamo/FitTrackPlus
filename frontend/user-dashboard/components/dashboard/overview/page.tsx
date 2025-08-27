"use client";

import { useState, useEffect } from "react";
import Image from "next/image";
import { FaDumbbell, FaCheckCircle, FaCalendarAlt, FaChartLine } from "react-icons/fa";

type Overview = {
  totalWorkouts: number;
  completedWorkouts: number;
  upcomingSessions: number;
  progressPercentage: number;
};

export default function OverviewPage() {
  const [overview, setOverview] = useState<Overview | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const mockOverview: Overview = {
      totalWorkouts: 25,
      completedWorkouts: 18,
      upcomingSessions: 3,
      progressPercentage: 72,
    };

    setTimeout(() => {
      setOverview(mockOverview);
      setLoading(false);
    }, 500);
  }, []);

  if (loading) return <p className="p-8">Loading overview...</p>;

  return (
    <div className="p-8 space-y-8">
      {/* Hero Section */}
      <div className="relative w-full h-64 rounded-xl overflow-hidden shadow-lg">
        <Image
          src="/assets/image/gym.png"
          alt="Gym Hero"
          fill
          className="object-cover"
          priority
        />
        <div className="absolute inset-0 bg-black bg-opacity-40 flex items-center justify-center">
          <h1 className="text-4xl font-bold text-white text-center">
            Welcome to Your Fitness Dashboard
          </h1>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid md:grid-cols-4 gap-6 mt-6">
        <div className="bg-blue-50 text-blue-800 p-6 rounded-xl shadow hover:shadow-xl transition flex items-center gap-4">
          <FaDumbbell className="text-3xl" />
          <div>
            <h2 className="text-lg font-semibold">Total Workouts</h2>
            <p className="text-2xl font-bold">{overview?.totalWorkouts}</p>
          </div>
        </div>

        <div className="bg-green-50 text-green-800 p-6 rounded-xl shadow hover:shadow-xl transition flex items-center gap-4">
          <FaCheckCircle className="text-3xl" />
          <div>
            <h2 className="text-lg font-semibold">Completed</h2>
            <p className="text-2xl font-bold">{overview?.completedWorkouts}</p>
          </div>
        </div>

        <div className="bg-yellow-50 text-yellow-800 p-6 rounded-xl shadow hover:shadow-xl transition flex items-center gap-4">
          <FaCalendarAlt className="text-3xl" />
          <div>
            <h2 className="text-lg font-semibold">Upcoming Sessions</h2>
            <p className="text-2xl font-bold">{overview?.upcomingSessions}</p>
          </div>
        </div>

        <div className="bg-purple-50 text-purple-800 p-6 rounded-xl shadow hover:shadow-xl transition flex items-center gap-4">
          <FaChartLine className="text-3xl" />
          <div>
            <h2 className="text-lg font-semibold">Progress</h2>
            <p className="text-2xl font-bold">{overview?.progressPercentage}%</p>
          </div>
        </div>
      </div>

      {/* Overall Progress Bar */}
      <div className="mt-6">
        <h2 className="text-xl font-semibold mb-2">Overall Progress</h2>
        <div className="w-full bg-gray-200 rounded-full h-6">
          <div
            className="bg-purple-600 h-6 rounded-full transition-all duration-500"
            style={{ width: `${overview?.progressPercentage}%` }}
          />
        </div>
      </div>
    </div>
  );
}
