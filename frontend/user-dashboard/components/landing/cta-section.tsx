import { Button } from "@/components/ui/button"
import { ArrowRight, CheckCircle2 } from "lucide-react"
import Link from "next/link"

export function CTASection() {
  return (
    <section id="get-started" className="py-24">
      <div className="max-w-7xl mx-auto px-6">
        <div className="relative overflow-hidden rounded-2xl border border-emerald-900/30 bg-[linear-gradient(to_bottom_right,rgba(2,44,34,0.95),rgba(15,23,42,0.95))] text-white p-8 md:p-12">
          <div className="pointer-events-none absolute -right-16 -top-16 h-48 w-48 rounded-full bg-white/5 blur-3xl" />
          <div className="pointer-events-none absolute -left-20 -bottom-20 h-52 w-52 rounded-full bg-white/5 blur-3xl" />

          <div className="relative grid grid-cols-1 lg:grid-cols-3 gap-10 items-center">
            <div className="lg:col-span-2 text-center lg:text-left">
              <h2 className="text-3xl md:text-4xl font-bold tracking-tight">
                Level up your fitness —
                <span className="block text-primary">the elegant way</span>
              </h2>
              <p className="mt-4 text-base md:text-lg text-white/80 max-w-2xl mx-auto lg:mx-0">
                Create momentum with smart plans, seamless booking, and beautifully simple progress tracking. Stay
                consistent without the clutter.
              </p>

              <div className="mt-7 grid grid-cols-1 sm:grid-cols-3 gap-3">
                <div className="rounded-xl border border-white/10 bg-white/5 p-4 text-left">
                  <div className="text-sm text-white/90">No setup</div>
                  <div className="text-xs text-white/60">Start in minutes</div>
                </div>
                <div className="rounded-xl border border-white/10 bg-white/5 p-4 text-left">
                  <div className="text-sm text-white/90">Flexible booking</div>
                  <div className="text-xs text-white/60">24/7 scheduling</div>
                </div>
                <div className="rounded-xl border border-white/10 bg-white/5 p-4 text-left">
                  <div className="text-sm text-white/90">Real progress</div>
                  <div className="text-xs text-white/60">Clear insights</div>
                </div>
              </div>

              <div className="mt-7 flex flex-col sm:flex-row gap-3 justify-center lg:justify-start">
                <Button asChild size="lg" className="rounded-full px-7 h-11 bg-primary text-primary-foreground hover:bg-primary/90">
                  <Link href="/auth/register" className="flex items-center gap-2">
                    Create your account
                    <ArrowRight className="h-4 w-4" />
                  </Link>
                </Button>
                <Button asChild variant="secondary" size="lg" className="rounded-full px-7 h-11 bg-white text-black hover:bg-white/90">
                  <Link href="/auth/login">I already have an account</Link>
                </Button>
              </div>
            </div>

            <div className="hidden lg:block">
              <div className="rounded-xl border border-white/10 bg-white/5 p-6">
                <div className="grid grid-cols-3 gap-6 text-center">
                  <div>
                    <div className="text-3xl font-bold text-primary">500+</div>
                    <div className="text-sm text-white/70">Active Members</div>
                  </div>
                  <div>
                    <div className="text-3xl font-bold text-primary">50+</div>
                    <div className="text-sm text-white/70">Expert Trainers</div>
                  </div>
                  <div>
                    <div className="text-3xl font-bold text-primary">10k+</div>
                    <div className="text-sm text-white/70">Sessions Booked</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
