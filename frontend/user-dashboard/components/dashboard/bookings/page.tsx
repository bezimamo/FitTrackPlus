// src/app/dashboard/bookings/page.tsx
"use client";

import { useState, useEffect } from "react";

type Booking = {
  id: number;
  trainer: string;
  type: "Workout" | "Physio";
  date: string;
  time: string;
  status: "Available" | "Booked";
};

export default function BookingsPage() {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const mockBookings: Booking[] = [
      {
        id: 1,
        trainer: "John Doe",
        type: "Workout",
        date: "2025-08-16",
        time: "10:00 AM",
        status: "Available",
      },
      {
        id: 2,
        trainer: "Jane Smith",
        type: "Physio",
        date: "2025-08-17",
        time: "2:00 PM",
        status: "Booked",
      },
      {
        id: 3,
        trainer: "Mike Johnson",
        type: "Workout",
        date: "2025-08-18",
        time: "6:00 PM",
        status: "Available",
      },
    ];

    setTimeout(() => {
      setBookings(mockBookings);
      setLoading(false);
    }, 500);
  }, []);

  if (loading) return <p className="p-8">Loading bookings...</p>;

  return (
    <div className="p-8 space-y-8">
      <h1 className="text-3xl font-bold">Book Your Sessions</h1>

      <div className="grid md:grid-cols-3 gap-6">
        {bookings.map((booking) => (
          <div
            key={booking.id}
            className={`p-6 rounded-lg shadow-md hover:shadow-xl transition ${
              booking.status === "Booked" ? "bg-gray-200" : "bg-white"
            }`}
          >
            <h2 className="text-xl font-semibold">{booking.type} Session</h2>
            <p className="text-gray-700 mt-2">Trainer: {booking.trainer}</p>
            <p className="text-gray-700">Date: {booking.date}</p>
            <p className="text-gray-700">Time: {booking.time}</p>
            <p
              className={`mt-2 font-medium ${
                booking.status === "Booked"
                  ? "text-red-600"
                  : "text-green-600"
              }`}
            >
              {booking.status}
            </p>
            {booking.status === "Available" && (
              <button className="mt-4 w-full bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600 transition">
                Book Now
              </button>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
