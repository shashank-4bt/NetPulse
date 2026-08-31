"use client";

import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

export function InteractiveControls() {
  return (
    <div className="grid gap-6 md:grid-cols-2">
      <div className="flex flex-col gap-2">
        <label htmlFor="probe-target" className="text-sm font-medium">
          Input
        </label>
        <Input
          id="probe-target"
          name="probe-target"
          placeholder="example.com"
          autoComplete="off"
        />
      </div>
      <div className="flex flex-col gap-2">
        <span className="text-sm font-medium" id="layer-label">
          Select
        </span>
        <Select defaultValue="dns" aria-labelledby="layer-label">
          <SelectTrigger className="w-full">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="dns">DNS</SelectItem>
            <SelectItem value="tls">TLS</SelectItem>
            <SelectItem value="http">HTTP</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="flex flex-col gap-2 md:col-span-2">
        <label htmlFor="notes" className="text-sm font-medium">
          Textarea
        </label>
        <Textarea id="notes" name="notes" placeholder="Operator notes" />
      </div>
      <label className="flex min-h-11 items-center gap-2 text-sm">
        <Checkbox defaultChecked={false} />
        Include traceroute when available
      </label>
      <label className="flex min-h-11 items-center gap-2 text-sm">
        <Switch aria-label="Verbose probe logging" />
        Verbose probe logging
      </label>
      <div className="flex flex-wrap gap-2">
        <Dialog>
          <DialogTrigger render={<Button variant="outline" />}>
            Open dialog
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Confirm probe window</DialogTitle>
              <DialogDescription>
                This dialog is a design-system example. It does not start a
                measurement.
              </DialogDescription>
            </DialogHeader>
            <DialogFooter showCloseButton>
              <Button type="button" variant="outline" disabled>
                Start probe (unavailable)
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
        <Drawer>
          <DrawerTrigger render={<Button variant="outline" />}>
            Open drawer
          </DrawerTrigger>
          <DrawerContent>
            <DrawerHeader>
              <DrawerTitle>Layer detail</DrawerTitle>
              <DrawerDescription>
                Example drawer. No live layer data is loaded.
              </DrawerDescription>
            </DrawerHeader>
            <DrawerFooter>
              <DrawerClose render={<Button variant="outline" />}>
                Close
              </DrawerClose>
            </DrawerFooter>
          </DrawerContent>
        </Drawer>
        <DropdownMenu>
          <DropdownMenuTrigger render={<Button variant="outline" />}>
            Dropdown
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem>Copy run id</DropdownMenuItem>
            <DropdownMenuItem>Export evidence</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <Tooltip>
          <TooltipTrigger render={<Button variant="outline" />}>
            Tooltip
          </TooltipTrigger>
          <TooltipContent>Evidence class is always labeled in text</TooltipContent>
        </Tooltip>
        <Button
          type="button"
          onClick={() =>
            toast("Toast example", {
              description: "This is a UI notification, not a live incident.",
            })
          }
        >
          Show toast
        </Button>
      </div>
    </div>
  );
}
