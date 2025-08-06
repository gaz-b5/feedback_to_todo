"use client";
import { useState, useCallback, useEffect } from "react";
import { AddMemberForm } from "@/components/add-member-form";
import { EditMembers } from "@/components/edit-members-form";
import { Separator } from "@/components/ui/separator";
import { DialogTitle } from "@/components/ui/dialog";

const baseUrl = process.env.NEXT_S_PUBLIC_BASE_URL || "http://localhost:3000";

export function MembersSection({ projectId }: { projectId: string }) {
    const [members, setMembers] = useState([]);
    const [loading, setLoading] = useState(false);

    const fetchMembers = useCallback(async () => {
        setLoading(true);
        try {
            const res = await fetch(`${baseUrl}/api/projects/members/getall?projectId=${projectId}`, {
                method: "GET",
                credentials: "include",
            });
            if (!res.ok) throw new Error("Failed to fetch members");
            const data = await res.json();
            setMembers(Array.isArray(data) ? data : data.members); // adapt for your API shape
        } catch (err) {
            // error handling if needed
        } finally {
            setLoading(false);
        }
    }, [projectId]);

    useEffect(() => {
        fetchMembers();
    }, [fetchMembers]);

    return (
        <div>
            <DialogTitle>Add member</DialogTitle>
            <AddMemberForm projectId={projectId} onMemberAdded={fetchMembers} />

            <Separator className="my-8" />

            <DialogTitle className="mb-2">Edit members</DialogTitle>
            <EditMembers
                projectId={projectId}
                members={members}
                loading={loading}
                refetchMembers={fetchMembers}
            />

        </div>
    );
}
