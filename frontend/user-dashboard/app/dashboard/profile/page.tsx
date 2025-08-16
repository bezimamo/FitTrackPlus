"use client";

import { useState, useEffect } from "react";
import Cookies from "js-cookie";
import Image from "next/image";

type Profile = {
  name?: string;
  email?: string;
  phone?: string;
  age?: number;
  gender?: string;
  height?: number; // cm
  weight?: number; // kg
  goal?: string;
  membership?: string;
  beforeImage?: string;
  afterImage?: string;
  lastAfterUpdate?: string; // ISO date
};

export default function ProfilePage() {
  const [profile, setProfile] = useState<Profile>({});
  const [bmi, setBmi] = useState<number | null>(null);
  const [bmiCategory, setBmiCategory] = useState<string>("");

  // Load profile from cookies
  useEffect(() => {
    const saved = Cookies.get("fit_profile");
    if (saved) setProfile(JSON.parse(saved));
  }, []);

  // Save profile to cookies & calculate BMI whenever height or weight changes
  useEffect(() => {
    Cookies.set("fit_profile", JSON.stringify(profile), { expires: 30 });
    if (profile.height && profile.weight) {
      const hM = profile.height / 100;
      const value = profile.weight / (hM * hM);
      setBmi(Number(value.toFixed(1)));
      if (value < 18.5) setBmiCategory("Underweight");
      else if (value < 24.9) setBmiCategory("Normal");
      else if (value < 29.9) setBmiCategory("Overweight");
      else setBmiCategory("Obese");
    } else {
      setBmi(null);
      setBmiCategory("");
    }
  }, [profile.height, profile.weight]);

  const isProfileComplete = profile.name && profile.email && profile.height && profile.weight;

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>, type: "before" | "after") => {
    const file = e.target.files?.[0];
    if (!file) return;
    const url = URL.createObjectURL(file);
    if (type === "before" && !profile.beforeImage) {
      setProfile((p) => ({ ...p, beforeImage: url }));
    } else {
      setProfile((p) => ({ ...p, afterImage: url, lastAfterUpdate: new Date().toISOString() }));
    }
  };

  // Refresh afterImage every 15 days
  useEffect(() => {
    if (!profile.lastAfterUpdate) return;
    const lastUpdate = new Date(profile.lastAfterUpdate);
    const now = new Date();
    const diff = Math.floor((now.getTime() - lastUpdate.getTime()) / (1000 * 60 * 60 * 24));
    if (diff >= 15) setProfile((p) => ({ ...p, afterImage: undefined }));
  }, [profile.lastAfterUpdate]);

  return (
    <div className="p-8 max-w-5xl mx-auto space-y-8">
      {/* Profile Completion Banner */}
      {!isProfileComplete && (
        <div className="bg-yellow-100 border-l-4 border-yellow-500 text-yellow-800 p-4 rounded-lg shadow">
          <p className="font-semibold">Complete your profile</p>
          <p className="text-sm">Please add your name, email, height, and weight to unlock full features.</p>
        </div>
      )}

      {/* Static UserProfile Card */}
      <div className="max-w-2xl mx-auto bg-white shadow-lg rounded-2xl p-6">
        {/* Profile Image */}
        <div className="flex flex-col items-center">
          <Image
            src={profile.beforeImage || "/assets/users/profile.jpg"}
            alt="Profile"
            width={120}
            height={120}
            className="rounded-full border-4 border-blue-500 shadow-md"
          />
          <h2 className="text-2xl font-bold mt-3">{profile.name || "Bezawit Mamo"}</h2>
          <p className="text-gray-500">{profile.email || "bezawitmamo27@gmail.com"}</p>
        </div>

        {/* Personal Information */}
        <div className="mt-6 grid grid-cols-2 gap-4 text-sm">
          <div className="p-3 bg-blue-50 rounded-lg shadow-sm">
            <span className="font-semibold">Phone:</span> {profile.phone || "+251914370232"}
          </div>
          <div className="p-3 bg-blue-50 rounded-lg shadow-sm">
            <span className="font-semibold">Age:</span> {profile.age || 22}
          </div>
          <div className="p-3 bg-blue-50 rounded-lg shadow-sm">
            <span className="font-semibold">Gender:</span> {profile.gender || "Female"}
          </div>
          <div className="p-3 bg-blue-50 rounded-lg shadow-sm">
            <span className="font-semibold">Height:</span>
            <input
              type="number"
              value={profile.height || ""}
              placeholder="cm"
              className="ml-2 w-20 border rounded p-1 text-sm"
              onChange={(e) => setProfile((p) => ({ ...p, height: Number(e.target.value) }))}
            />
          </div>
          <div className="p-3 bg-blue-50 rounded-lg shadow-sm">
            <span className="font-semibold">Weight:</span>
            <input
              type="number"
              value={profile.weight || ""}
              placeholder="kg"
              className="ml-2 w-20 border rounded p-1 text-sm"
              onChange={(e) => setProfile((p) => ({ ...p, weight: Number(e.target.value) }))}
            />
          </div>
          <div className="p-3 bg-blue-50 rounded-lg shadow-sm col-span-2">
            <span className="font-semibold">Goal:</span> {profile.goal || "Stay fit and healthy"}
          </div>
          <div className="p-3 bg-blue-50 rounded-lg shadow-sm col-span-2">
            <span className="font-semibold">Membership:</span> {profile.membership || "Premium"}
          </div>
        </div>
      </div>

      {/* Weight, Height & BMI Card */}
      <div className="bg-white p-6 rounded-xl shadow-md grid md:grid-cols-3 gap-6 text-center">
        <div className="p-4 bg-blue-50 rounded-lg shadow hover:shadow-lg transition">
          <p className="font-semibold text-gray-600">Weight</p>
          <p className="text-2xl font-bold">{profile.weight ? `${profile.weight} kg` : "-"}</p>
        </div>
        <div className="p-4 bg-green-50 rounded-lg shadow hover:shadow-lg transition">
          <p className="font-semibold text-gray-600">Height</p>
          <p className="text-2xl font-bold">{profile.height ? `${profile.height} cm` : "-"}</p>
        </div>
        <div className="p-4 bg-purple-50 rounded-lg shadow hover:shadow-lg transition">
          <p className="font-semibold text-gray-600">BMI</p>
          <p className="text-2xl font-bold">{bmi || "-"}</p>
          <p className="text-sm text-gray-500">{bmiCategory || ""}</p>
        </div>
      </div>

      {/* Progress Images */}
      <div className="bg-white p-6 rounded-xl shadow-md space-y-6">
        <h2 className="text-xl font-bold">Progress Tracker</h2>
        <div className="grid md:grid-cols-2 gap-6">
          <div className="text-center">
            <p className="font-semibold mb-2">Before</p>
            {profile.beforeImage ? (
              <img src={profile.beforeImage} alt="Before" className="rounded-lg shadow-md w-full h-64 object-cover" />
            ) : (
              <label className="cursor-pointer block p-6 border-2 border-dashed rounded-lg text-gray-500 hover:bg-gray-50">
                Upload Before Image
                <input type="file" accept="image/*" className="hidden" onChange={(e) => handleImageUpload(e, "before")} />
              </label>
            )}
          </div>

          <div className="text-center">
            <p className="font-semibold mb-2">After</p>
            {profile.afterImage ? (
              <img src={profile.afterImage} alt="After" className="rounded-lg shadow-md w-full h-64 object-cover" />
            ) : (
              <label className="cursor-pointer block p-6 border-2 border-dashed rounded-lg text-gray-500 hover:bg-gray-50">
                Upload After Image (every 15 days)
                <input type="file" accept="image/*" className="hidden" onChange={(e) => handleImageUpload(e, "after")} />
              </label>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
