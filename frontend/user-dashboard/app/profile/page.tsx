"use client";

import { useEffect, useState } from "react";
import apiFetch from "@/lib/api";
import { UserProfileResponse } from "@/lib/types/profile";
import { AppSidebar } from "@/components/app-sidebar"
import ProfileHeader from "@/components/dashboard/profile/ProfileHeader";
import ProfileDetails from "@/components/dashboard/profile/ProfileDetails";
import ProfileCompletion from "@/components/dashboard/profile/ProfileCompletion";
import BMICard from "@/components/dashboard/profile/BMICard";
import RoleProfile from "@/components/dashboard/profile/RoleProfile";

export default function ProfilePage() {
  const [profileData, setProfileData] = useState<UserProfileResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const token = document.cookie
      .split("; ")
      .find(row => row.startsWith("ft_token="))
      ?.split("=")[1];

    if (!token) {
      setError("You are not logged in");
      setLoading(false);
      return;
    }

    apiFetch
      .get<UserProfileResponse>("/users/profile", { headers: { Authorization: `Bearer ${token}` } })
      .then(res => setProfileData(res.data))
      .catch(() => setError("Failed to load profile"))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-center mt-10">Loading profile...</p>;
  if (error) return <p className="text-center mt-10 text-red-600">{error}</p>;
  if (!profileData) return <p className="text-center mt-10">No profile data found</p>;

  const { profile, completion_percentage, is_complete } = profileData;

  return (
    <div className="flex">
      {/* Sidebar fixed */}
      <div className="w-64 flex-shrink-0">
        <AppSidebar />
      </div>

      {/* Scrollable Main Content */}
      <main className="flex-1 overflow-auto h-screen p-6 bg-gray-50">
        <div className="max-w-5xl mx-auto space-y-6">
          <ProfileHeader
            name={`User #${profile.user_id}`}
            profileImage={profile.profile_image_url}
            isComplete={profile.is_profile_complete}
          />

          <ProfileCompletion completion={completion_percentage} isComplete={is_complete} />

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <ProfileDetails profile={profile} />
            <BMICard height={profile.height} weight={profile.weight} />
          </div>

          <RoleProfile profile={profile} />
        </div>
      </main>
    </div>
  );
}
