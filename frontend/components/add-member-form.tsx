"use client";
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
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

const baseUrl = process.env.NEXT_S_PUBLIC_BASE_URL || "http://localhost:3000";

export function AddMemberForm({ projectId, onMemberAdded }: { projectId: string; onMemberAdded: () => void }) {
    async function handleAddMember(event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault();
        const form = event.currentTarget;
        const email = (form.elements.namedItem("email") as HTMLInputElement)?.value.trim();
        const role = (form.elements.namedItem("role") as HTMLSelectElement)?.value;

        await fetch(`${baseUrl}/api/projects/members/add`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, role, project_id: projectId }),
        });
        form.reset();
        onMemberAdded(); // <--- refresh list!
    }

    return (
        <form className="flex flex-row w-full justify-between gap-4 mt-4" onSubmit={handleAddMember}>
            <Input name="email" type="email" placeholder="Email" required />
            <Select name="role" required>
                <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder="Select role" />
                </SelectTrigger>
                <SelectContent>
                    <SelectGroup>
                        <SelectLabel>Role</SelectLabel>
                        {roles.map((role) => (
                            <SelectItem key={role.value} value={role.value}>
                                {role.label}
                            </SelectItem>
                        ))}
                    </SelectGroup>
                </SelectContent>
            </Select>
            <Button type="submit" variant="default">Submit</Button>
        </form>
    );
}
