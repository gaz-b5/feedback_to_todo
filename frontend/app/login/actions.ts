'use server';

import { NextResponse } from 'next/server';
import { encrypt, SessionPayload } from '@/app/_lib/session';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL as string;

export async function handleLogin(formData: FormData) {
  // Collect credentials from form data
  const identity = formData.get("identity");
  const password = formData.get("password");

  // Forward to PocketBase
  const res = await fetch(`${API_BASE_URL}/collections/users/auth-with-password`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ identity, password }),
    // credentials: 'include' // Not needed, runs server-side
  });

  const data = await res.json();

  if (!res.ok) {
    // You can throw error or return object for UI to use
    return { error: data.message || "Login failed" };
  }

  const token = data.token;
  if (!token) {
    return { error: "No token received" };
  }

  // Optionally return data for the client
  return { success: true };
}
