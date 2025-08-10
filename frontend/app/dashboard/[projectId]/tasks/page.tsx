import { Metadata } from "next"
import Image from "next/image"
import { z } from "zod"

import { columns } from "@/components/columns"
import { DataTable } from "@/components/data-table"
import { taskSchema } from "@/data/schema"
import { headers } from "next/headers";
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
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import { Delete, Settings } from 'lucide-react';
import { Separator } from "@/components/ui/separator"
import { ScrollArea } from "@/components/ui/scroll-area"
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { roles } from "@/data/data"
import { AddMemberForm } from "@/components/add-member-form"
import { EditMembers } from "@/components/edit-members-form"
import { MembersSection } from "@/components/project-settings"
import { DeleteProjectButton } from "@/components/delete-project"


const baseUrl = process.env.NEXT_S_PUBLIC_BASE_URL || "http://localhost:3000";

interface PageProps {
    params: {
        projectId: string;
    };
}

export const metadata: Metadata = {
    title: "Tasks",
    description: "A task and issue tracker build using Tanstack Table.",
}

// dummy data fetching function
// async function getTasks({ params }: PageProps) {
//     const { projectId } = params;

//     const res = await fetch(`${baseUrl}/api/tasks?projectId=${projectId}`, {
//         method: "GET",
//         credentials: "include", // Sends cookies securely (including httpOnly)
//     });

//     if (!res.ok) {
//         throw new Error("Failed to fetch tasks " + res.statusText);
//     }

//     const tasks = await res.json();

//     // Validate the data using zod schema as before
//     return z.array(taskSchema).parse(tasks);
// }

export default async function TaskPage({ params }: PageProps) {

    const hdrs = await headers();
    const cookie = hdrs.get("cookie");

    const { projectId } = await params
    const res = await fetch(`${baseUrl}/api/tasks?projectId=${projectId}`, {
        method: "GET",
        headers: {
            cookie: cookie || "",      // <-- Forward cookies
        },
        cache: "no-store"
    });

    if (!res.ok) {
        throw new Error("Failed to fetch tasks " + res.statusText);
    }

    const tasks = z.array(taskSchema).parse(await res.json());

    return (
        <>
            <div className="md:hidden">
                <Image
                    src="/examples/tasks-light.png"
                    width={1280}
                    height={998}
                    alt="Playground"
                    className="block dark:hidden"
                />
                <Image
                    src="/examples/tasks-dark.png"
                    width={1280}
                    height={998}
                    alt="Playground"
                    className="hidden dark:block"
                />
            </div>
            <div className="hidden h-full flex-1 flex-col gap-8 p-8 md:flex">
                <div className="flex items-center justify-between gap-2">
                    <div className="flex flex-col gap-1">
                        <h2 className="text-2xl font-semibold tracking-tight">
                            Welcome back!
                        </h2>
                        <p className="text-muted-foreground">
                            Here&apos;s a list of the tasks for your project.
                        </p>
                    </div>
                    <div className="flex items-center gap-2">
                        <Dialog>
                            <DialogTrigger asChild>
                                <Button variant="secondary" className="h-9 px-4">
                                    <Settings />
                                </Button>
                            </DialogTrigger>


                            <DialogContent className="sm:max-w-[600px] sm:max-h-[600px]" showCloseButton={false}>

                                <ScrollArea>
                                    <DialogHeader>
                                        <DialogTitle>Project settings</DialogTitle>
                                    </DialogHeader>

                                    <Separator className="my-8" />

                                    <MembersSection projectId={projectId} />

                                    <Separator className="my-8" />

                                    <DeleteProjectButton projectId={projectId} />

                                </ScrollArea>
                            </DialogContent>

                        </Dialog >
                    </div >
                </div >
                <DataTable data={tasks} columns={columns} params={{ projectId }} />
            </div >
        </>
    )
}

