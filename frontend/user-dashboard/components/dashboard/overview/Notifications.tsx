"use client"

import type { AppNotification } from "@/lib/types/dashboard"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

export default function Notifications({ notifications }: { notifications: AppNotification[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Notifications</CardTitle>
      </CardHeader>
      <CardContent className="space-y-2">
        {notifications.length === 0 && (
          <p className="text-sm text-muted-foreground">No notifications</p>
        )}
        {notifications.map((note) => (
          <div
            key={note.id}
            className="p-2 border rounded-md text-sm"
          >
            <p className="font-medium">{note.title}</p>
            <p className="text-xs text-muted-foreground">{note.message}</p>
          </div>
        ))}
      </CardContent>
    </Card>
  )
}
