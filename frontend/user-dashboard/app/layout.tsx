"use client"; // <-- Make this a client component

import './globals.css';
import { usePathname } from 'next/navigation';
import Navbar from '../components/Navbar';
import Sidebar from '../components/Sidebar';
import DashboardNavbar from '../components/DashboardNavbar';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();

  const isDashboard = pathname?.startsWith('/dashboard');

  return (
    <html lang="en">
      <body className="flex flex-col min-h-screen">
        {isDashboard ? (
          <div className="flex">
            <Sidebar />
            <div className="flex-1 flex flex-col">
              <DashboardNavbar />
              <main className="pt-20 px-4 md:px-8">{children}</main>
            </div>
          </div>
        ) : (
          <div className="flex flex-col">
            <Navbar />
            <main className="pt-24 md:pt-32 px-4 md:px-8 lg:px-16">{children}</main>
          </div>
        )}
      </body>
    </html>
  );
}
