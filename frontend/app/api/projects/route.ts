import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL as string;

export async function GET(req: NextRequest) {
  // Extract token from httpOnly cookie securely on server
  const token = req.cookies.get("session_token")?.value;

  if (!token) {
    return NextResponse.json({ error: "Unauthorized" }, { status: 401 });
  }

  // Call your backend securely with token in the Authorization header
  const backendRes = await fetch(`${API_BASE_URL}/forms/projects`, {
    method: "GET",
    headers: {
      "Authorization": `Bearer ${token}`,
    },
  });

  if (!backendRes.ok) {
    return NextResponse.json({ error: "Failed to fetch projects"}, { status: backendRes.status });
  }

  const data = await backendRes.json();

  // Transform projects array
  const transformedProjects = Array.isArray(data.projects)
    ? data.projects.map((p: any) => ({
        id: p.id,
        title: p.name,
        isActive: false,
      }))
    : [];

  return NextResponse.json({ projects: transformedProjects });
}
