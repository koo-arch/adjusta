import { apiClient } from '@/lib/api/client';
import type { AuthUser } from '@/features/auth/types';

export const fetchAccount = async () => {
    const response = await apiClient.get<AuthUser>('/api/account/list');
    return response.data;
};
