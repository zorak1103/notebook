import { useState } from 'react';
import { fetchWhoAmI } from '../api/client';
import type { UserInfo } from '../api/types';

interface UseWhoAmIResult {
  user: UserInfo | null;
  loading: boolean;
  error: string | null;
  fetchUser: () => Promise<void>;
}

/**
 * Custom hook for managing WhoAmI API state
 * @returns Object containing user data, loading state, error, and fetch function
 */
export function useWhoAmI(): UseWhoAmIResult {
  const [user, setUser] = useState<UserInfo | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchUser = async () => {
    setLoading(true);
    setError(null);

    try {
      const userData = await fetchWhoAmI();
      setUser(userData);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error occurred';
      setError(message);
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  return {
    user,
    loading,
    error,
    fetchUser,
  };
}
