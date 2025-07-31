import { promises as fs } from "fs"
import path from "path"
import { Metadata } from "next"
import Image from "next/image"
import { z } from "zod"

import { columns } from "@/components/columns"
import { DataTable } from "@/components/data-table"
import { UserNav } from "@/components/user-nav"
import { taskSchema } from "@/data/schema"
import { headers } from "next/headers";

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
                            Here&apos;s a list of your tasks for this month.
                        </p>
                    </div>
                    <div className="flex items-center gap-2">
                        <UserNav />
                    </div>
                </div>
                <DataTable data={tasks} columns={columns} />
            </div>
        </>
    )
}
