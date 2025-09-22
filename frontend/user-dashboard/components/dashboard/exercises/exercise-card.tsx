"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Play, Info, Dumbbell } from "lucide-react"
import type { Exercise } from "@/lib/types/exercises"
import Image from "next/image"

interface ExerciseCardProps {
  exercise: Exercise
  onViewDetails: (exercise: Exercise) => void
  onAddToWorkout?: (exercise: Exercise) => void
  showAddButton?: boolean
}

export function ExerciseCard({ exercise, onViewDetails, onAddToWorkout, showAddButton = false }: ExerciseCardProps) {
  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case "beginner":
        return "bg-green-100 text-green-800 hover:bg-green-200"
      case "intermediate":
        return "bg-yellow-100 text-yellow-800 hover:bg-yellow-200"
      case "advanced":
        return "bg-red-100 text-red-800 hover:bg-red-200"
      default:
        return "bg-gray-100 text-gray-800 hover:bg-gray-200"
    }
  }

  const getCategoryColor = (category: string) => {
    switch (category) {
      case "strength":
        return "bg-blue-100 text-blue-800 hover:bg-blue-200"
      case "cardio":
        return "bg-orange-100 text-orange-800 hover:bg-orange-200"
      case "flexibility":
        return "bg-purple-100 text-purple-800 hover:bg-purple-200"
      case "balance":
        return "bg-teal-100 text-teal-800 hover:bg-teal-200"
      default:
        return "bg-gray-100 text-gray-800 hover:bg-gray-200"
    }
  }

  return (
    <Card className="group hover:shadow-lg transition-all duration-200 border-0 bg-white">
      <CardHeader className="pb-3">
        <div className="relative aspect-video w-full mb-3 rounded-lg overflow-hidden bg-gray-100">
          {exercise.image_url ? (
            <Image
              src={exercise.image_url || "/placeholder.svg"}
              alt={exercise.name}
              fill
              className="object-cover group-hover:scale-105 transition-transform duration-200"
            />
          ) : (
            <div className="flex items-center justify-center h-full">
              <Dumbbell className="h-12 w-12 text-gray-400" />
            </div>
          )}
          {exercise.video_url && (
            <div className="absolute top-2 right-2">
              <div className="bg-black/70 rounded-full p-1.5">
                <Play className="h-4 w-4 text-white" />
              </div>
            </div>
          )}
        </div>

        <div className="space-y-2">
          <CardTitle className="text-lg font-semibold text-gray-900 line-clamp-2">{exercise.name}</CardTitle>

          <div className="flex flex-wrap gap-2">
            <Badge className={getCategoryColor(exercise.category)}>{exercise.category}</Badge>
            <Badge className={getDifficultyColor(exercise.difficulty)}>{exercise.difficulty}</Badge>
          </div>
        </div>
      </CardHeader>

      <CardContent className="pt-0">
        <div className="space-y-3">
          <div className="text-sm text-gray-600">
            <p className="font-medium">Target: {exercise.muscle_group}</p>
            <p>Equipment: {exercise.equipment}</p>
          </div>

          {exercise.description && <p className="text-sm text-gray-700 line-clamp-2">{exercise.description}</p>}

          <div className="flex gap-2 pt-2">
            <Button variant="outline" size="sm" onClick={() => onViewDetails(exercise)} className="flex-1">
              <Info className="h-4 w-4 mr-1" />
              Details
            </Button>

            {showAddButton && onAddToWorkout && (
              <Button size="sm" onClick={() => onAddToWorkout(exercise)} className="flex-1">
                Add to Workout
              </Button>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
