"use client"

import { useEffect, useState } from "react"
import { Badge } from "@/components/ui/badge"
import type { DashboardData } from "@/lib/types/dashboard"
import MemberDashboard from "@/components/dashboard/overview/MemberDashboard"
import TrainerDashboard from "@/components/dashboard/overview/TrainerDashboard"
import AdminDashboard from "@/components/dashboard/overview/AdminDashboard"
import QuickActions from "@/components/dashboard/overview/QuickActions"
import RecentActivity from "@/components/dashboard/overview/RecentActivity"
import Notifications from "@/components/dashboard/overview/Notifications"

export default function DashboardPage() {
  const [dashboardData, setDashboardData] = useState<DashboardData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Replace this with API call later
    const mockData: DashboardData = {
      user_role: "member",
      user_info: {
        id: 1,
        first_name: "John",
        last_name: "Doe",
        email: "john@example.com",
        role: "member",
        is_active: true,
      },
      stats: {
        total_users: 150,
        active_plans: 12,
        total_sessions: 24,
        completion_rate: 85,
      },
      recent_activity: [
        {
          id: 1,
          type: "workout_completed",
          description: "Completed Upper Body Strength Training",
          created_at: new Date().toISOString(),
        },
      ],
      quick_actions: [
        {
          id: "book_session",
          title: "Book Session",
          description: "Schedule a training session",
          icon: "calendar",
          url: "/booking",
          color: "green",
        },
      ],
      notifications: [
        {
          id: 1,
          type: "warning",
          title: "Profile Incomplete",
          message: "Complete your profile to access all features",
          is_read: false,
          created_at: new Date().toISOString(),
        },
      ],
      member_data: {
        current_plan: {
          id: 1,
          name: "Beginner Fitness Plan",
          type: "fitness",
          status: "active",
          progress: 65,
        },
        progress_summary: {
          current_weight: 75.0,
          target_weight: 70.0,
          weight_lost: 2.5,
          overall_progress: 45,
        },
        upcoming_sessions: [],
        goals: [],
      },
    }

    setDashboardData(mockData)
    setLoading(false)
  }, [])

  if (loading) {
    return <p className="p-6 text-muted-foreground">Loading...</p>
  }

  if (!dashboardData) {
    return <p className="p-6 text-muted-foreground">Failed to load dashboard</p>
  }

  const { user_info, recent_activity, quick_actions, notifications } = dashboardData

  return (
    <div className="flex-1 space-y-6 p-6">
      {/* Header */}
      <div className="flex items-center justify-between flex-wrap gap-2">
        <div>
          <h1 className="text-2xl md:text-3xl font-bold text-foreground">
            Welcome back, {user_info.first_name} {user_info.last_name}!
          </h1>
          <p className="text-sm md:text-base text-muted-foreground">
            Here's what's happening with your fitness journey today.
          </p>
        </div>
        <Badge variant="secondary" className="capitalize">
          {user_info.role}
        </Badge>
      </div>

      {/* Role-based dashboard */}
      {user_info.role === "member" && <MemberDashboard data={dashboardData} />}
      {user_info.role === "trainer" && <TrainerDashboard data={dashboardData} />}
      {user_info.role === "admin" && <AdminDashboard data={dashboardData} />}

      {/* Common Sections */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <QuickActions actions={quick_actions} />
        <RecentActivity activities={recent_activity} />
        <Notifications notifications={notifications} />
      </div>
    </div>
  )
}
