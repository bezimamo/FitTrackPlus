"use client";

import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from "recharts";

const data = [
  { week: "Week 1", weight: 78 },
  { week: "Week 2", weight: 76 },
  { week: "Week 3", weight: 75 },
  { week: "Week 4", weight: 73 },
];

export default function PlanProgress() {
  return (
    <div className="bg-white rounded-xl shadow-md p-6 mt-8">
      <h3 className="text-xl font-semibold text-gray-900 mb-4">Progress Tracking</h3>
      <ResponsiveContainer width="100%" height={250}>
        <LineChart data={data}>
          <XAxis dataKey="week" />
          <YAxis />
          <Tooltip />
          <Line type="monotone" dataKey="weight" stroke="#3b82f6" strokeWidth={3} />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
