"use client"

import "./globals.css"
import { usePathname } from "next/navigation"
import { AppSidebar } from "@/components/app-sidebar"
import DashboardNavbar from "@/components/DashboardNavbar"
import { SidebarProvider } from "@/components/ui/sidebar"

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const isDashboard = pathname?.startsWith("/dashboard")

  return (
    <html lang="en">
      <body className="flex flex-col min-h-screen">
        <SidebarProvider>
          {isDashboard ? (
            <div className="flex flex-1 relative">
              {/* Sidebar stays fixed */}
              <AppSidebar />

              {/* Content area with left margin so it's not covered */}
              <div className="flex-1 flex flex-col ml-64 min-h-screen">
                <DashboardNavbar />
                <main className="flex-1 pt-20 px-4 md:px-8">{children}</main>
                
                {/* Footer wrapper with z-index and full width */}
                <div className="w-full z-10">
                </div>
              </div>
            </div>
          ) : (
            <div className="flex flex-col flex-1">
              <main className="flex-1 p-0">{children}</main>
            </div>
          )}
        </SidebarProvider>
      </body>
    </html>
  )
}
