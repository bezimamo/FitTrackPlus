"use client";

import { useState } from "react";
import Sidebar from "@/components/Sidebar";
import BookingHeader from "@/components/dashboard/bookings/BookingHeader";
import BookingFilters from "@/components/dashboard/bookings/BookingFilters";
import BookingGrid from "@/components/dashboard/bookings/BookingGrid";

// Mock session data
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
    description:
      "One-on-one personal training session focused on strength building and form correction.",
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
    description:
      "High-intensity interval training to boost your cardiovascular fitness and burn calories.",
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
];

export default function BookingPage() {
  const [searchTerm, setSearchTerm] = useState("");
  const [filterType, setFilterType] = useState("all");
  const [filterDate, setFilterDate] = useState("");

  const handleBookSession = (session: typeof mockSessions[0]) => {
    alert(`Successfully booked: ${session.title} with ${session.trainer}`);
  };

  const filteredSessions = mockSessions.filter((session) => {
    const matchesSearch =
      session.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      session.trainer.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesType = filterType === "all" || session.type === filterType;
    const matchesDate = !filterDate || session.date === filterDate;
    return matchesSearch && matchesType && matchesDate;
  });

  return (
    <div className="flex min-h-screen bg-gray-50">
      {/* Sidebar */}
      <aside className="fixed left-0 top-0 h-full w-64">
        <Sidebar />
      </aside>

      {/* Main Content */}
      <main className="flex-1 md:ml-64 p-6">
        <header className="mb-6">
          <BookingHeader />
        </header>

        <div className="mb-6">
          <BookingFilters
            searchTerm={searchTerm}
            setSearchTerm={setSearchTerm}
            filterType={filterType}
            setFilterType={setFilterType}
            filterDate={filterDate}
            setFilterDate={setFilterDate}
          />
        </div>

        <BookingGrid sessions={filteredSessions} onBook={handleBookSession} />

        {filteredSessions.length === 0 && (
          <p className="text-center text-muted-foreground mt-8">
            No sessions found matching your criteria.
          </p>
        )}
      </main>
    </div>
  );
}
