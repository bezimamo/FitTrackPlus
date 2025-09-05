import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

interface Workout {
  id: number
  date: string
  type: string
  duration: number
  calories: number
  exercises: number
}

interface RecentWorkoutsProps {
  workouts: Workout[]
}

export function RecentWorkouts({ workouts }: RecentWorkoutsProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Recent Workouts</CardTitle>
        <CardDescription>Your latest training sessions</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {workouts.map((workout) => (
          <div key={workout.id} className="flex items-center justify-between p-3 border rounded-lg">
            <div>
              <h4 className="font-medium">{workout.type}</h4>
              <p className="text-sm text-muted-foreground">{new Date(workout.date).toLocaleDateString()}</p>
            </div>
            <div className="text-right">
              <p className="font-medium">{workout.duration} min</p>
              <p className="text-sm text-muted-foreground">{workout.calories} cal</p>
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
  )
}
