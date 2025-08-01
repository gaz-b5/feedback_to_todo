import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL as string;

export async function GET(req: NextRequest) {
    // Extract token from httpOnly cookie securely on server
    const token = req.cookies.get("session_token")?.value;

    if (!token) {
        return NextResponse.json({ error: "Unauthorized" }, { status: 401 });
    }

    // Get projectId from query params
    const { searchParams } = new URL(req.url);
    const projectId = searchParams.get("projectId");

    // Call your backend securely with token in the Authorization header
    if (!projectId) {
        return NextResponse.json({ error: "Missing projectId" }, { status: 400 });
    }

    // Call your backend securely with token in the Authorization header
    const backendRes = await fetch(`${API_BASE_URL}/forms/project/task/members?projectId=${projectId}`, {
        method: "GET",
        headers: {
            "Authorization": `Bearer ${token}`,
        },
    });

    if (!backendRes.ok) {
        return NextResponse.json({ error: "Failed to fetch tasks" }, { status: backendRes.status });
    }

    const members = await backendRes.json();

    console.log("Fetched members:", members);

    return NextResponse.json(members.members || []);
}