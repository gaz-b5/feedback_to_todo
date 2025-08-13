"use client"

import { ColumnDef } from "@tanstack/react-table"

import { Badge } from "@/components/ui/badge"
import { Checkbox } from "@/components/ui/checkbox"

import { priorities, statuses } from "../data/data"
import { Task } from "../data/schema"
import { DataTableColumnHeader } from "./data-table-column-header"
import { DataTableRowActions } from "./data-table-row-actions"
import { CheckIcon, ChevronsUpDownIcon } from "lucide-react"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"

import { useCallback, useEffect } from "react"

import React from "react"

const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || "http://localhost:3000";

type Member = {
  user_id: string;
  name: string;
  email: string;
  role: string;
};

export const columns: ColumnDef<Task>[] = [
  // {
  //   id: "select",
  //   header: ({ table }) => (
  //     <Checkbox
  //       checked={
  //         table.getIsAllPageRowsSelected() ||
  //         (table.getIsSomePageRowsSelected() && "indeterminate")
  //       }
  //       onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
  //       aria-label="Select all"
  //       className="translate-y-[2px]"
  //     />
  //   ),
  //   cell: ({ row }) => (
  //     <Checkbox
  //       checked={row.getIsSelected()}
  //       onCheckedChange={(value) => row.toggleSelected(!!value)}
  //       aria-label="Select row"
  //       className="translate-y-[2px]"
  //     />
  //   ),
  //   enableSorting: false,
  //   enableHiding: false,
  // },
  // {
  //   accessorKey: "id",
  //   header: ({ column }) => (
  //     <DataTableColumnHeader column={column} title="Task" />
  //   ),
  //   cell: ({ row }) => <div className="w-[80px]">{row.getValue("id")}</div>,
  //   enableSorting: false,
  //   enableHiding: false,
  // },
  {
    accessorKey: "nature",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Nature" />
    ),
    cell: ({ row }) => {

      return (
        <div className="flex gap-2">
          {row.original.nature && <Badge variant="outline">{row.original.nature}</Badge>}
        </div>
      )
    },
    enableSorting: false,
  },
  {
    accessorKey: "title",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Title" />
    ),
    cell: ({ row }) => {
      // const label = labels.find((label) => label.value === row.original.nature)

      return (
        <div className=" gap-2 ">
          {/* {row.original.nature && <Badge variant="outline">{row.original.nature}</Badge>} */}
          <div className="text-xs text-muted-foreground">{row.original.created}</div>
          <div className="max-w-[500px] truncate font-medium">
            {row.getValue("title")}
          </div>
        </div>
      )
    },
    enableSorting: false,
  },
  {
    accessorKey: "status",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Status" />
    ),
    cell: ({ row }) => {
      const status = statuses.find(
        (status) => status.value === row.getValue("status")
      )

      if (!status) {
        return null
      }

      return (
        <div className="flex w-[100px] items-center gap-2">
          {status.icon && (
            <status.icon className="text-muted-foreground size-4" />
          )}
          <span>{status.label}</span>
        </div>
      )
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id))
    },
    enableSorting: false,
  },
  {
    accessorKey: "priority",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Priority" />
    ),
    cell: ({ row }) => {
      const priority = priorities.find(
        (priority) => priority.value === row.getValue("priority")
      )

      if (!priority) {
        return null
      }

      return (
        <div className="flex items-center gap-2">
          {priority.icon && (
            <priority.icon className="text-muted-foreground size-4" />
          )}
          <span>{priority.label}</span>
        </div>
      )
    },
    filterFn: (row, id, value) => {
      return value.includes(row.getValue(id))
    },
  },
  {
    accessorKey: "assigned",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Assigned To" />
    ),
    cell: ({ row }) => {
      // const label = labels.find((label) => label.value === row.original.nature)
      const [open, setOpen] = React.useState(false)
      const [value, setValue] = React.useState(row.original.assigned)
      const [members, setMembers] = React.useState<Member[]>([]);

      const fetchMembers = useCallback(async () => {
        try {
          const res = await fetch(`${baseUrl}/api/projects/members/getall?projectId=${row.original.project}`, {
            method: "GET",
            credentials: "include",
          });
          if (!res.ok) throw new Error("Failed to fetch members");
          const data = await res.json();
          setMembers(Array.isArray(data) ? data : data.members); // adapt for your API shape
        } catch (err) {
          // error handling if needed
        }
      }, [row.original.project]);

      useEffect(() => {
        fetchMembers();
      }, [fetchMembers]);


      async function updateTask(fields: Partial<{ status: string; priority: string; nature: string; assigned: string }>) {
        try {
          console.log(fields)
          const res = await fetch(`${baseUrl}/api/tasks/update`, {
            method: "PATCH",
            credentials: "include", // Sends cookies securely (including httpOnly)
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              task_id: row.original.id,
              ...fields,
            }),
          });

          if (!res.ok) {
            throw new Error("Update failed");
          }
        } catch (err) {
          alert("Failed to update task");
          console.error(err);
        }
      }

      return (
        <div className="flex gap-2">
          <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
              <Button
                variant="outline"
                role="combobox"
                aria-expanded={open}
                className="w-[200px] justify-between"
              >
                {value ? members.find((member) => member.user_id === value)?.name : "Select member..."}
                <ChevronsUpDownIcon className="ml-2 h-4 w-4 shrink-0 opacity-50" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-[200px] p-0">
              <Command>
                <CommandInput placeholder="Search member..." />
                <CommandList>
                  <CommandEmpty>No member found.</CommandEmpty>
                  <CommandGroup>
                    {members.map((member) => (
                      <CommandItem
                        key={member.user_id}
                        value={member.user_id}
                        onSelect={(currentValue) => {
                          setValue(currentValue === value ? "" : currentValue)
                          updateTask({ assigned: currentValue === value ? "" : member.user_id })
                          setOpen(false)
                        }}
                      >
                        <CheckIcon
                          className={cn(
                            "mr-2 h-4 w-4",
                            value === member.user_id ? "opacity-100" : "opacity-0"
                          )}
                        />
                        {member.name}
                      </CommandItem>
                    ))}
                  </CommandGroup>
                </CommandList>
              </Command>
            </PopoverContent>
          </Popover>
          {/* <span className="max-w-[500px] truncate font-medium">
            {row.original.assigned ? row.original.assigned : "Unassigned"}
          </span> */}
        </div>
      )
    },
    enableSorting: true,
  },
  {
    accessorKey: "occurence",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Occurence" />
    ),
    cell: ({ row }) => {
      // const label = labels.find((label) => label.value === row.original.nature)

      return (
        <div className="flex gap-2">
          {/* {row.original.nature && <Badge variant="outline">{row.original.nature}</Badge>} */}
          <span className="max-w-[500px] truncate font-medium">
            {row.original.occurrence}
          </span>
        </div>
      )
    },
    enableSorting: true,
  },
  {
    id: "actions",
    cell: ({ row }) => {
      return (
        <>
          <DataTableRowActions row={row} />
        </>
      )
    },
  },
]
