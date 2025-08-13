"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@radix-ui/react-scroll-area";

const baseUrl = process.env.NEXT_PUBLIC_BASE_URL || "http://localhost:3000";

interface MailItem {
  id: string;
  content: string;
  created: string;
}

export function TaskMailsList({ taskId }: { taskId: string }) {
  const [mails, setMails] = useState<MailItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchMails() {
      try {
        setLoading(true);
        setError(null);
        const res = await fetch(`${baseUrl}/api/tasks/mails?taskId=${taskId}`, {
          method: "GET",
          credentials: "include",
        });
        if (!res.ok) throw new Error("Failed to fetch mails");
        const data = await res.json();
        setMails(data.emails || []);
      } catch (err: any) {
        setError(err.message || "Failed to load mails");
      } finally {
        setLoading(false);
      }
    }
    if (taskId) fetchMails();
  }, [taskId]);

  if (loading) return <p>Loading mails...</p>;
  if (error) return <p className="text-red-500">{error}</p>;

  if (mails.length === 0) {
    return <p className="text-muted-foreground">No mails linked to this task.</p>;
  }

  return (
    <ScrollArea className="space-y-4 m-4">
      {mails.map((mail) => (
        <Card key={mail.id}>
          <CardHeader>
            <CardDescription>
              {new Date(mail.created).toLocaleString()}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {mail.content}
          </CardContent>
        </Card>
      ))}
    </ScrollArea>
  );
}
