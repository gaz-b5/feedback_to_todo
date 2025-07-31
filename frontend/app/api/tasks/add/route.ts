import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL as string;

export async function POST(req: NextRequest) {
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

  console.log("Received body:", JSON.stringify(body));


  const backendRes = await fetch(`${API_BASE_URL}/forms/project/task/add`, {
    method: "POST",
    headers: {
      "Authorization": `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
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
      { error: "Failed to add task", details: result },
      { status: backendRes.status }
    );
  }

  return NextResponse.json(result);
}
