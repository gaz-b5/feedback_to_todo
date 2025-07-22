import 'server-only';

import { SignJWT, jwtVerify } from 'jose';

const secretKey = process.env.SECRET!;
const key = new TextEncoder().encode(secretKey);

export type SessionPayload = {
  token: string;
  expiresAt: number; // use timestamp in ms to avoid serialization issues
};

export async function encrypt(payload: SessionPayload) {
  return new SignJWT(payload)
    .setProtectedHeader({ alg: 'HS256' })
    .setIssuedAt()
    .setExpirationTime('7d')
    .sign(key);
}

export async function decrypt(token: string | undefined) {
  if (!token) return null;
  try {
    const { payload } = await jwtVerify(token, key, {
      algorithms: ['HS256'],
    });
    return payload as SessionPayload;
  } catch {
    return null;
  }
}


