"use client";
import { useState, useEffect } from "react";
import type { PlanFilters } from "@/lib/types/plans";

type Props = {
  value: PlanFilters;
  onChange: (f: PlanFilters) => void;
};

export default function PlanFilters({ value, onChange }: Props) {
  const [local, setLocal] = useState<PlanFilters>(value);

  useEffect(() => setLocal(value), [value]);

  return (
    <div className="bg-white/80 backdrop-blur rounded-2xl shadow p-4 md:p-6 flex flex-col md:flex-row gap-4 md:items-end">
      {/* Search */}
      <div className="flex-1">
        <label className="text-sm text-gray-600">Search</label>
        <input
          value={local.q ?? ""}
          onChange={(e) => setLocal({ ...local, q: e.target.value })}
          placeholder="Search plans by name or description..."
          className="w-full mt-1 rounded-xl border border-gray-200 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      {/* Goal Filter */}
      <div className="flex-1">
        <label className="text-sm text-gray-600">Goal</label>
        <select
          value={local.goal_type ?? ""}
          onChange={(e) =>
            setLocal({ ...local, goal_type: e.target.value || undefined })
          }
          className="w-full mt-1 rounded-xl border border-gray-200 px-3 py-2 bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">All Goals</option>
          <option value="lose_weight">Lose Weight</option>
          <option value="gain_muscle">Gain Muscle</option>
          <option value="flexibility">Flexibility</option>
          <option value="rehab">Rehab</option>
        </select>
      </div>

      {/* Plan Type Filter */}
      <div className="flex-1">
        <label className="text-sm text-gray-600">Plan Type</label>
        <select
          value={local.plan_type ?? ""}
          onChange={(e) =>
            setLocal({ ...local, plan_type: e.target.value || undefined })
          }
          className="w-full mt-1 rounded-xl border border-gray-200 px-3 py-2 bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">All Plan Types</option>
          <option value="fitness">Fitness</option>
          <option value="physiotherapy">Physiotherapy</option>
        </select>
      </div>

      {/* Apply Button */}
      <button
        onClick={() => onChange(local)}
        className="px-4 py-2 rounded-xl bg-blue-600 text-white font-medium hover:bg-blue-700 transition"
      >
        Apply
      </button>
    </div>
  );
}
