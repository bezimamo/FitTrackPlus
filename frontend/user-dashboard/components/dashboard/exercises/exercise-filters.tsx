"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Search, Filter, X } from "lucide-react"
import { type ExerciseFilters, type ExerciseMetadata, exerciseAPI } from "@/lib/types/exercises"

interface ExerciseFiltersProps {
  filters: ExerciseFilters
  onFiltersChange: (filters: ExerciseFilters) => void
  searchQuery: string
  onSearchChange: (query: string) => void
}

export function ExerciseFiltersComponent({
  filters,
  onFiltersChange,
  searchQuery,
  onSearchChange,
}: ExerciseFiltersProps) {
  const [metadata, setMetadata] = useState<ExerciseMetadata | null>(null)
  const [isExpanded, setIsExpanded] = useState(false)

  useEffect(() => {
    const loadMetadata = async () => {
      try {
        const data = await exerciseAPI.getExerciseMetadata()
        setMetadata(data)
      } catch (error) {
        console.error("Failed to load exercise metadata:", error)
      }
    }

    loadMetadata()
  }, [])

  const handleFilterChange = (key: keyof ExerciseFilters, value: string) => {
    onFiltersChange({
      ...filters,
      [key]: value === "all" ? undefined : value,
    })
  }

  const clearFilters = () => {
    onFiltersChange({})
    onSearchChange("")
  }

  const activeFiltersCount = Object.values(filters).filter(Boolean).length

  return (
    <Card className="mb-6">
      <CardHeader className="pb-4">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg font-semibold flex items-center gap-2">
            <Filter className="h-5 w-5" />
            Filters
            {activeFiltersCount > 0 && (
              <Badge variant="secondary" className="ml-2">
                {activeFiltersCount}
              </Badge>
            )}
          </CardTitle>
          <div className="flex gap-2">
            {(activeFiltersCount > 0 || searchQuery) && (
              <Button variant="outline" size="sm" onClick={clearFilters}>
                <X className="h-4 w-4 mr-1" />
                Clear
              </Button>
            )}
            <Button variant="outline" size="sm" onClick={() => setIsExpanded(!isExpanded)}>
              {isExpanded ? "Less" : "More"} Filters
            </Button>
          </div>
        </div>
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
          <Input
            placeholder="Search exercises..."
            value={searchQuery}
            onChange={(e) => onSearchChange(e.target.value)}
            className="pl-10"
          />
        </div>

        {/* Quick Filters */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div>
            <Label htmlFor="category" className="text-sm font-medium">
              Category
            </Label>
            <Select value={filters.category || "all"} onValueChange={(value) => handleFilterChange("category", value)}>
              <SelectTrigger>
                <SelectValue placeholder="All Categories" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Categories</SelectItem>
                {metadata?.categories.map((category) => (
                  <SelectItem key={category} value={category}>
                    {category.charAt(0).toUpperCase() + category.slice(1)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div>
            <Label htmlFor="difficulty" className="text-sm font-medium">
              Difficulty
            </Label>
            <Select
              value={filters.difficulty || "all"}
              onValueChange={(value) => handleFilterChange("difficulty", value)}
            >
              <SelectTrigger>
                <SelectValue placeholder="All Levels" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Levels</SelectItem>
                {metadata?.difficulties.map((difficulty) => (
                  <SelectItem key={difficulty} value={difficulty}>
                    {difficulty.charAt(0).toUpperCase() + difficulty.slice(1)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {isExpanded && (
            <>
              <div>
                <Label htmlFor="muscle_group" className="text-sm font-medium">
                  Muscle Group
                </Label>
                <Select
                  value={filters.muscle_group || "all"}
                  onValueChange={(value) => handleFilterChange("muscle_group", value)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="All Muscles" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Muscles</SelectItem>
                    {metadata?.muscle_groups.map((group) => (
                      <SelectItem key={group} value={group}>
                        {group.charAt(0).toUpperCase() + group.slice(1)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div>
                <Label htmlFor="equipment" className="text-sm font-medium">
                  Equipment
                </Label>
                <Select
                  value={filters.equipment || "all"}
                  onValueChange={(value) => handleFilterChange("equipment", value)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="All Equipment" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Equipment</SelectItem>
                    {metadata?.equipment.map((equip) => (
                      <SelectItem key={equip} value={equip}>
                        {equip.charAt(0).toUpperCase() + equip.slice(1)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </>
          )}
        </div>

        {/* Active Filters Display */}
        {activeFiltersCount > 0 && (
          <div className="flex flex-wrap gap-2 pt-2 border-t">
            <span className="text-sm text-gray-600 font-medium">Active filters:</span>
            {Object.entries(filters).map(([key, value]) =>
              value ? (
                <Badge
                  key={key}
                  variant="secondary"
                  className="cursor-pointer hover:bg-gray-200"
                  onClick={() => handleFilterChange(key as keyof ExerciseFilters, "all")}
                >
                  {key.replace("_", " ")}: {value}
                  <X className="h-3 w-3 ml-1" />
                </Badge>
              ) : null,
            )}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
