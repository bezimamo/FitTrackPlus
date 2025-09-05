"use client"

import { Card, CardContent } from "@/components/ui/card"
import BookingCard from "./BookingCard"

interface Session {
  id: number
  title: string
  trainer: string
  type: string
  duration: number
  price: number
  date: string
  time: string
  location: string
  description: string
  available: boolean
  maxParticipants: number
  currentParticipants: number
}

interface Props {
  sessions: Session[]
  onBook: (session: Session) => void
}

export default function BookingGrid({ sessions, onBook }: Props) {
  if (sessions.length === 0) {
    return (
      <Card className="text-center py-12">
        <CardContent>
          <p className="text-muted-foreground text-lg">No sessions found matching your criteria.</p>
          <p className="text-sm text-muted-foreground mt-2">
            Try adjusting your filters or search terms.
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {sessions.map((session) => (
        <BookingCard key={session.id} session={session} onBook={onBook} />
      ))}
    </div>
  )
}
