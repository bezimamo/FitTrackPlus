"use client";

import { Card, CardContent } from "@/components/ui/card";
import { FaWeight, FaRulerVertical } from "react-icons/fa";

export default function BMICard({ weight, height }: { weight: number; height: number }) {
  const bmi = (weight / ((height / 100) * (height / 100))).toFixed(1);

  let status = "Normal";
  let bgColor = "from-green-100 to-green-50";
  if (parseFloat(bmi) < 18.5) {
    status = "Underweight";
    bgColor = "from-blue-100 to-blue-50";
  } else if (parseFloat(bmi) >= 25) {
    status = "Overweight";
    bgColor = "from-red-100 to-red-50";
  }

  return (
    <Card className={`rounded-2xl shadow-xl p-6 bg-gradient-to-r ${bgColor} hover:shadow-2xl transition-shadow duration-300`}>
      <h3 className="text-2xl font-bold text-gray-800 border-b pb-2 mb-4">BMI Information</h3>
      <CardContent className="space-y-3">
        <div className="flex items-center gap-3">
          <FaWeight className="text-gray-600 w-5 h-5" />
          <span className="text-gray-700 font-medium">Weight: {weight} kg</span>
        </div>
        <div className="flex items-center gap-3">
          <FaRulerVertical className="text-gray-600 w-5 h-5" />
          <span className="text-gray-700 font-medium">Height: {height} cm</span>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-gray-700 font-medium">BMI: </span>
          <span className={`font-bold ${status === "Normal" ? "text-green-600" : status === "Underweight" ? "text-blue-600" : "text-red-600"}`}>
            {bmi} ({status})
          </span>
        </div>
      </CardContent>
    </Card>
  );
}
