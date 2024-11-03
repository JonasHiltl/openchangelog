"use client"

import * as React from "react"
import Link from "next/link"

import {
    NavigationMenu,
    NavigationMenuItem,
    NavigationMenuLink,
    NavigationMenuList,
    navigationMenuTriggerStyle,
    NavigationMenuTrigger
} from "@/components/ui/navigation-menu"
import { Button } from "./ui/button"

export function NavigationMenuDemo() {
    return (
        <nav className="border-b border-border">
            <div className="w-full mx-auto max-w-screen-lg flex justify-between py-4">
                <a href="/">
                    <h1 className="text-3xl font-bold">Logo</h1>
                </a>
                <div className="flex gap-8">
                    <NavigationMenu>
                        <NavigationMenuList>
                            <NavigationMenuItem>
                                <NavigationMenuTrigger>Products</NavigationMenuTrigger>
                            </NavigationMenuItem>
                            <NavigationMenuItem>
                                <Link href="/changelog" legacyBehavior passHref>
                                    <NavigationMenuLink className={navigationMenuTriggerStyle()}>
                                        Changelog
                                    </NavigationMenuLink>
                                </Link>
                            </NavigationMenuItem>
                            <NavigationMenuItem>
                                <NavigationMenuLink className={navigationMenuTriggerStyle()}>
                                    Pricing
                                </NavigationMenuLink>
                            </NavigationMenuItem>
                        </NavigationMenuList>
                    </NavigationMenu>
                    <Button variant="default">Sign up</Button>
                </div>
            </div>
        </nav>
    )
}
