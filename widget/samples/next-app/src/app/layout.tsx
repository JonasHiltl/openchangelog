import type { Metadata } from "next";
import { Inter } from 'next/font/google'
import "./globals.css";
import { NavigationMenuDemo } from "@/components/navigation-menu-demo";
import { cn } from "@/lib/utils";

export const metadata: Metadata = {
  title: "Next Changelog Demo",
  description: "Embed Openchangelog Changelog into Next.js app",
};

const inter = Inter({ subsets: ['latin'] })

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body className={cn(inter.className)}>
        <NavigationMenuDemo />
        {children}
      </body>
    </html>
  );
}
