"use client"

import { Row } from "@tanstack/react-table"
import { MoreHorizontal } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

import { labels, statuses, priorities } from "../data/data"
import { taskSchema } from "../data/schema"
import { useRouter } from "next/navigation";

const baseUrl = process.env.NEXT_S_PUBLIC_BASE_URL || "http://localhost:3000";

interface DataTableRowActionsProps<TData> {
  row: Row<TData>
}

export function DataTableRowActions<TData>({
  row,
}: DataTableRowActionsProps<TData>) {
  const router = useRouter();
  const task = taskSchema.parse(row.original)

  // Generic update function, calls router.refresh() after success
  async function updateTask(fields: Partial<{ status: string; priority: string; nature: string }>) {
    try {
      const res = await fetch(`${baseUrl}/api/tasks/update`, {
        method: "PATCH",
        credentials: "include", // Sends cookies securely (including httpOnly)
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          task_id: task.id,
          ...fields,
        }),
      });

      if (!res.ok) {
        throw new Error("Update failed");
      }

      router.refresh(); // Forces table/page to reload new data
    } catch (err) {
      alert("Failed to update task");
      console.error(err);
    }
  }

  async function deleteTask() {
    if (!confirm("Are you sure you want to delete this task?")) {
      return; // Cancel delete if user aborts
    }

    try {
      const res = await fetch(`${baseUrl}/api/tasks/delete`, {
        method: "DELETE",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          task_id: task.id,
        }),
      });

      if (!res.ok) {
        throw new Error("Delete failed");
      }

      router.refresh(); // Refresh table to reflect deletion
    } catch (err) {
      alert("Failed to delete task");
      console.error(err);
    }
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          className="data-[state=open]:bg-muted size-8"
        >
          <MoreHorizontal />
          <span className="sr-only">Open menu</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-[160px]">
        {/* <DropdownMenuItem>Edit</DropdownMenuItem>
        <DropdownMenuItem>Make a copy</DropdownMenuItem>
        <DropdownMenuItem>Favorite</DropdownMenuItem> */}
        {/* <DropdownMenuSeparator /> */}
        {/* <DropdownMenuSub>
          <DropdownMenuSubTrigger>Nature</DropdownMenuSubTrigger>
          <DropdownMenuSubContent>
            <DropdownMenuRadioGroup value={task.nature}>
              {labels.map((label) => (
                <DropdownMenuRadioItem key={label.value} value={label.value}>
                  {label.label}
                </DropdownMenuRadioItem>
              ))}
            </DropdownMenuRadioGroup>
          </DropdownMenuSubContent>
        </DropdownMenuSub> */}
        <DropdownMenuSub>
          <DropdownMenuSubTrigger>Status</DropdownMenuSubTrigger>
          <DropdownMenuSubContent>
            <DropdownMenuRadioGroup
              value={task.status}
              onValueChange={value => updateTask({ status: value })}
            >
              {statuses.map((status) => (
                <DropdownMenuRadioItem key={status.value} value={status.value}>
                  {status.label}
                </DropdownMenuRadioItem>
              ))}
            </DropdownMenuRadioGroup>
          </DropdownMenuSubContent>
        </DropdownMenuSub>
        <DropdownMenuSub>
          <DropdownMenuSubTrigger>Priority</DropdownMenuSubTrigger>
          <DropdownMenuSubContent>
            <DropdownMenuRadioGroup
              value={task.priority}
              onValueChange={value => updateTask({ priority: value })}
            >
              {priorities.map((priority) => (
                <DropdownMenuRadioItem key={priority.value} value={priority.value}>
                  {priority.label}
                </DropdownMenuRadioItem>
              ))}
            </DropdownMenuRadioGroup>
          </DropdownMenuSubContent>
        </DropdownMenuSub>
        <DropdownMenuSeparator />
        <DropdownMenuItem variant="destructive" onClick={deleteTask}>
          Delete
          <DropdownMenuShortcut>⌘⌫</DropdownMenuShortcut>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
