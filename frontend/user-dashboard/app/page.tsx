"use client";
import { HeroSection } from "@/components/landing/hero-section"
import { FeaturesSection } from "@/components/landing/features-section"
import { TestimonialsSection } from "@/components/landing/testimonials-section"
import { CTASection } from "@/components/landing/cta-section"
import { LandingNavbar } from "@/components/landing/landing-navbar"
import { SiteFooter } from "@/components/footer"

export default function LandingPage() {
  return (
    <div>
    <LandingNavbar />
    <main className="min-h-screen">
      <HeroSection />
      <FeaturesSection />
      <TestimonialsSection />
      <CTASection />
    </main>
    <SiteFooter />
    </div>
  )
}
