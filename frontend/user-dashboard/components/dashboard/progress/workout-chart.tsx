import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from "recharts"

interface WorkoutData {
  week: string
  sessions: number
  duration: number
}

interface WorkoutChartProps {
  data: WorkoutData[]
}

export function WorkoutChart({ data }: WorkoutChartProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Workout Activity</CardTitle>
        <CardDescription>Weekly workout sessions and duration</CardDescription>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={data}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="week" />
            <YAxis />
            <Tooltip />
            <Bar dataKey="sessions" fill="hsl(var(--primary))" name="Sessions" />
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
}
