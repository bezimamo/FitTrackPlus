"use client"

import { useParams, useRouter } from "next/navigation"
import { useState, useEffect } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { ArrowLeft, Play, Target, Zap, Wrench, Info, Plus, AlertCircle } from "lucide-react"
import { type Exercise, exerciseAPI } from "@/lib/types/exercises"
import Image from "next/image"

export default function ExerciseDetailPage() {
  const params = useParams()
  const router = useRouter()
  const exerciseId = Number.parseInt(params.id as string)

  const [exercise, setExercise] = useState<Exercise | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const loadExercise = async () => {
      try {
        setLoading(true)
        setError(null)
        const response = await exerciseAPI.getExercise(exerciseId)
        setExercise(response.exercise)
      } catch (err) {
        setError("Failed to load exercise details. Please try again.")
        console.error("Error loading exercise:", err)
      } finally {
        setLoading(false)
      }
    }

    if (exerciseId && !isNaN(exerciseId)) {
      loadExercise()
    } else {
      setError("Invalid exercise ID")
      setLoading(false)
    }
  }, [exerciseId])

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

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-6 space-y-6">
        <div className="flex items-center gap-4">
          <Skeleton className="h-10 w-10 rounded-lg" />
          <Skeleton className="h-8 w-48" />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Skeleton className="aspect-video rounded-lg" />
          <div className="space-y-4">
            <Skeleton className="h-6 w-3/4" />
            <div className="flex gap-2">
              <Skeleton className="h-6 w-16" />
              <Skeleton className="h-6 w-20" />
            </div>
            <Skeleton className="h-20 w-full" />
          </div>
        </div>
      </div>
    )
  }

  if (error || !exercise) {
    return (
      <div className="container mx-auto px-4 py-6">
        <Alert>
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>{error || "Exercise not found."}</AlertDescription>
        </Alert>
        <Button variant="outline" onClick={() => router.push("/exercises")} className="mt-4">
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Exercises
        </Button>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-6 space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button variant="outline" size="sm" onClick={() => router.push("/exercises")}>
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back
        </Button>
        <div className="flex-1">
          <h1 className="text-3xl font-bold text-gray-900">{exercise.name}</h1>
          <div className="flex flex-wrap gap-2 mt-2">
            <Badge className={getCategoryColor(exercise.category)}>{exercise.category}</Badge>
            <Badge className={getDifficultyColor(exercise.difficulty)}>{exercise.difficulty}</Badge>
          </div>
        </div>
        <Button>
          <Plus className="h-4 w-4 mr-2" />
          Add to Workout
        </Button>
      </div>

      {/* Main Content */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Media Section */}
        <div className="space-y-4">
          {/* Image */}
          {exercise.image_url && (
            <Card>
              <CardContent className="p-0">
                <div className="relative aspect-video rounded-lg overflow-hidden bg-gray-100">
                  <Image
                    src={exercise.image_url || "/placeholder.svg"}
                    alt={exercise.name}
                    fill
                    className="object-cover"
                  />
                </div>
              </CardContent>
            </Card>
          )}

          {/* Video */}
          {exercise.video_url && (
            <Card>
              <CardHeader>
                <CardTitle className="text-lg flex items-center gap-2">
                  <Play className="h-5 w-5" />
                  Exercise Demonstration
                </CardTitle>
              </CardHeader>
              <CardContent className="p-0 pb-6">
                <div className="relative aspect-video rounded-lg overflow-hidden bg-gray-100 mx-6">
                  <video
                    src={exercise.video_url}
                    controls
                    className="w-full h-full object-cover"
                    poster={exercise.image_url}
                  >
                    Your browser does not support the video tag.
                  </video>
                </div>
              </CardContent>
            </Card>
          )}
        </div>

        {/* Exercise Information */}
        <div className="space-y-6">
          {/* Basic Info */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg flex items-center gap-2">
                <Info className="h-5 w-5" />
                Exercise Details
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {exercise.description && (
                <div>
                  <h4 className="font-semibold text-gray-900 mb-2">Description</h4>
                  <p className="text-gray-700 leading-relaxed">{exercise.description}</p>
                </div>
              )}

              <div className="grid grid-cols-1 gap-4">
                <div className="flex items-center gap-3">
                  <Target className="h-5 w-5 text-gray-500" />
                  <div>
                    <span className="font-medium text-gray-900">Target Muscle:</span>
                    <span className="ml-2 capitalize text-gray-700">{exercise.muscle_group}</span>
                  </div>
                </div>

                <div className="flex items-center gap-3">
                  <Wrench className="h-5 w-5 text-gray-500" />
                  <div>
                    <span className="font-medium text-gray-900">Equipment:</span>
                    <span className="ml-2 capitalize text-gray-700">{exercise.equipment}</span>
                  </div>
                </div>

                <div className="flex items-center gap-3">
                  <Zap className="h-5 w-5 text-gray-500" />
                  <div>
                    <span className="font-medium text-gray-900">Difficulty:</span>
                    <span className="ml-2 capitalize text-gray-700">{exercise.difficulty}</span>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Instructions */}
          {exercise.instructions && (
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">How to Perform</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="bg-gray-50 rounded-lg p-4">
                  <p className="text-gray-700 leading-relaxed whitespace-pre-line">{exercise.instructions}</p>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Tips */}
          {exercise.tips && (
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Tips & Form Cues</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="bg-blue-50 rounded-lg p-4 border-l-4 border-blue-400">
                  <p className="text-blue-800 leading-relaxed whitespace-pre-line">{exercise.tips}</p>
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </div>

      {/* Related Exercises Section */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Related Exercises</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-gray-600 text-center py-8">Related exercises feature coming soon...</p>
        </CardContent>
      </Card>
    </div>
  )
}
