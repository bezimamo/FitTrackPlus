import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Calendar, TrendingUp, Users, Target, Clock, Award } from "lucide-react"

const features = [
  {
    icon: Calendar,
    title: "Easy Booking",
    description: "Schedule personal training, group classes, and physiotherapy sessions with just a few clicks.",
  },
  {
    icon: TrendingUp,
    title: "Progress Tracking",
    description: "Monitor your fitness journey with detailed analytics, charts, and achievement tracking.",
  },
  {
    icon: Target,
    title: "Personalized Plans",
    description: "Get custom fitness and nutrition plans tailored to your specific goals and needs.",
  },
  {
    icon: Users,
    title: "Expert Trainers",
    description: "Work with certified professionals who guide you every step of your fitness journey.",
  },
  {
    icon: Clock,
    title: "Flexible Scheduling",
    description: "Book sessions that fit your schedule with our 24/7 online booking system.",
  },
  {
    icon: Award,
    title: "Goal Achievement",
    description: "Set, track, and celebrate your fitness milestones with our comprehensive goal system.",
  },
]

export function FeaturesSection() {
  return (
    <section className="py-24 bg-muted/30">
      <div className="max-w-7xl mx-auto px-6">
        <div className="text-center mb-16">
          <h2 className="text-4xl md:text-5xl font-bold text-foreground mb-6 text-balance">
            Everything You Need to
            <span className="text-primary block">Succeed</span>
          </h2>
          <p className="text-xl text-muted-foreground max-w-3xl mx-auto text-pretty">
            FitTrack+ provides all the tools and support you need to achieve your fitness goals, whether you're just
            starting out or looking to take your training to the next level.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {features.map((feature, index) => (
            <Card
              key={index}
              className="group hover:shadow-lg transition-all duration-300 border-border/50 hover:border-primary/20"
            >
              <CardHeader className="text-center pb-4">
                <div className="mx-auto w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mb-4 group-hover:bg-primary/20 transition-colors">
                  <feature.icon className="h-8 w-8 text-primary" />
                </div>
                <CardTitle className="text-xl font-semibold">{feature.title}</CardTitle>
              </CardHeader>
              <CardContent className="text-center">
                <CardDescription className="text-base leading-relaxed">{feature.description}</CardDescription>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  )
}
