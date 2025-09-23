"use client"

import { Button } from "@/components/ui/button"
import { ArrowRight } from "lucide-react"
import Link from "next/link"

export function HeroSection() {
  return (
    <section id="home" className="relative w-screen h-screen flex flex-col">
      {/* Background Image */}
      <div className="absolute inset-0 -z-10">
        <img
          src="/diverse-people-working-out-in-modern-gym-with-equi.png"
          alt="People working out in modern gym"
          className="w-full h-full object-cover object-center"
        />
        {/* Overlay */}
        <div className="absolute inset-0 bg-black/70 bg-gradient-to-r from-primary/40 via-black/60 to-secondary/40" />
      </div>

      {/* Content - keep hero pinned to viewport; navbar floats above */}
      <div className="flex flex-1 items-center justify-center text-center text-white px-6">
        <div className="max-w-4xl mx-auto">
          <h1 className="text-4xl md:text-6xl font-bold mb-6 tracking-tight">
            Transform Your
            <span className="text-primary block">Fitness Journey</span>
          </h1>

          <p className="text-lg md:text-xl text-gray-200/90 mb-8 max-w-3xl mx-auto">
            Track progress, book sessions, and achieve your goals with FitTrack+ —
            the complete fitness management platform for members, trainers, and staff.
          </p>

          {/* Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center mb-12">
            <Button asChild size="lg" className="text-lg px-8 py-6 h-auto rounded-full bg-primary text-primary-foreground hover:bg-primary/90 shadow-lg shadow-primary/20">
              <Link href="/auth/register" className="flex items-center gap-2">
                Get Started Free
                <ArrowRight className="h-5 w-5" />
              </Link>
            </Button>

            <Button
              asChild
              variant="outline"
              size="lg"
              className="text-lg px-8 py-6 h-auto rounded-full border-white/40 bg-white/10 text-white hover:bg-white/20 backdrop-blur"
            >
              <Link href="/auth/login" className="flex items-center gap-2">
                Sign In
              </Link>
            </Button>
          </div>

          {/* Stats */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-2xl mx-auto">
            <div className="text-center">
              <div className="text-3xl font-bold text-primary mb-2">500+</div>
              <div className="text-gray-300">Active Members</div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold text-primary mb-2">50+</div>
              <div className="text-gray-300">Expert Trainers</div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold text-primary mb-2">10k+</div>
              <div className="text-gray-300">Sessions Booked</div>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
