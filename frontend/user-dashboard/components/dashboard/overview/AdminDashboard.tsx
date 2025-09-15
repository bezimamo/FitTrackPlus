"use client"

import type { DashboardData } from "@/lib/types/dashboard"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

export default function AdminDashboard({ data }: { data: DashboardData }) {
  return (
    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <Card>
        <CardHeader>
          <CardTitle>Admin Stats</CardTitle>
        </CardHeader>
        <CardContent>
          <p>Total Users: {data.stats.total_users}</p>
          <p>Completion Rate: {data.stats.completion_rate}%</p>
        </CardContent>
      </Card>
    </div>
  )
}
