"use client"

import { useEffect, useMemo, useRef, useState } from "react"
import { Card, CardContent } from "@/components/ui/card"
import { Star, ChevronLeft, ChevronRight } from "lucide-react"

const testimonials = [
  {
    name: "Saron Bekele",
    role: "Member",
    content:
      "FitTrack+ made it simple to book sessions and actually stay consistent. I'm finally seeing steady progress!",
    rating: 5,
    image: "/ethiopia-saron-bekele.jpg",
  },
  {
    name: "Abel Tadesse",
    role: "Personal Trainer",
    content:
      "Managing clients is so much easier now. The tools for plans and tracking help me deliver better results.",
    rating: 5,
    image: "/ethiopia-abel-tadesse.jpg",
  },
  {
    name: "Martha Abate",
    role: "Gym Member",
    content:
      "I love the personalized plans and the clear progress charts. It's motivating and easy to follow.",
    rating: 5,
    image: "/ethiopia-martha-abate.jpg",
  },
]

export function TestimonialsSection() {
  const [index, setIndex] = useState(0)
  const [perView, setPerView] = useState(1)
  const timerRef = useRef<NodeJS.Timeout | null>(null)

  // Responsive items per view
  useEffect(() => {
    const update = () => {
      if (window.matchMedia("(min-width: 1024px)").matches) setPerView(3)
      else if (window.matchMedia("(min-width: 768px)").matches) setPerView(2)
      else setPerView(1)
    }
    update()
    window.addEventListener("resize", update)
    return () => window.removeEventListener("resize", update)
  }, [])

  const maxIndex = useMemo(() => Math.max(0, testimonials.length - perView), [perView])
  const next = () => setIndex((i) => (i >= maxIndex ? 0 : i + 1))
  const prev = () => setIndex((i) => (i <= 0 ? maxIndex : i - 1))

  useEffect(() => {
    timerRef.current && clearInterval(timerRef.current)
    timerRef.current = setInterval(() => {
      setIndex((i) => (i >= maxIndex ? 0 : i + 1))
    }, 6000)
    return () => {
      if (timerRef.current) clearInterval(timerRef.current)
    }
  }, [maxIndex])

  const slideWidthPercent = useMemo(() => 100 / perView, [perView])
  const trackTranslate = useMemo(() => index * slideWidthPercent, [index, slideWidthPercent])

  const getInitials = (name: string) =>
    name
      .split(" ")
      .map((n) => n[0])
      .slice(0, 2)
      .join("")

  return (
    <section id="testimonials" className="py-24 bg-background">
      <div className="max-w-6xl mx-auto px-6">
        <div className="text-center mb-12">
          <h2 className="text-4xl md:text-5xl font-bold text-foreground mb-4 text-balance">
            Voices from
            <span className="text-primary block">Our Community</span>
          </h2>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto text-pretty">
            Real stories from members and trainers using FitTrack+ every day.
          </p>
        </div>

        <div className="relative">
          <div className="overflow-hidden">
            <div
              className="flex transition-transform duration-500"
              style={{ transform: `translateX(-${trackTranslate}%)` }}
            >
              {testimonials.map((t, i) => (
                <div key={i} className="px-2" style={{ minWidth: `${slideWidthPercent}%` }}>
                  <Card className="relative overflow-hidden border border-primary/10 bg-background/70 backdrop-blur-sm rounded-xl shadow-sm">
                    <div className="pointer-events-none absolute -right-8 -top-8 h-24 w-24 rounded-full bg-primary/10 blur-2xl" />
                    <CardContent className="p-6 md:p-8">
                      <div className="flex items-center mb-3">
                        {[...Array(t.rating)].map((_, ri) => (
                          <Star key={ri} className="h-5 w-5 fill-primary text-primary" />
                        ))}
                      </div>
                      <blockquote className="text-muted-foreground mb-5 leading-relaxed">"{t.content}"</blockquote>
                      <div className="flex items-center gap-4">
                        {t.image ? (
                          <img
                            src={t.image}
                            alt={t.name}
                            onError={(e) => {
                              ;(e.currentTarget as HTMLImageElement).style.display = "none"
                              const sibling = (e.currentTarget.nextSibling as HTMLElement)!
                              if (sibling) sibling.style.display = "grid"
                            }}
                            className="w-12 h-12 rounded-full object-cover"
                          />
                        ) : null}
                        <div
                          style={{ display: t.image ? "none" : "grid" }}
                          className="w-12 h-12 rounded-full bg-primary/10 text-primary grid place-items-center font-semibold"
                        >
                          {getInitials(t.name)}
                        </div>
                        <div>
                          <div className="font-semibold text-foreground">{t.name}</div>
                          <div className="text-sm text-muted-foreground">{t.role}</div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </div>
              ))}
            </div>
          </div>

          <button
            aria-label="Previous"
            onClick={prev}
            className="absolute left-2 top-1/2 -translate-y-1/2 h-10 w-10 rounded-full bg-primary text-primary-foreground grid place-items-center shadow hover:bg-primary/90"
          >
            <ChevronLeft className="h-5 w-5" />
          </button>
          <button
            aria-label="Next"
            onClick={next}
            className="absolute right-2 top-1/2 -translate-y-1/2 h-10 w-10 rounded-full bg-primary text-primary-foreground grid place-items-center shadow hover:bg-primary/90"
          >
            <ChevronRight className="h-5 w-5" />
          </button>

          <div className="mt-6 flex justify-center gap-2">
            {Array.from({ length: maxIndex + 1 }).map((_, i) => (
              <button
                key={i}
                onClick={() => setIndex(i)}
                aria-label={`Go to slide ${i + 1}`}
                className={`h-2.5 w-2.5 rounded-full transition-colors ${
                  i === index ? "bg-primary" : "bg-muted"
                }`}
              />
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}
