"use client";

import { useState, useEffect } from "react";

type Membership = {
  planName: string;
  status: string;
  startDate: string;
  endDate: string;
  nextPayment: string;
};

export default function MembershipPage() {
  const [membership, setMembership] = useState<Membership | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Mock membership data
    const mockMembership: Membership = {
      planName: "Premium Plan",
      status: "Active",
      startDate: "2025-01-01",
      endDate: "2025-12-31",
      nextPayment: "2025-09-30",
    };

    setTimeout(() => {
      setMembership(mockMembership);
      setLoading(false);
    }, 500);
  }, []);

  if (loading) return <p className="p-8">Loading membership...</p>;

  return (
    <div className="p-8 max-w-5xl mx-auto space-y-8">
      <h1 className="text-3xl font-bold">Membership Details</h1>

      <div className="bg-white rounded-xl shadow-md p-6 grid md:grid-cols-2 gap-6">
        <div>
          <h2 className="text-lg font-semibold text-gray-700">Plan Name</h2>
          <p className="text-xl font-bold text-gray-900">{membership?.planName}</p>
        </div>

        <div>
          <h2 className="text-lg font-semibold text-gray-700">Status</h2>
          <p className={`text-xl font-bold ${
            membership?.status === "Active" ? "text-green-500" : "text-red-500"
          }`}>
            {membership?.status}
          </p>
        </div>

        <div>
          <h2 className="text-lg font-semibold text-gray-700">Start Date</h2>
          <p className="text-xl text-gray-900">{membership?.startDate}</p>
        </div>

        <div>
          <h2 className="text-lg font-semibold text-gray-700">End Date</h2>
          <p className="text-xl text-gray-900">{membership?.endDate}</p>
        </div>

        <div className="md:col-span-2">
          <h2 className="text-lg font-semibold text-gray-700">Next Payment</h2>
          <p className="text-xl text-gray-900">{membership?.nextPayment}</p>
        </div>
      </div>

      {/* Optional: Add a button to renew or upgrade */}
      <div className="text-right">
        <button className="px-6 py-3 bg-blue-600 text-white rounded-lg shadow hover:bg-blue-700 transition">
          Renew / Upgrade
        </button>
      </div>
    </div>
  );
}
