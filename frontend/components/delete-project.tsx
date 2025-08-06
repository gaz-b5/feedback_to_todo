// delete project functionality
"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Trash, Trash2, X } from "lucide-react";

const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || "http://localhost:3000";

export function DeleteProjectButton({
    projectId,
    // onDeleteSuccess,
}: {
    projectId: string;
    // onDeleteSuccess: () => void;
}) {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    async function handleDeleteProject() {
        setLoading(true);
        setError(null);

        try {
            const res = await fetch(`${baseUrl}/api/projects/delete`, {
                method: "DELETE",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ project_id: projectId }),
            });

            if (!res.ok) {
                throw new Error("Failed to delete project");
            }

            // onDeleteSuccess();
        } catch (err) {
            setError("Network error: Failed to delete project");
        } finally {
            setLoading(false);
        }
    }

    return (
        <Button
            variant="destructive"
            className="w-full"
            onClick={handleDeleteProject}
            disabled={loading}
        >
            <a href="/dashboard" className="w-full">
                <div className="w-full flex flex-row items-center justify-between">Delete Project
                    <Trash2 />
                </div>
            </a>
            {error && <p className="text-red-500" > {error} </p>
            }
        </Button>
    );
}