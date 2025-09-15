"use client"

import type { DashboardData } from "@/lib/types/dashboard"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

export default function TrainerDashboard({ data }: { data: DashboardData }) {
  return (
    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <Card>
        <CardHeader>
          <CardTitle>Trainer Overview</CardTitle>
        </CardHeader>
        <CardContent>
          <p>Total Users: {data.stats.total_users}</p>
          <p>Active Plans: {data.stats.active_plans}</p>
          <p>Total Sessions: {data.stats.total_sessions}</p>
        </CardContent>
      </Card>
    </div>
  )
}
