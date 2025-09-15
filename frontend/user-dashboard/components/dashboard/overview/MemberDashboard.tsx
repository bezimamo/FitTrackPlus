"use client"

import type { DashboardData } from "@/lib/types/dashboard"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"

export default function MemberDashboard({ data }: { data: DashboardData }) {
  if (!data.member_data) return null

  const { current_plan, progress_summary } = data.member_data

  return (
    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      {/* Current Plan */}
      <Card>
        <CardHeader>
          <CardTitle>Current Plan</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="font-medium">{current_plan.name}</p>
          <p className="text-sm text-muted-foreground capitalize">{current_plan.status}</p>
          <Progress value={current_plan.progress} className="mt-2" />
          <p className="text-xs mt-1">{current_plan.progress}% completed</p>
        </CardContent>
      </Card>

      {/* Progress Summary */}
      <Card>
        <CardHeader>
          <CardTitle>Progress Summary</CardTitle>
        </CardHeader>
        <CardContent>
          <p>Current Weight: {progress_summary.current_weight} kg</p>
          <p>Target Weight: {progress_summary.target_weight} kg</p>
          <p>Weight Lost: {progress_summary.weight_lost} kg</p>
          <Progress value={progress_summary.overall_progress} className="mt-2" />
          <p className="text-xs mt-1">{progress_summary.overall_progress}% progress</p>
        </CardContent>
      </Card>
    </div>
  )
}
