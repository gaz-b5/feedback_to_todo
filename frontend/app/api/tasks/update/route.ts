import { assert } from "console";
import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL as string;

export async function PATCH(req: NextRequest) {
  // Get token from httpOnly cookie
  const token = req.cookies.get("session_token")?.value;

  if (!token) {
    return NextResponse.json({ error: "Unauthorized" }, { status: 401 });
  }

  let body: any;
  try {
    body = await req.json();
  } catch (err) {
    return NextResponse.json({ error: "Invalid JSON" }, { status: 400 });
  }

  const { task_id, status, priority, nature, assigned } = body || {};
  if (!task_id) {
    return NextResponse.json({ error: "Missing task_id" }, { status: 400 });
  }

  // Build the PATCH body with only defined fields
  const patchBody: any = { task_id };
  if (status !== undefined) patchBody.status = status;
  if (priority !== undefined) patchBody.priority = priority;
  if (nature !== undefined) patchBody.nature = nature;
  if (assigned !== undefined) patchBody.assigned = assigned;

  // Forward to backend (endpoint should accept PATCH for updating a single task)
  const backendRes = await fetch(`${API_BASE_URL}/forms/project/task/update`, {
    method: "PATCH",
    headers: {
      "Authorization": `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(patchBody),
  });

  const text = await backendRes.text();
  let result;
  try {
    result = JSON.parse(text);
  } catch {
    result = text;
  }

  if (!backendRes.ok) {
    return NextResponse.json(
      { error: "Failed to update task", details: result },
      { status: backendRes.status }
    );
  }

  return NextResponse.json(result);
}
