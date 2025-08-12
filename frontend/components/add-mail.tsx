"use client";

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
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import React from "react";
import { MailPlus } from "lucide-react";
import { Separator } from "./ui/separator";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { Textarea } from "@/components/ui/textarea";

interface AddMailProps {
    projectId: string;
}


export function AddMail({ projectId }: AddMailProps) {
    const router = useRouter();

    async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault();

        const form = event.currentTarget;

        // Extract values
        const content = (form.elements.namedItem("content") as HTMLTextAreaElement).value.trim();
        const project_id = projectId;

        // Prepare the request body
        const body = JSON.stringify({ project_id, content });

        // Send the POST request to the API
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/feedback/new`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body,
        });

        if (!response.ok) {
            router.refresh();
            console.error("Failed to add mail");
            return;
        }

        // Handle success (e.g., close dialog, refresh data)
        router.refresh();
    }

    return (
        <Dialog>
            <DialogTrigger asChild>
                <Button variant="default" size="sm"><MailPlus />Add Mail</Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[600px] sm:max-h-[600px]">
                <DialogHeader>
                    <DialogTitle>Add Mail</DialogTitle>
                    <DialogDescription>
                        Enter the content of the mail to be added.
                    </DialogDescription>
                </DialogHeader>
                <Separator className="my-4" />
                <div className="text-sm text-red">Please beware that it takes about a minute to process the sent message, so it might seem that it is not working.</div>

                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="grid gap-3">
                        <Label htmlFor="content">Message</Label>
                        <Textarea
                            id="content"
                            name="content"
                            placeholder="Enter your message here..."
                            required
                            className="h-40"
                        />
                    </div>
                    <DialogFooter>

                        <DialogClose asChild>
                            <Button type="submit">Add Mail</Button>
                        </DialogClose>
                        <DialogClose asChild>
                            <Button variant="outline">Cancel</Button>
                        </DialogClose>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
}