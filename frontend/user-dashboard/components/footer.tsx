import Link from "next/link"

export function SiteFooter() {
	return (
		<footer className="border-t border-emerald-900/30 bg-[linear-gradient(to_bottom_right,rgba(2,44,34,1),rgba(15,23,42,0.96))] text-white">
			<div className="max-w-7xl mx-auto px-6 py-12 grid grid-cols-1 md:grid-cols-4 gap-8">
				<div>
					<div className="font-bold text-lg">FitTrack<span className="text-primary">+</span></div>
					<p className="mt-3 text-sm text-white/70 max-w-xs">
						Track progress, book sessions, and achieve your goals with a modern, elegant platform.
					</p>
				</div>
				<div>
					<div className="font-semibold mb-3">Product</div>
					<ul className="space-y-2 text-sm text-white/70">
						<li><Link href="#features" className="hover:text-primary">Features</Link></li>
						<li><Link href="#testimonials" className="hover:text-primary">Testimonials</Link></li>
						<li><Link href="#get-started" className="hover:text-primary">Get Started</Link></li>
					</ul>
				</div>
				<div>
					<div className="font-semibold mb-3">Company</div>
					<ul className="space-y-2 text-sm text-white/70">
						<li><Link href="/about" className="hover:text-primary">About</Link></li>
						<li><Link href="/contact" className="hover:text-primary">Contact</Link></li>
						<li><Link href="/careers" className="hover:text-primary">Careers</Link></li>
					</ul>
				</div>
				<div>
					<div className="font-semibold mb-3">Legal</div>
					<ul className="space-y-2 text-sm text-white/70">
						<li><Link href="/terms" className="hover:text-primary">Terms</Link></li>
						<li><Link href="/privacy" className="hover:text-primary">Privacy</Link></li>
					</ul>
				</div>
			</div>
			<div className="border-t border-white/10">
				<div className="max-w-7xl mx-auto px-6 py-6 text-xs text-white/70 flex flex-col items-center justify-center gap-3 text-center">
					<p>© {new Date().getFullYear()} FitTrack+. All rights reserved.</p>
					<div className="flex items-center gap-4">
						<Link href="/auth/login" className="hover:text-primary">Sign In</Link>
						<Link href="/auth/register" className="hover:text-primary">Create account</Link>
					</div>
				</div>
			</div>
		</footer>
	)
}
