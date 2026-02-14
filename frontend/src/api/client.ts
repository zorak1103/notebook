import type { UserInfo } from './types';

/**
 * Fetches the current user's Tailscale information from the backend
 * @returns Promise resolving to UserInfo
 * @throws Error if the request fails
 */
export async function fetchWhoAmI(): Promise<UserInfo> {
  const response = await fetch('/api/whoami');

  if (!response.ok) {
    throw new Error(`Failed to fetch user info: ${response.status} ${response.statusText}`);
  }

  return response.json();
}
