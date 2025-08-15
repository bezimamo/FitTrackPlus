"use client";

import React from "react"; 

type StatCardProps = {
  title: string;
  value: string | number;
  icon?: React.ReactNode; 
};

export default function StatCard({ title, value, icon }: StatCardProps) {
  return (
    <div className="bg-white rounded-xl shadow-md p-6 flex items-center gap-4 hover:shadow-xl transition">
      {icon && <div className="text-2xl text-blue-600">{icon}</div>}
      <div>
        <p className="text-gray-500 text-sm">{title}</p>
        <p className="text-2xl font-bold text-gray-900">{value}</p>
      </div>
    </div>
  );
}
