import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL as string;

export async function DELETE(req: NextRequest) {
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

  const { task_id } = body || {};
  if (!task_id) {
    return NextResponse.json({ error: "Missing task_id" }, { status: 400 });
  }

  // Forward DELETE request to backend API
  const backendRes = await fetch(`${API_BASE_URL}/forms/project/task/delete`, {
    method: "DELETE",
    headers: {
      "Authorization": `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ task_id }),
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
      { error: "Failed to delete task", details: result },
      { status: backendRes.status }
    );
  }

  return NextResponse.json(result);
}
