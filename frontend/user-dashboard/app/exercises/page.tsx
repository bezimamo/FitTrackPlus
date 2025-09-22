"use client"

import { AppSidebar } from "@/components/app-sidebar"
import { ExerciseLibrary } from "@/components/dashboard/exercises/exercise-library"

export default function ExercisesPage() {
  return (
    <div className="flex min-h-screen bg-gray-50">
      {/* Sidebar */}
      <aside className="fixed left-0 top-0 h-full w-64">
        <AppSidebar />
      </aside>

      {/* Main Content */}
      <main className="flex-1 md:ml-64 p-6">
        <ExerciseLibrary />
      </main>
    </div>
  )
}
