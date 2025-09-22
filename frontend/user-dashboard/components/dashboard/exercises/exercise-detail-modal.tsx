"use client"

import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Play, Plus, Target, Zap, Wrench, Info } from "lucide-react"
import type { Exercise } from "@/lib/types/exercises"
import Image from "next/image"

interface ExerciseDetailModalProps {
  exercise: Exercise
  isOpen: boolean
  onClose: () => void
  onAddToWorkout?: (exercise: Exercise) => void
  showAddButton?: boolean
}

export function ExerciseDetailModal({
  exercise,
  isOpen,
  onClose,
  onAddToWorkout,
  showAddButton = false,
}: ExerciseDetailModalProps) {
  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case "beginner":
        return "bg-green-100 text-green-800"
      case "intermediate":
        return "bg-yellow-100 text-yellow-800"
      case "advanced":
        return "bg-red-100 text-red-800"
      default:
        return "bg-gray-100 text-gray-800"
    }
  }

  const getCategoryColor = (category: string) => {
    switch (category) {
      case "strength":
        return "bg-blue-100 text-blue-800"
      case "cardio":
        return "bg-orange-100 text-orange-800"
      case "flexibility":
        return "bg-purple-100 text-purple-800"
      case "balance":
        return "bg-teal-100 text-teal-800"
      default:
        return "bg-gray-100 text-gray-800"
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-4xl max-h-[90vh]">
        <DialogHeader>
          <DialogTitle className="text-2xl font-bold text-gray-900">{exercise.name}</DialogTitle>
        </DialogHeader>

        <ScrollArea className="max-h-[calc(90vh-120px)]">
          <div className="space-y-6">
            {/* Media Section */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {/* Image */}
              {exercise.image_url && (
                <div className="relative aspect-video rounded-lg overflow-hidden bg-gray-100">
                  <Image
                    src={exercise.image_url || "/placeholder.svg"}
                    alt={exercise.name}
                    fill
                    className="object-cover"
                  />
                </div>
              )}

              {/* Video */}
              {exercise.video_url && (
                <div className="relative aspect-video rounded-lg overflow-hidden bg-gray-100">
                  <video
                    src={exercise.video_url}
                    controls
                    className="w-full h-full object-cover"
                    poster={exercise.image_url}
                  >
                    Your browser does not support the video tag.
                  </video>
                  <div className="absolute top-2 left-2">
                    <Badge className="bg-black/70 text-white">
                      <Play className="h-3 w-3 mr-1" />
                      Video
                    </Badge>
                  </div>
                </div>
              )}
            </div>

            {/* Exercise Info */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-4">
                <div>
                  <h3 className="text-lg font-semibold mb-3 flex items-center gap-2">
                    <Info className="h-5 w-5" />
                    Exercise Details
                  </h3>
                  <div className="space-y-3">
                    <div className="flex items-center gap-2">
                      <Badge className={getCategoryColor(exercise.category)}>{exercise.category}</Badge>
                      <Badge className={getDifficultyColor(exercise.difficulty)}>{exercise.difficulty}</Badge>
                    </div>

                    <div className="space-y-2 text-sm">
                      <div className="flex items-center gap-2">
                        <Target className="h-4 w-4 text-gray-500" />
                        <span className="font-medium">Target:</span>
                        <span className="capitalize">{exercise.muscle_group}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <Wrench className="h-4 w-4 text-gray-500" />
                        <span className="font-medium">Equipment:</span>
                        <span className="capitalize">{exercise.equipment}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <Zap className="h-4 w-4 text-gray-500" />
                        <span className="font-medium">Difficulty:</span>
                        <span className="capitalize">{exercise.difficulty}</span>
                      </div>
                    </div>
                  </div>
                </div>

                {showAddButton && onAddToWorkout && (
                  <Button onClick={() => onAddToWorkout(exercise)} className="w-full">
                    <Plus className="h-4 w-4 mr-2" />
                    Add to Workout
                  </Button>
                )}
              </div>

              <div className="space-y-4">
                {exercise.description && (
                  <div>
                    <h4 className="font-semibold text-gray-900 mb-2">Description</h4>
                    <p className="text-gray-700 text-sm leading-relaxed">{exercise.description}</p>
                  </div>
                )}
              </div>
            </div>

            <Separator />

            {/* Instructions */}
            {exercise.instructions && (
              <div>
                <h3 className="text-lg font-semibold mb-3">Instructions</h3>
                <div className="bg-gray-50 rounded-lg p-4">
                  <p className="text-gray-700 text-sm leading-relaxed whitespace-pre-line">{exercise.instructions}</p>
                </div>
              </div>
            )}

            {/* Tips */}
            {exercise.tips && (
              <div>
                <h3 className="text-lg font-semibold mb-3">Tips & Form Cues</h3>
                <div className="bg-blue-50 rounded-lg p-4 border-l-4 border-blue-400">
                  <p className="text-blue-800 text-sm leading-relaxed whitespace-pre-line">{exercise.tips}</p>
                </div>
              </div>
            )}
          </div>
        </ScrollArea>
      </DialogContent>
    </Dialog>
  )
}
