"use client"

import type { Action } from "@/lib/types/dashboard"
import type { DashboardData } from "@/lib/types/dashboard"
import Link from "next/link"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

export default function QuickActions({ actions }: { actions: Action[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Quick Actions</CardTitle>
      </CardHeader>
      <CardContent className="space-y-2">
        {actions.map((action) => (
          <Link
            key={action.id}
            href={action.url}
            className="block p-2 rounded-md border hover:bg-muted"
          >
            {action.title}
          </Link>
        ))}
      </CardContent>
    </Card>
  )
}
