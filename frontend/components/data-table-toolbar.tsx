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
    const nature = (form.elements.namedItem("nature") as RadioNodeList | null)?.value ?? labels[0]?.value
    const priority = (form.elements.namedItem("priority") as RadioNodeList | null)?.value ?? priorities[0]?.value
    const status = (form.elements.namedItem("status") as RadioNodeList | null)?.value ?? statuses[0]?.value
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
        <Dialog >
          {/* <form onSubmit={handleSubmit}> */}
          <DialogTrigger asChild>
            <Button size="sm">Add Task</Button>
          </DialogTrigger>

          <DialogContent className="sm:max-w-[600px]">

            <DialogHeader>
              <DialogTitle>Add a task</DialogTitle>
              <DialogDescription>
                Add a task, which will be added to the tasklist unedited.
              </DialogDescription>

            </DialogHeader>
            <form onSubmit={handleSubmit}>
              <div className="grid gap-4">
                {/* Task description textarea */}
                <div className="grid gap-3 ">
                  <Label htmlFor="description-1">Task Description</Label>
                  <Textarea
                    id="description-1"
                    name="description-1"
                    placeholder="Type task description here."
                  />
                </div>

                {/* Container for the three radio groups side by side */}
                <div className="grid grid-cols-3 gap-8">
                  {/* Labels RadioGroup */}
                  <div className="grid gap-3 place-items-start grid-flow-col grid-rows-5">
                    <Label className="row-span-1">Nature</Label>
                    <RadioGroup
                      aria-label="Nature"
                      name="nature"
                      defaultValue={labels[0]?.value}
                      className="flex flex-col space-y-2 row-span-4"
                    >
                      {labels.map((label) => (
                        <div key={label.value} className="flex items-center space-x-2">
                          <RadioGroupItem value={label.value} id={`nature-${label.value}`} />
                          <Label htmlFor={`nature-${label.value}`}>{label.label}</Label>
                        </div>
                      ))}
                    </RadioGroup>
                  </div>

                  {/* Priorities RadioGroup */}
                  <div className="grid gap-1 place-items-start grid-flow-col grid-rows-5">
                    <Label className="row-span-1">Priority</Label>
                    <RadioGroup
                      aria-label="Priority"
                      name="priority"
                      defaultValue={priorities[0]?.value}
                      className="flex flex-col space-y-2 row-span-4"
                    >
                      {priorities.map((priority) => (
                        <div key={priority.value} className="flex items-center space-x-2">
                          <RadioGroupItem value={priority.value} id={`priority-${priority.value}`} />
                          <Label htmlFor={`priority-${priority.value}`}>{priority.label}</Label>
                        </div>
                      ))}
                    </RadioGroup>
                  </div>

                  {/* Statuses RadioGroup */}
                  <div className="grid grid-rows-5 gap-1 place-items-start grid-flow-col">
                    <Label className="row-span-1">Status</Label>
                    <RadioGroup
                      aria-label="Status"
                      name="status"
                      defaultValue={statuses[0]?.value}
                      className="flex flex-col space-y-2 row-span-4"
                    >
                      {statuses.map((status) => (
                        <div key={status.value} className="flex items-center space-x-2">
                          <RadioGroupItem value={status.value} id={`status-${status.value}`} />
                          <Label htmlFor={`status-${status.value}`}>{status.label}</Label>
                        </div>
                      ))}
                    </RadioGroup>
                  </div>
                </div>
              </div>

              <DialogFooter>
                <DialogClose asChild>
                  <Button variant="outline">Cancel</Button>
                </DialogClose>
                <DialogClose asChild>
                  <Button type="submit">Add task</Button>
                </DialogClose>
              </DialogFooter>
            </form>
          </DialogContent>
          {/* </form> */}
        </Dialog>

      </div>
    </div >
  )
}
