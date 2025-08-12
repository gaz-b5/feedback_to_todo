'use server'


const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL as string;


export async function handleRegister(formData: FormData) {
  const name = formData.get("name");
  const identity = formData.get("identity");
  const password = formData.get("password");
  const passwordConfirm = formData.get("passwordConfirm");

  try {

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

    const verifyRes = await fetch(`${API_BASE_URL}/collections/users/request-verification`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: identity }),
    });

    if (!verifyRes.ok) {
      const verifyData = await verifyRes.json();
      return { error: verifyData.message || "Failed to send verification email" };
    }

    return { success: true };

  } catch (err) {
    console.error("Registration error:", err);
    return { error: "Something went wrong during registration." };
  }
}
