"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"

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

export default function BookingDialog({ session, onBook }: Props) {
  const [notes, setNotes] = useState("")

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button disabled={!session.available}>{session.available ? "Book Now" : "Full"}</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Book Session</DialogTitle>
          <DialogDescription>Confirm your booking for {session.title}</DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <Label className="font-medium">Trainer</Label>
              <p>{session.trainer}</p>
            </div>
            <div>
              <Label className="font-medium">Date & Time</Label>
              <p>
                {new Date(session.date).toLocaleDateString()} at {session.time}
              </p>
            </div>
            <div>
              <Label className="font-medium">Duration</Label>
              <p>{session.duration} minutes</p>
            </div>
            <div>
              <Label className="font-medium">Price</Label>
              <p className="text-primary font-bold">${session.price}</p>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="notes">Special Notes (Optional)</Label>
            <Textarea
              id="notes"
              placeholder="Any special requirements or notes for your trainer..."
              className="resize-none"
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
            />
          </div>

          <div className="flex gap-2 pt-4">
            <Button className="flex-1" onClick={() => onBook(session)}>
              Confirm Booking
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
