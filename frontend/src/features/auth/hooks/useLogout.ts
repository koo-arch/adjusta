'use client'
import { useMutation } from '@tanstack/react-query';
import { logout as requestLogout } from '@/features/auth/api/logout';

export const useLogout = () => {
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

            // フルリロードで query cache / Jotai atom の認証済みデータを確実に破棄する(401経路と同じ理屈)
            window.location.assign('/login');
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
