import { TrendingUp, TrendingDown, Target, Calendar, Activity } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

interface StatsOverviewProps {
  currentWeight: number // assumed in lbs
  weightChange: number  // assumed in lbs
  totalWorkouts: number
  totalDuration: number
}

export function StatsOverview({ currentWeight, weightChange, totalWorkouts, totalDuration }: StatsOverviewProps) {
  // Convert lbs to kg for Ethiopia standard (keep as number)
  const currentWeightKg = currentWeight * 0.453592
  const weightChangeKg = weightChange * 0.453592

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Current Weight</CardTitle>
          {weightChange < 0 ? (
            <TrendingDown className="h-4 w-4 text-primary" />
          ) : (
            <TrendingUp className="h-4 w-4 text-destructive" />
          )}
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{currentWeightKg.toFixed(1)} kg</div>
          <p className="text-xs text-muted-foreground">
            {weightChange < 0 ? "" : "+"}
            {weightChangeKg.toFixed(1)} kg from last week
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Total Workouts</CardTitle>
          <Activity className="h-4 w-4 text-primary" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{totalWorkouts}</div>
          <p className="text-xs text-muted-foreground">This month</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Total Duration</CardTitle>
          <Calendar className="h-4 w-4 text-primary" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {Math.floor(totalDuration / 60)}h {totalDuration % 60}m
          </div>
          <p className="text-xs text-muted-foreground">This month</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">BMI</CardTitle>
          <Target className="h-4 w-4 text-primary" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">
            {currentWeightKg ? (currentWeightKg / 2.89).toFixed(1) : 0}
          </div>
          <p className="text-xs text-muted-foreground">Normal range</p>
        </CardContent>
      </Card>
    </div>
  )
}
