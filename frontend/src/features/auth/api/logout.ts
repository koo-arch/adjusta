import { APIClientError, apiClient } from '@/lib/api/client';
import { buildFormErrorsFromAPIError } from '@/lib/form/errors';
import type { SubmitResult } from '@/lib/form/submit';

export const logout = async (): Promise<SubmitResult<null, never>> => {
    try {
        await apiClient.get<void>('/api/auth/logout');

        return {
            ok: true,
            data: null,
        };
    } catch (error) {
        if (!(error instanceof APIClientError)) {
            throw error;
        }

        return {
            ok: false,
            type: 'request',
            errors: buildFormErrorsFromAPIError<never>(
                error,
                'ログアウトに失敗しました。時間をおいて再度お試しください。',
            ),
        };
    }
};
