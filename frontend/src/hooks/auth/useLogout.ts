'use client'
import { useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api/client';
import { currentUserKey } from './useAuth';

export const useLogout = () => {
    const queryClient = useQueryClient();
    const router = useRouter();

    const logout = () => {
        apiClient.get('/api/auth/logout')
            .then(() => {
                queryClient.clear();
                queryClient.setQueryData(currentUserKey, null);
                router.push('/login');
            })
            .catch(err => console.error(err));
    }

    return logout 
}
