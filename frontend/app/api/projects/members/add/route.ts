import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL as string;

export async function POST(req: NextRequest) {
  // Extract token from httpOnly cookie securely on server
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
  console.log(body)


  // Call your backend securely with token in the Authorization header
  const backendRes = await fetch(`${API_BASE_URL}/forms/project/addUser`, {
    method: "POST",
    headers: {
      "Authorization": `Bearer ${token}`,
    },
    body: JSON.stringify(body),
  });

  if (!backendRes.ok) {
    return NextResponse.json({ error: "Failed to Add member" }, { status: backendRes.status });
  }

  const data = await backendRes.json();

  return NextResponse.json(data);
}