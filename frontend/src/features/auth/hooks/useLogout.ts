'use client'
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { logout as requestLogout } from '@/features/auth/api/logout';
import { buildCurrentUserQueryKey } from '@/features/auth/queryKeys';

export const useLogout = () => {
    const queryClient = useQueryClient();
    const router = useRouter();

    const mutation = useMutation({
        mutationFn: requestLogout,
        onSuccess: (result) => {
            if (!result.ok) {
                if (result.type === 'request') {
                    console.error(result.errors.formErrors[0] ?? 'ログアウトに失敗しました。');
                    return;
                }

                console.error('ログアウトに失敗しました。');
                return;
            }

            queryClient.clear();
            queryClient.setQueryData(buildCurrentUserQueryKey(), null);
            router.push('/login');
        },
        onError: (error) => {
            console.error(error);
        },
    });

    const logout = async () => {
        const result = await mutation.mutateAsync();
        return result.ok;
    };

    return {
        logout,
        isPending: mutation.isPending,
    };
};
