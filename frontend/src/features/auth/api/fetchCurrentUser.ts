import { apiClient } from '@/lib/api/client';
import type { AuthUser } from '@/features/auth/types';

export const fetchCurrentUser = async (): Promise<AuthUser | null> => {
    const response = await apiClient.get<AuthUser>('/api/users/me', {
        allowStatuses: [401],
    });

    if (response.status === 401) {
        return null;
    }

    return response.data ?? null;
};
