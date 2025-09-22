const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080"

export interface Exercise {
  id: number
  name: string
  description: string
  category: "strength" | "cardio" | "flexibility" | "balance"
  muscle_group: string
  difficulty: "beginner" | "intermediate" | "advanced"
  equipment: string
  video_url: string
  image_url: string
  instructions: string
  tips: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface ExerciseFilters {
  category?: string
  muscle_group?: string
  difficulty?: string
  equipment?: string
  status?: "active" | "inactive"
}

export interface ExerciseResponse {
  message: string
  exercises: Exercise[]
  total: number
  filters: ExerciseFilters
}

export interface ExerciseDetailResponse {
  message: string
  exercise: Exercise
}

export interface ExerciseMetadata {
  categories: string[]
  muscle_groups: string[]
  difficulties: string[]
  equipment: string[]
}

class ExerciseAPI {
  private async fetchWithAuth(url: string, options: RequestInit = {}) {
    const token = localStorage.getItem("auth_token")

    return fetch(`${API_BASE_URL}${url}`, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...(token && { Authorization: `Bearer ${token}` }),
        ...options.headers,
      },
    })
  }

  async getExercises(filters: ExerciseFilters = {}): Promise<ExerciseResponse> {
    const params = new URLSearchParams()

    Object.entries(filters).forEach(([key, value]) => {
      if (value) params.append(key, value)
    })

    const response = await this.fetchWithAuth(`/exercises?${params.toString()}`)

    if (!response.ok) {
      throw new Error("Failed to fetch exercises")
    }

    return response.json()
  }

  async getExercise(id: number): Promise<ExerciseDetailResponse> {
    const response = await this.fetchWithAuth(`/exercises/${id}`)

    if (!response.ok) {
      throw new Error("Failed to fetch exercise details")
    }

    return response.json()
  }

  async getExerciseCategories(): Promise<{ message: string; categories: string[] }> {
    const response = await this.fetchWithAuth("/exercises/categories")

    if (!response.ok) {
      throw new Error("Failed to fetch exercise categories")
    }

    return response.json()
  }

  async getMuscleGroups(): Promise<{ message: string; muscle_groups: string[] }> {
    const response = await this.fetchWithAuth("/exercises/muscle-groups")

    if (!response.ok) {
      throw new Error("Failed to fetch muscle groups")
    }

    return response.json()
  }

  async getDifficulties(): Promise<{ message: string; difficulties: string[] }> {
    const response = await this.fetchWithAuth("/exercises/difficulties")

    if (!response.ok) {
      throw new Error("Failed to fetch difficulties")
    }

    return response.json()
  }

  async getEquipment(): Promise<{ message: string; equipment: string[] }> {
    const response = await this.fetchWithAuth("/exercises/equipment")

    if (!response.ok) {
      throw new Error("Failed to fetch equipment types")
    }

    return response.json()
  }

  async getExerciseMetadata(): Promise<ExerciseMetadata> {
    const [categories, muscleGroups, difficulties, equipment] = await Promise.all([
      this.getExerciseCategories(),
      this.getMuscleGroups(),
      this.getDifficulties(),
      this.getEquipment(),
    ])

    return {
      categories: categories.categories,
      muscle_groups: muscleGroups.muscle_groups,
      difficulties: difficulties.difficulties,
      equipment: equipment.equipment,
    }
  }
}

export const exerciseAPI = new ExerciseAPI()
