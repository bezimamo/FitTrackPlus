"use client";

import Sidebar from "@/components/Sidebar";
import { StatsOverview } from "@/components/dashboard/progress/stats-overview";
import { WeightChart } from "@/components/dashboard/progress/weight-chart";
import { WorkoutChart } from "@/components/dashboard/progress/workout-chart";
import { GoalsSection } from "@/components/dashboard/progress/goals-section";
import { RecentWorkouts } from "@/components/dashboard/progress/recent-workouts";
import { AchievementsSection } from "@/components/dashboard/progress/achievements-section";

// Mock data for progress tracking
const weightData = [
  { date: "2024-01-01", weight: 180, bmi: 25.2 },
  { date: "2024-01-08", weight: 178, bmi: 24.9 },
  { date: "2024-01-15", weight: 176, bmi: 24.6 },
  { date: "2024-01-22", weight: 175, bmi: 24.5 },
  { date: "2024-01-29", weight: 173, bmi: 24.2 },
];

const workoutData = [
  { week: "Week 1", sessions: 3, duration: 180 },
  { week: "Week 2", sessions: 4, duration: 240 },
  { week: "Week 3", sessions: 3, duration: 195 },
  { week: "Week 4", sessions: 5, duration: 300 },
];

const mockGoals = [
  { id: 1, title: "Lose 10 lbs", target: 10, current: 7, unit: "lbs", deadline: "2024-03-01", category: "Weight Loss" },
  { id: 2, title: "Workout 4x per week", target: 4, current: 3, unit: "sessions", deadline: "2024-02-01", category: "Fitness" },
  { id: 3, title: "Run 5K in under 25 min", target: 25, current: 28, unit: "minutes", deadline: "2024-04-01", category: "Cardio" },
];

const mockAchievements = [
  { id: 1, title: "First Week Complete", description: "Completed your first week of workouts", date: "2024-01-08", icon: "🎯" },
  { id: 2, title: "Weight Loss Milestone", description: "Lost your first 5 pounds", date: "2024-01-15", icon: "⚖️" },
  { id: 3, title: "Consistency Champion", description: "Worked out 4 times in one week", date: "2024-01-22", icon: "🏆" },
];

const mockWorkouts = [
  { id: 1, date: "2024-01-29", type: "Strength Training", duration: 60, calories: 320, exercises: 8 },
  { id: 2, date: "2024-01-27", type: "Cardio", duration: 45, calories: 280, exercises: 5 },
  { id: 3, date: "2024-01-25", type: "HIIT", duration: 30, calories: 250, exercises: 6 },
  { id: 4, date: "2024-01-23", type: "Yoga", duration: 75, calories: 180, exercises: 12 },
];

export default function ProgressPage() {
  const currentWeight = weightData[weightData.length - 1]?.weight || 0;
  const previousWeight = weightData[weightData.length - 2]?.weight || 0;
  const weightChange = currentWeight - previousWeight;
  const totalWorkouts = workoutData.reduce((sum, week) => sum + week.sessions, 0);
  const totalDuration = workoutData.reduce((sum, week) => sum + week.duration, 0);

  return (
    <div className="flex min-h-screen bg-gray-50">
      {/* Sidebar */}
      <aside className="fixed left-0 top-0 h-full w-64">
        <Sidebar />
      </aside>

      {/* Main Content */}
      <main className="flex-1 md:ml-64 p-6">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-foreground mb-2">Progress Tracking</h1>
          <p className="text-muted-foreground text-lg">
            Monitor your fitness journey and celebrate your achievements
          </p>
        </div>

        {/* Stats Overview */}
        <StatsOverview
          currentWeight={currentWeight}
          weightChange={weightChange}
          totalWorkouts={totalWorkouts}
          totalDuration={totalDuration}
        />

        {/* Charts Section */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
          <WeightChart data={weightData} />
          <WorkoutChart data={workoutData} />
        </div>

        {/* Goals and Recent Workouts */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
          <GoalsSection goals={mockGoals} />
          <RecentWorkouts workouts={mockWorkouts} />
        </div>

        {/* Achievements */}
        <AchievementsSection achievements={mockAchievements} />
      </main>
    </div>
  );
}
