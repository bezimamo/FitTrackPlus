"use client";

import Navbar from "@/components/Navbar";
import Hero from "@/components/dashboard/landingPage/Hero";
import Features from "@/components/dashboard/landingPage/Features";
import HowItWorks from "@/components/dashboard/landingPage/HowItWorks";
import PlansPreview from "@/components/dashboard/landingPage/PlansPreview";
import CTA from "@/components/dashboard/landingPage/CTA";

export default function LandingPage() {
  return (
    <div className="overflow-x-hidden">
      <Navbar /> {/* Only Navbar, no Sidebar */}
      <main className="pt-24 md:pt-32 px-4 md:px-8 lg:px-16">
        <Hero />
        <Features />
        <HowItWorks />
        <PlansPreview />
        <CTA />
      </main>
    </div>
  );
}
