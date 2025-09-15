"use client"

import type { DashboardData } from "@/lib/types/dashboard"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

export default function RecentActivity({ activities }: { activities: Activity[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Recent Activity</CardTitle>
      </CardHeader>
      <CardContent className="space-y-2">
        {activities.length === 0 && <p className="text-sm text-muted-foreground">No recent activity</p>}
        {activities.map((act) => (
          <div key={act.id} className="text-sm">
            <p>{act.description}</p>
            <p className="text-xs text-muted-foreground">{new Date(act.created_at).toLocaleString()}</p>
          </div>
        ))}
      </CardContent>
    </Card>
  )
}
