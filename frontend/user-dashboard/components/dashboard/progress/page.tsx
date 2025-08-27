"use client";

import { useEffect, useState } from "react";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from "recharts";

type ProgressData = {
  week: string;
  completedWorkouts: number;
};

export default function ProgressPage() {
  const [progress, setProgress] = useState<ProgressData[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const mockProgress: ProgressData[] = [
      { week: "Week 1", completedWorkouts: 3 },
      { week: "Week 2", completedWorkouts: 4 },
      { week: "Week 3", completedWorkouts: 5 },
      { week: "Week 4", completedWorkouts: 4 },
      { week: "Week 5", completedWorkouts: 6 },
    ];

    setTimeout(() => {
      setProgress(mockProgress);
      setLoading(false);
    }, 500);
  }, []);

  if (loading) return <p className="p-8">Loading progress...</p>;

  return (
    <div className="p-8 space-y-8">
      <h1 className="text-3xl font-bold">Your Progress</h1>

      {/* Stats Cards */}
      <div className="grid md:grid-cols-4 gap-6">
        <div className="bg-blue-500 text-white p-6 rounded-lg shadow-md hover:shadow-xl transition">
          <h2 className="text-xl font-semibold">Total Workouts</h2>
          <p className="mt-2 text-2xl">{progress.reduce((sum, p) => sum + p.completedWorkouts, 0)}</p>
        </div>
        <div className="bg-green-500 text-white p-6 rounded-lg shadow-md hover:shadow-xl transition">
          <h2 className="text-xl font-semibold">Average / Week</h2>
          <p className="mt-2 text-2xl">{Math.round(progress.reduce((sum, p) => sum + p.completedWorkouts, 0) / progress.length)}</p>
        </div>
        <div className="bg-yellow-500 text-white p-6 rounded-lg shadow-md hover:shadow-xl transition">
          <h2 className="text-xl font-semibold">Best Week</h2>
          <p className="mt-2 text-2xl">{Math.max(...progress.map(p => p.completedWorkouts))}</p>
        </div>
        <div className="bg-purple-500 text-white p-6 rounded-lg shadow-md hover:shadow-xl transition">
          <h2 className="text-xl font-semibold">Progress %</h2>
          <p className="mt-2 text-2xl">{Math.round((progress.reduce((sum, p) => sum + p.completedWorkouts, 0) / (progress.length * 7)) * 100)}%</p>
        </div>
      </div>

      {/* Progress Chart */}
      <div className="bg-white p-6 rounded-xl shadow-md">
        <h2 className="text-xl font-semibold mb-4">Weekly Workouts</h2>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={progress} margin={{ top: 10, right: 20, bottom: 0, left: 0 }}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="week" />
            <YAxis allowDecimals={false} />
            <Tooltip />
            <Line type="monotone" dataKey="completedWorkouts" stroke="#6B46C1" strokeWidth={3} activeDot={{ r: 8 }} />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
