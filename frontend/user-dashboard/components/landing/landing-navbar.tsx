"use client"

import { useState } from "react"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Menu, X } from "lucide-react"

export function LandingNavbar() {
  const [open, setOpen] = useState(false)

  const scrollToId = (id: string) => {
    const el = document.getElementById(id)
    if (!el) return
    const headerOffset = 72
    const elementPosition = el.getBoundingClientRect().top + window.pageYOffset
    const offsetPosition = elementPosition - headerOffset
    window.scrollTo({ top: offsetPosition, behavior: "smooth" })
  }

  return (
    <header className="fixed top-0 left-0 right-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 mt-6">
        <div className="h-12 flex items-center justify-between rounded-full border border-border/60 bg-background shadow-md px-3 sm:px-5">
          <Link href="/" className="font-semibold text-sm sm:text-base">
            <span className="inline-flex items-center justify-center h-8 w-8 rounded-full bg-primary/15 text-primary mr-2">F+</span>
            <span className="hidden xs:inline">FitTrack</span><span className="hidden xs:inline text-primary">+</span>
          </Link>

          <nav className="hidden md:flex items-center gap-1.5">
            <button onClick={() => scrollToId("home")} className="px-3 h-8 rounded-full text-sm text-muted-foreground hover:text-primary hover:bg-primary/10 transition-colors">Home</button>
            <button onClick={() => scrollToId("features")} className="px-3 h-8 rounded-full text-sm text-muted-foreground hover:text-primary hover:bg-primary/10 transition-colors">Features</button>
            <button onClick={() => scrollToId("testimonials")} className="px-3 h-8 rounded-full text-sm text-muted-foreground hover:text-primary hover:bg-primary/10 transition-colors">Testimonials</button>
            <button onClick={() => scrollToId("get-started")} className="px-3 h-8 rounded-full text-sm text-muted-foreground hover:text-primary hover:bg-primary/10 transition-colors">Get Started</button>
          </nav>

          <div className="hidden md:flex items-center gap-1.5">
            <Button asChild variant="ghost" className="h-8 rounded-full px-3 hover:text-primary">
              <Link href="/auth/login">Sign in</Link>
            </Button>
            <Button asChild className="h-8 rounded-full px-3 bg-primary text-primary-foreground hover:bg-primary/90">
              <Link href="/auth/register">Create account</Link>
            </Button>
          </div>

          <Button variant="ghost" size="icon" className="md:hidden h-8 w-8 rounded-full" onClick={() => setOpen(v => !v)} aria-label="Toggle menu">
            {open ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
          </Button>
        </div>
      </div>

      {open && (
        <div className="md:hidden border-t border-border/50 bg-background/95 backdrop-blur">
          <div className="max-w-7xl mx-auto px-6 py-4 flex flex-col gap-3">
            <button onClick={() => { setOpen(false); scrollToId("home") }} className="h-10 rounded-full text-left px-3 hover:bg-primary/10 hover:text-primary">Home</button>
            <button onClick={() => { setOpen(false); scrollToId("features") }} className="h-10 rounded-full text-left px-3 hover:bg-primary/10 hover:text-primary">Features</button>
            <button onClick={() => { setOpen(false); scrollToId("testimonials") }} className="h-10 rounded-full text-left px-3 hover:bg-primary/10 hover:text-primary">Testimonials</button>
            <button onClick={() => { setOpen(false); scrollToId("get-started") }} className="h-10 rounded-full text-left px-3 hover:bg-primary/10 hover:text-primary">Get Started</button>
            <div className="flex gap-3 pt-2">
              <Button asChild variant="ghost" className="flex-1 h-10 rounded-full">
                <Link href="/auth/login" onClick={() => setOpen(false)}>Sign in</Link>
              </Button>
              <Button asChild className="flex-1 h-10 rounded-full">
                <Link href="/auth/register" onClick={() => setOpen(false)}>Create account</Link>
              </Button>
            </div>
          </div>
        </div>
      )}
    </header>
  )
}


