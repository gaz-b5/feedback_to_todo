import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL as string;

export async function GET(req: NextRequest) {
    // Extract token from httpOnly cookie securely on server
    const token = req.cookies.get("session_token")?.value;

    if (!token) {
        return NextResponse.json({ error: "Unauthorized" }, { status: 401 });
    }

    // Get taskId from query params
    const { searchParams } = new URL(req.url);
    const taskId = searchParams.get("taskId");


    if (!taskId) {
        return NextResponse.json({ error: "Missing taskId" }, { status: 400 });
    }

    // Call your backend securely with token in the Authorization header
    const backendRes = await fetch(`${API_BASE_URL}/forms/project/task/mails?taskId=${taskId}`, {
        method: "GET",
        headers: {
            "Authorization": `Bearer ${token}`,
        },
    });

    if (!backendRes.ok) {
        return NextResponse.json({ error: "Failed to fetch mails" }, { status: backendRes.status });
    }

    const data = await backendRes.json();

    // const mails = Array.isArray(data.emails) ? data.emails : [];

    // console.log(data);

    return NextResponse.json(data);
}