"use client"

import { useState, useEffect, useMemo } from "react"
import { Card, CardContent, CardHeader } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Search, Grid, List, AlertCircle } from "lucide-react"
import { type Exercise, type ExerciseFilters, exerciseAPI } from "@/lib/types/exercises"
import { ExerciseCard } from "./exercise-card"
import { ExerciseFiltersComponent } from "./exercise-filters"
import { ExerciseDetailModal } from "./exercise-detail-modal"

interface ExerciseLibraryProps {
  onAddToWorkout?: (exercise: Exercise) => void
  showAddButton?: boolean
}

export function ExerciseLibrary({ onAddToWorkout, showAddButton = false }: ExerciseLibraryProps) {
  const [exercises, setExercises] = useState<Exercise[]>([]) // ✅ always array
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [filters, setFilters] = useState<ExerciseFilters>({})
  const [searchQuery, setSearchQuery] = useState("")
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid")
  const [selectedExercise, setSelectedExercise] = useState<Exercise | null>(null)

  // Load exercises
  useEffect(() => {
    const loadExercises = async () => {
      try {
        setLoading(true)
        setError(null)
        const response = await exerciseAPI.getExercises({ ...filters, status: "active" })
        setExercises(response?.exercises || []) // ✅ fallback to empty array
      } catch (err) {
        setError("Failed to load exercises. Please try again.")
        console.error("Error loading exercises:", err)
        setExercises([]) // ✅ avoid null
      } finally {
        setLoading(false)
      }
    }

    loadExercises()
  }, [filters])

  // Filter exercises by search query
  const filteredExercises = useMemo(() => {
    if (!exercises) return [] // ✅ safeguard
    if (!searchQuery) return exercises

    const query = searchQuery.toLowerCase()
    return exercises.filter(
      (exercise) =>
        exercise.name?.toLowerCase().includes(query) ||
        exercise.description?.toLowerCase().includes(query) ||
        exercise.muscle_group?.toLowerCase().includes(query) ||
        exercise.equipment?.toLowerCase().includes(query) ||
        exercise.category?.toLowerCase().includes(query),
    )
  }, [exercises, searchQuery])

  const handleViewDetails = (exercise: Exercise) => {
    setSelectedExercise(exercise)
  }

  const handleCloseModal = () => {
    setSelectedExercise(null)
  }

  if (error) {
    return (
      <Alert className="mb-6">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>{error}</AlertDescription>
      </Alert>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Exercise Library</h1>
          <p className="text-gray-600 mt-1">Discover and learn new exercises for your fitness journey</p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant={viewMode === "grid" ? "default" : "outline"}
            size="sm"
            onClick={() => setViewMode("grid")}
          >
            <Grid className="h-4 w-4" />
          </Button>
          <Button
            variant={viewMode === "list" ? "default" : "outline"}
            size="sm"
            onClick={() => setViewMode("list")}
          >
            <List className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* Filters */}
      <ExerciseFiltersComponent
        filters={filters}
        onFiltersChange={setFilters}
        searchQuery={searchQuery}
        onSearchChange={setSearchQuery}
      />

      {/* Results Summary */}
      <div className="flex items-center justify-between">
        <p className="text-sm text-gray-600">
          {loading ? "Loading..." : `${filteredExercises.length} exercises found`}
        </p>
      </div>

      {/* Exercise Grid/List */}
      {loading ? (
        <div
          className={`grid gap-6 ${
            viewMode === "grid" ? "grid-cols-1 md:grid-cols-2 lg:grid-cols-3" : "grid-cols-1"
          }`}
        >
          {Array.from({ length: 6 }).map((_, i) => (
            <Card key={i}>
              <CardHeader>
                <Skeleton className="h-48 w-full rounded-lg" />
                <Skeleton className="h-6 w-3/4" />
                <div className="flex gap-2">
                  <Skeleton className="h-6 w-16" />
                  <Skeleton className="h-6 w-20" />
                </div>
              </CardHeader>
              <CardContent>
                <Skeleton className="h-4 w-full mb-2" />
                <Skeleton className="h-4 w-2/3" />
              </CardContent>
            </Card>
          ))}
        </div>
      ) : filteredExercises.length === 0 ? (
        <Card className="text-center py-12">
          <CardContent>
            <Search className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-semibold text-gray-900 mb-2">No exercises found</h3>
            <p className="text-gray-600 mb-4">Try adjusting your filters or search terms to find exercises.</p>
            <Button
              variant="outline"
              onClick={() => {
                setFilters({})
                setSearchQuery("")
              }}
            >
              Clear all filters
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div
          className={`grid gap-6 ${
            viewMode === "grid" ? "grid-cols-1 md:grid-cols-2 lg:grid-cols-3" : "grid-cols-1"
          }`}
        >
          {filteredExercises.map((exercise) => (
            <ExerciseCard
              key={exercise.id}
              exercise={exercise}
              onViewDetails={handleViewDetails}
              onAddToWorkout={onAddToWorkout}
              showAddButton={showAddButton}
            />
          ))}
        </div>
      )}

      {/* Exercise Detail Modal */}
      {selectedExercise && (
        <ExerciseDetailModal
          exercise={selectedExercise}
          isOpen={!!selectedExercise}
          onClose={handleCloseModal}
          onAddToWorkout={onAddToWorkout}
          showAddButton={showAddButton}
        />
      )}
    </div>
  )
}
