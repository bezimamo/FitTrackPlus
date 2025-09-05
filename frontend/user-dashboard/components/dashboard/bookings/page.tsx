"use client"

import { useState } from "react"
import { Calendar, Clock, User, MapPin, Filter, Search } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
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

// Mock data for sessions
const mockSessions = [
  {
    id: 1,
    title: "Personal Training",
    trainer: "Sarah Johnson",
    type: "Personal Training",
    duration: 60,
    price: 75,
    date: "2024-01-15",
    time: "09:00",
    location: "Gym Floor A",
    description: "One-on-one personal training session focused on strength building and form correction.",
    available: true,
    maxParticipants: 1,
    currentParticipants: 0,
  },
  {
    id: 2,
    title: "HIIT Cardio Blast",
    trainer: "Mike Chen",
    type: "Group Class",
    duration: 45,
    price: 25,
    date: "2024-01-15",
    time: "10:30",
    location: "Studio B",
    description: "High-intensity interval training to boost your cardiovascular fitness and burn calories.",
    available: true,
    maxParticipants: 15,
    currentParticipants: 8,
  },
  {
    id: 3,
    title: "Yoga Flow",
    trainer: "Emma Davis",
    type: "Group Class",
    duration: 75,
    price: 20,
    date: "2024-01-15",
    time: "18:00",
    location: "Studio A",
    description: "Relaxing yoga session to improve flexibility and reduce stress.",
    available: true,
    maxParticipants: 20,
    currentParticipants: 12,
  },
  {
    id: 4,
    title: "Physiotherapy Session",
    trainer: "Dr. Alex Rodriguez",
    type: "Physiotherapy",
    duration: 45,
    price: 90,
    date: "2024-01-16",
    time: "14:00",
    location: "Therapy Room 1",
    description: "Rehabilitation session for injury recovery and movement improvement.",
    available: true,
    maxParticipants: 1,
    currentParticipants: 0,
  },
  {
    id: 5,
    title: "Strength Training",
    trainer: "John Smith",
    type: "Group Class",
    duration: 60,
    price: 30,
    date: "2024-01-16",
    time: "07:00",
    location: "Weight Room",
    description: "Build muscle and increase strength with guided weight training.",
    available: false,
    maxParticipants: 10,
    currentParticipants: 10,
  },
  {
    id: 6,
    title: "Pilates Core",
    trainer: "Lisa Wang",
    type: "Group Class",
    duration: 50,
    price: 25,
    date: "2024-01-17",
    time: "12:00",
    location: "Studio C",
    description: "Core-focused Pilates session to improve stability and posture.",
    available: true,
    maxParticipants: 12,
    currentParticipants: 5,
  },
]

export default function BookingPage() {
  const [selectedSession, setSelectedSession] = useState<(typeof mockSessions)[0] | null>(null)
  const [searchTerm, setSearchTerm] = useState("")
  const [filterType, setFilterType] = useState("all")
  const [filterDate, setFilterDate] = useState("")

  const filteredSessions = mockSessions.filter((session) => {
    const matchesSearch =
      session.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      session.trainer.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesType = filterType === "all" || session.type === filterType
    const matchesDate = !filterDate || session.date === filterDate

    return matchesSearch && matchesType && matchesDate
  })

  const handleBookSession = (session: (typeof mockSessions)[0]) => {
    // Mock booking logic
    alert(`Successfully booked: ${session.title} with ${session.trainer}`)
  }

  return (
    <div className="min-h-screen bg-background p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-foreground mb-2">Book Your Session</h1>
          <p className="text-muted-foreground text-lg">
            Choose from personal training, group classes, or physiotherapy sessions
          </p>
        </div>

        {/* Filters */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Filter className="h-5 w-5" />
              Filter Sessions
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label htmlFor="search">Search</Label>
                <div className="relative">
                  <Search className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                  <Input
                    id="search"
                    placeholder="Search sessions or trainers..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="pl-10"
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="type">Session Type</Label>
                <Select value={filterType} onValueChange={setFilterType}>
                  <SelectTrigger>
                    <SelectValue placeholder="All types" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Types</SelectItem>
                    <SelectItem value="Personal Training">Personal Training</SelectItem>
                    <SelectItem value="Group Class">Group Class</SelectItem>
                    <SelectItem value="Physiotherapy">Physiotherapy</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="date">Date</Label>
                <Input id="date" type="date" value={filterDate} onChange={(e) => setFilterDate(e.target.value)} />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Sessions Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredSessions.map((session) => (
            <Card
              key={session.id}
              className={`transition-all hover:shadow-lg ${!session.available ? "opacity-60" : ""}`}
            >
              <CardHeader>
                <div className="flex justify-between items-start">
                  <div>
                    <CardTitle className="text-xl mb-1">{session.title}</CardTitle>
                    <CardDescription className="flex items-center gap-1">
                      <User className="h-4 w-4" />
                      {session.trainer}
                    </CardDescription>
                  </div>
                  <Badge variant={session.available ? "default" : "secondary"}>
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
                    <span>
                      {session.time} ({session.duration} min)
                    </span>
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

                  <Dialog>
                    <DialogTrigger asChild>
                      <Button disabled={!session.available} onClick={() => setSelectedSession(session)}>
                        {session.available ? "Book Now" : "Full"}
                      </Button>
                    </DialogTrigger>
                    <DialogContent className="sm:max-w-md">
                      <DialogHeader>
                        <DialogTitle>Book Session</DialogTitle>
                        <DialogDescription>Confirm your booking for {selectedSession?.title}</DialogDescription>
                      </DialogHeader>
                      {selectedSession && (
                        <div className="space-y-4">
                          <div className="grid grid-cols-2 gap-4 text-sm">
                            <div>
                              <Label className="font-medium">Trainer</Label>
                              <p>{selectedSession.trainer}</p>
                            </div>
                            <div>
                              <Label className="font-medium">Date & Time</Label>
                              <p>
                                {new Date(selectedSession.date).toLocaleDateString()} at {selectedSession.time}
                              </p>
                            </div>
                            <div>
                              <Label className="font-medium">Duration</Label>
                              <p>{selectedSession.duration} minutes</p>
                            </div>
                            <div>
                              <Label className="font-medium">Price</Label>
                              <p className="text-primary font-bold">${selectedSession.price}</p>
                            </div>
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="notes">Special Notes (Optional)</Label>
                            <Textarea
                              id="notes"
                              placeholder="Any special requirements or notes for your trainer..."
                              className="resize-none"
                            />
                          </div>

                          <div className="flex gap-2 pt-4">
                            <Button className="flex-1" onClick={() => handleBookSession(selectedSession)}>
                              Confirm Booking
                            </Button>
                          </div>
                        </div>
                      )}
                    </DialogContent>
                  </Dialog>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {filteredSessions.length === 0 && (
          <Card className="text-center py-12">
            <CardContent>
              <p className="text-muted-foreground text-lg">No sessions found matching your criteria.</p>
              <p className="text-sm text-muted-foreground mt-2">Try adjusting your filters or search terms.</p>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  )
}
