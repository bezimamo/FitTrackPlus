"use client";

import { useState, useEffect } from "react";

type UserProfile = {
  name: string;
  email: string;
  membershipStatus: string;
  nextPaymentDue: string;
};

export default function ProfilePage() {
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Mock data for now
    const mockProfile: UserProfile = {
      name: "John Doe",
      email: "john@example.com",
      membershipStatus: "Active",
      nextPaymentDue: "2025-09-30",
    };

    setTimeout(() => {
      setProfile(mockProfile);
      setLoading(false);
    }, 500);
  }, []);

  if (loading) return <p className="p-8">Loading profile...</p>;

  return (
    <div className="p-8 space-y-8 max-w-5xl mx-auto">
      {/* Profile Card */}
      <div className="bg-white rounded-xl shadow-md p-6 flex flex-col md:flex-row items-center gap-6">
        <div className="flex-shrink-0">
          <div className="w-24 h-24 bg-gray-200 rounded-full flex items-center justify-center text-2xl font-bold text-gray-600">
            {profile?.name[0]}
          </div>
        </div>
        <div className="flex-1">
          <h1 className="text-2xl font-bold">{profile?.name}</h1>
          <p className="text-gray-600">{profile?.email}</p>
        </div>
      </div>

      {/* Membership & Payment Section */}
      <div className="grid md:grid-cols-3 gap-6">
        {/* Membership Status */}
        <div className="bg-green-500 text-white p-6 rounded-xl shadow-md hover:shadow-xl transition">
          <h2 className="text-lg font-semibold">Membership Status</h2>
          <p className="mt-2 text-2xl">{profile?.membershipStatus}</p>
        </div>

        {/* Next Payment Due */}
        <div className="bg-yellow-500 text-white p-6 rounded-xl shadow-md hover:shadow-xl transition">
          <h2 className="text-lg font-semibold">Next Payment Due</h2>
          <p className="mt-2 text-2xl">{profile?.nextPaymentDue}</p>
        </div>

        {/* Payment Progress */}
        <div className="bg-purple-500 text-white p-6 rounded-xl shadow-md hover:shadow-xl transition">
          <h2 className="text-lg font-semibold">Payment Progress</h2>
          <div className="mt-3 bg-purple-700 h-4 rounded-full overflow-hidden">
            <div
              className="bg-white h-4 rounded-full"
              style={{ width: "70%" }} // mock progress
            ></div>
          </div>
          <p className="mt-2 text-white text-sm">70% completed</p>
        </div>
      </div>
    </div>
  );
}
