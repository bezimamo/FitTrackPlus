"use client"

import { Card, CardContent } from "@/components/ui/card"
import { Star } from "lucide-react"

const testimonials = [
  {
    name: "Sarah Johnson",
    role: "Fitness Enthusiast",
    content:
      "FitTrack+ has completely transformed my fitness routine. The progress tracking keeps me motivated, and booking sessions is so easy!",
    rating: 5,
    image: "/professional-woman-smiling-headshot.png",
  },
  {
    name: "Mike Chen",
    role: "Personal Trainer",
    content:
      "As a trainer, FitTrack+ helps me manage my clients efficiently. The plan creation tools are fantastic and my clients love the progress tracking.",
    rating: 5,
    image: "/professional-headshot-man-fitness-trainer.jpg",
  },
  {
    name: "Emma Davis",
    role: "Gym Member",
    content:
      "The personalized plans and goal tracking have helped me achieve results I never thought possible. Highly recommend FitTrack+!",
    rating: 5,
    image: "/professional-headshot-woman-athlete.png",
  },
]

export function TestimonialsSection() {
  return (
    <section className="py-24 bg-background">
      <div className="max-w-7xl mx-auto px-6">
        <div className="text-center mb-16">
          <h2 className="text-4xl md:text-5xl font-bold text-foreground mb-6 text-balance">
            What Our Community
            <span className="text-primary block">Says</span>
          </h2>
          <p className="text-xl text-muted-foreground max-w-3xl mx-auto text-pretty">
            Join thousands of satisfied members who have transformed their fitness journey with FitTrack+.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {testimonials.map((testimonial, index) => (
            <Card key={index} className="hover:shadow-lg transition-all duration-300">
              <CardContent className="p-8">
                <div className="flex items-center mb-4">
                  {[...Array(testimonial.rating)].map((_, i) => (
                    <Star key={i} className="h-5 w-5 fill-primary text-primary" />
                  ))}
                </div>

                <blockquote className="text-muted-foreground mb-6 leading-relaxed">"{testimonial.content}"</blockquote>

                <div className="flex items-center gap-4">
                  <img
                    src={testimonial.image || "/placeholder.svg"}
                    alt={testimonial.name}
                    className="w-12 h-12 rounded-full object-cover"
                  />
                  <div>
                    <div className="font-semibold text-foreground">{testimonial.name}</div>
                    <div className="text-sm text-muted-foreground">{testimonial.role}</div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  )
}
