"use client";

import Link from "next/link";
import { MenuIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";
import { publicConfig } from "@/lib/config/public";

type MobileNavProps = {
  links: readonly { href: string; label: string }[];
};

export function MobileNav({ links }: MobileNavProps) {
  return (
    <Drawer>
      <DrawerTrigger
        render={
          <Button
            type="button"
            variant="outline"
            size="icon"
            className="size-11 md:hidden"
            aria-label="Open menu"
          />
        }
      >
        <MenuIcon aria-hidden="true" />
      </DrawerTrigger>
      <DrawerContent className="md:hidden">
        <DrawerHeader>
          <DrawerTitle>{publicConfig.appName}</DrawerTitle>
        </DrawerHeader>
        <nav aria-label="Mobile" className="flex flex-col gap-1 px-4 pb-6">
          {links.map((item) => (
            <DrawerClose
              key={item.href}
              render={
                <Link
                  href={item.href}
                  className="inline-flex min-h-11 items-center rounded-md px-2 text-sm text-foreground hover:bg-muted"
                />
              }
            >
              {item.label}
            </DrawerClose>
          ))}
        </nav>
      </DrawerContent>
    </Drawer>
  );
}
