import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Progress } from "@/components/ui/progress"

interface Goal {
  id: number
  title: string
  target: number
  current: number
  unit: string
  deadline: string
  category: string
}

interface GoalsSectionProps {
  goals: Goal[]
}

export function GoalsSection({ goals }: GoalsSectionProps) {
  // Convert weight goals from lbs to kg
  const convertToKg = (value: number, unit: string) => {
    if (unit.toLowerCase() === "lbs") {
      return value * 0.453592
    }
    return value
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Current Goals</CardTitle>
        <CardDescription>Track your fitness objectives</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {goals.map((goal) => {
          const current = convertToKg(goal.current, goal.unit)
          const target = convertToKg(goal.target, goal.unit)
          const unit = goal.unit.toLowerCase() === "lbs" ? "kg" : goal.unit

          return (
            <div key={goal.id} className="space-y-2">
              <div className="flex justify-between items-center">
                <h4 className="font-medium">{goal.title}</h4>
                <Badge variant="outline">{goal.category}</Badge>
              </div>
              <Progress value={(current / target) * 100} className="h-2" />
              <div className="flex justify-between text-sm text-muted-foreground">
                <span>
                  {current.toFixed(1)} / {target.toFixed(1)} {unit}
                </span>
                <span>Due: {new Date(goal.deadline).toLocaleDateString()}</span>
              </div>
            </div>
          )
        })}
      </CardContent>
    </Card>
  )
}
