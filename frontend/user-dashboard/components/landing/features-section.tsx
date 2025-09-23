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
    <section id="features" className="py-24 bg-muted/20">
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
            <div
              key={index}
              className="group relative overflow-hidden rounded-2xl border border-primary/15 bg-background/80 backdrop-blur-sm p-6 md:p-8 shadow-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:border-primary/25"
            >
              <div className="pointer-events-none absolute -right-8 -top-8 h-24 w-24 rounded-full bg-primary/10 blur-2xl" />
              <div className="w-14 h-14 rounded-full bg-primary/10 ring-1 ring-primary/25 grid place-items-center mb-4 transition-all duration-300 group-hover:bg-primary/15 group-hover:ring-primary/40">
                <feature.icon className="h-7 w-7 text-primary transition-transform duration-300 group-hover:scale-110" />
              </div>
              <h3 className="text-xl font-semibold text-foreground tracking-tight">{feature.title}</h3>
              <p className="mt-2 text-base leading-relaxed text-muted-foreground">{feature.description}</p>
              <div className="mt-6 h-px w-full bg-gradient-to-r from-transparent via-primary/20 to-transparent" />
              <button className="mt-4 inline-flex items-center gap-2 text-sm text-primary hover:underline">
                Learn more
                <svg className="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M5 12h14"/><path d="M12 5l7 7-7 7"/></svg>
              </button>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
