"use client"

import { Calendar, Clock, User, MapPin } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import BookingDialog from "./BookingDialog"

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
  session: Session
  onBook: (session: Session) => void
}

export default function BookingCard({ session, onBook }: Props) {
  return (
    <Card className={`transition-all hover:shadow-lg ${!session.available ? "opacity-60" : ""}`}>
      <CardHeader>
        <div className="flex justify-between items-start">
          <div>
            <CardTitle className="text-xl mb-1">{session.title}</CardTitle>
            <CardDescription className="flex items-center gap-1">
              <User className="h-4 w-4" />
              {session.trainer}
            </CardDescription>
          </div>
     <Badge
  className={`${
    session.available
      ? "bg-[var(--primary)] text-[var(--primary-foreground)]"
      : "bg-[var(--destructive)] text-[var(--destructive-foreground)]"
  }`}
>
  {session.available ? "Available" : "Full"}
</Badge>


        </div>
      </CardHeader>

      <CardContent className="space-y-4">
        <p className="text-sm text-muted-foreground">{session.description}</p>

        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm">
            <Calendar className="h-4 w-4 text-primary" />
            <span>{new Date(session.date).toLocaleDateString()}</span>
          </div>
          <div className="flex items-center gap-2 text-sm">
            <Clock className="h-4 w-4 text-primary" />
            <span>{session.time} ({session.duration} min)</span>
          </div>
          <div className="flex items-center gap-2 text-sm">
            <MapPin className="h-4 w-4 text-primary" />
            <span>{session.location}</span>
          </div>
        </div>

        <div className="flex justify-between items-center pt-2">
          <div>
            <span className="text-2xl font-bold text-primary">${session.price}</span>
            {session.type === "Group Class" && (
              <p className="text-xs text-muted-foreground">
                {session.currentParticipants}/{session.maxParticipants} spots filled
              </p>
            )}
          </div>

          {/* Pass onBook to BookingDialog as "onBook" prop */}
          <BookingDialog session={session} onBook={onBook} />
        </div>
      </CardContent>
    </Card>
  )
}
