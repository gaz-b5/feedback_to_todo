"use client";

import { Mail } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
    Sheet,
    SheetContent,
    SheetDescription,
    SheetHeader,
    SheetTitle,
    SheetTrigger,
} from "@/components/ui/sheet";
import { TaskMailsList } from "@/components/task-mails-list"; // ⬅️ Import the child we made
import { ScrollArea } from "./ui/scroll-area";

/**
 * Props:
 * - taskId: The ID of the task whose mails should be displayed
 */
export function ShowMail({ taskId }: { taskId: string }) {
    return (
        <Sheet>
            {/* Trigger button */}
            <SheetTrigger asChild>
                <Button
                    variant="ghost"
                    size="icon"
                    className="data-[state=open]:bg-muted size-8"
                >
                    <Mail />
                    <span className="sr-only">Open mails</span>
                </Button>
            </SheetTrigger>

            {/* Sheet sliding panel */}
            <SheetContent className="sm:max-w-[600px]">
                <SheetHeader>
                    <SheetTitle>Linked Messages</SheetTitle>
                    <SheetDescription>
                        Below are the messages linked to this task.
                    </SheetDescription>
                </SheetHeader>

                {/* Our new component that handles fetching & card rendering */}
                <div className="mt-4">
                    <ScrollArea>
                        <TaskMailsList taskId={taskId} />
                    </ScrollArea>
                </div>
            </SheetContent>
        </Sheet>
    );
}
