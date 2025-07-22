'use server'


const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL as string;


export async function handleRegister(formData: FormData) {
  const name = formData.get("name");
  const identity = formData.get("identity");
  const password = formData.get("password");
  const passwordConfirm = formData.get("passwordConfirm");

  const res = await fetch(`${API_BASE_URL}/collections/users/records`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      name,
      email: identity,
      password,
      passwordConfirm
    }),
  });

  const data = await res.json();

  if (!res.ok) {
    return { error: data.message || "Registration failed" };
  }

  return { success: true };
}
