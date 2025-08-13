"use client"

import { Table } from "@tanstack/react-table"
import { X } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { DataTableViewOptions } from "@/components/data-table-view-options"

import { labels, priorities, statuses } from "../data/data"
import { DataTableFacetedFilter } from "./data-table-faceted-filter"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { useRouter } from "next/navigation"
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Separator } from "./ui/separator"
import { AddMail } from "@/components/add-mail"
import { BookmarkCheck } from "lucide-react"

const baseUrl = process.env.NEXT_S_PUBLIC_BASE_URL || "http://localhost:3000"


interface DataTableToolbarProps<TData> {
  table: Table<TData>;
  params: { projectId: string; };
}

export function DataTableToolbar<TData>({
  table,
  params
}: DataTableToolbarProps<TData>) {
  const isFiltered = table.getState().columnFilters.length > 0

  const router = useRouter()

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const form = event.currentTarget

    // Extract values
    const description = (form.elements.namedItem("description-1") as HTMLTextAreaElement).value.trim()
    const nature = (form.elements.namedItem("nature") as HTMLSelectElement | null)?.value ?? labels[0]?.value;
    const priority = (form.elements.namedItem("priority") as HTMLSelectElement | null)?.value ?? priorities[0]?.value;
    const status = (form.elements.namedItem("status") as HTMLSelectElement | null)?.value ?? statuses[0]?.value;
    const { projectId } = await params
    const project_id = projectId

    // Basic validation (add more if needed)
    if (!description) {
      alert("Please enter a task description.")
      return
    }

    // Prepare payload
    const payload = {
      project_id,
      description,
      nature,
      priority,
      status,
      // Add any other required fields, e.g., project_id if you have it in context
    }

    try {
      const res = await fetch(`${baseUrl}/api/tasks/add`, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      })

      if (!res.ok) {
        const errorData = await res.json().catch(() => null)
        alert(`Failed to add task: ${errorData?.error ?? res.statusText}`)
        return
      }

      form.reset()
      router.refresh()

      // If you control the dialog opening state, close it here as well

    } catch (error) {
      alert("An unexpected error occurred.")
      console.error(error)
    }
  }

  return (
    <div className="flex items-center justify-between">
      <div className="flex flex-1 items-center gap-2">
        <Input
          placeholder="Filter tasks..."
          value={(table.getColumn("title")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("title")?.setFilterValue(event.target.value)
          }
          className="h-8 w-[150px] lg:w-[250px]"
        />
        {table.getColumn("status") && (
          <DataTableFacetedFilter
            column={table.getColumn("status")}
            title="Status"
            options={statuses}
          />
        )}
        {table.getColumn("priority") && (
          <DataTableFacetedFilter
            column={table.getColumn("priority")}
            title="Priority"
            options={priorities}
          />
        )}
        {table.getColumn("nature") && (
          <DataTableFacetedFilter
            column={table.getColumn("nature")}
            title="Nature"
            options={labels}
          />
        )}
        {/* {table.getColumn("assigned") && (
          <DataTableFacetedFilter
            column={table.getColumn("assigned")}
            title="Assigned To"
            options={labels}
          />
        )} */}
        {isFiltered && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => table.resetColumnFilters()}
          >
            Reset
            <X />
          </Button>
        )}
      </div>
      <div className="flex items-center gap-2">
        <DataTableViewOptions table={table} />
        <AddMail projectId={params.projectId} />
        <Dialog >
          {/* <form onSubmit={handleSubmit}> */}
          <DialogTrigger asChild>
            <Button size="sm"> <BookmarkCheck />Add Task</Button>
          </DialogTrigger>

          <DialogContent className="sm:max-w-[600px]">

            <DialogHeader>
              <DialogTitle>Add a task</DialogTitle>
              <DialogDescription>
                Add a task, which will be added to the tasklist unedited.
              </DialogDescription>
            </DialogHeader>

            <Separator className="my-4" />

            <form onSubmit={handleSubmit}>
              <div className="grid gap-4">
                {/* Task description textarea */}
                <div className="grid gap-4 ">
                  <Label htmlFor="description-1">Task Description</Label>
                  <Textarea
                    id="description-1"
                    name="description-1"
                    placeholder="Type task description here."
                  />
                </div>

                <div className="grid grid-cols-3 gap-4">
                  <div className="flex flex-col gap-4">
                    <Label>Nature</Label>
                    <Select name="nature">
                      <SelectTrigger className="w-[180px]">
                        <SelectValue placeholder="Select nature" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectGroup>
                          {labels.map((label) => (
                            <SelectItem key={label.value} value={label.value}>{label.label}</SelectItem>
                          ))}
                        </SelectGroup>
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="flex flex-col gap-4">
                    <Label>Status</Label>
                    <Select name="status">
                      <SelectTrigger className="w-[180px]">
                        <SelectValue placeholder="Select status" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectGroup>
                          {statuses.map((status) => (
                            <SelectItem key={status.value} value={status.value}>{status.label}</SelectItem>
                          ))}
                        </SelectGroup>
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="flex flex-col gap-4">
                    <Label>Priority</Label>
                    <Select name="priority">
                      <SelectTrigger className="w-[180px]">
                        <SelectValue placeholder="Select priority" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectGroup>
                          {priorities.map((priority) => (
                            <SelectItem key={priority.value} value={priority.value}>{priority.label}</SelectItem>
                          ))}
                        </SelectGroup>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              </div>

              <Separator className="my-8" />

              <DialogFooter>
                <DialogClose asChild>
                  <Button type="submit">Add task</Button>
                </DialogClose>
                <DialogClose asChild>
                  <Button variant="outline">Cancel</Button>
                </DialogClose>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>

      </div>
    </div >
  )
}
