"use client";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableFooter,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { roles } from "@/data/data";
import { Trash } from "lucide-react";

const baseUrl = process.env.NEXT_S_PUBLIC_BASE_URL || "http://localhost:3000";

type Member = {
    user_id: string;
    name: string;
    email: string;
    role: string;
};

export function EditMembers({ projectId }: { projectId: string }) {
    const [members, setMembers] = useState<Member[]>([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        async function fetchMembers() {
            setLoading(true);
            try {
                const res = await fetch(`${baseUrl}/api/projects/members/getall?projectId=${projectId}`, {
                    method: "GET",
                    credentials: "include",
                });
                if (!res.ok) throw new Error("Failed to fetch members");
                const data = await res.json();
                setMembers(data);
            } catch (err) {
                // Optionally show error message
            } finally {
                setLoading(false);
            }
        }

        fetchMembers();
    }, [projectId]);

    async function handleRemoveMember(memberId: string) {
        try {
            const res = await fetch(`${baseUrl}/api/projects/members/remove`, {
                method: "DELETE",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ member_id: memberId, project: projectId }),
            });
            if (!res.ok) {
                alert(res.statusText);
                throw new Error("Failed to remove member");
            };
            setMembers((prev) => prev.filter((m) => m.user_id !== memberId));
        } catch (err) {
            // Handle error (show notification etc.)
        }
    }

    async function handleRoleChange(memberId: string, newRole: string) {
        try {
            const res = await fetch(`${baseUrl}/api/projects/members/update-role`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ member_id: memberId, project_id: projectId, role: newRole }),
            });
            if (!res.ok) throw new Error("Failed to update member role");

            // Update local state if successful
            setMembers((prev) =>
                prev.map((m) => (m.user_id === memberId ? { ...m, role: newRole } : m))
            );
        } catch (err) {
            // Optionally handle error
        }
    }

    return (
        <Table>
            <TableHeader>
                <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Email</TableHead>
                    <TableHead>Role</TableHead>
                    <TableHead>Remove</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {loading ? (
                    <TableRow>
                        <TableCell colSpan={4} className="text-center">
                            Loading...
                        </TableCell>
                    </TableRow>
                ) : members.length === 0 ? (
                    <TableRow>
                        <TableCell colSpan={4} className="text-center">
                            No members found.
                        </TableCell>
                    </TableRow>
                ) : (
                    members.map((member) => (
                        <TableRow key={member.user_id}>
                            <TableCell>
                                {member.name}
                            </TableCell>
                            <TableCell>{member.email}</TableCell>
                            <TableCell>
                                <Select
                                    value={member.role}
                                    onValueChange={(value) => handleRoleChange(member.user_id, value)}
                                    name="role"
                                >
                                    <SelectTrigger>
                                        <SelectValue placeholder="Select role" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectGroup>
                                            <SelectLabel>Role</SelectLabel>
                                            {roles.map((roleOption) => (
                                                <SelectItem key={roleOption.value} value={roleOption.value}>
                                                    {roleOption.label}
                                                </SelectItem>
                                            ))}
                                        </SelectGroup>
                                    </SelectContent>
                                </Select>
                            </TableCell>
                            <TableCell>
                                <Button
                                    variant="destructive"
                                    size="sm"
                                    onClick={() => handleRemoveMember(member.user_id)}
                                >
                                    <Trash />
                                </Button>
                            </TableCell>
                        </TableRow>
                    ))
                )}
            </TableBody>
        </Table>
    );
}
