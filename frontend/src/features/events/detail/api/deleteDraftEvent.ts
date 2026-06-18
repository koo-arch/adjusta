import { APIClientError, apiClient } from '@/lib/api/client';
import { buildFormErrorsFromAPIError } from '@/lib/form/errors';
import type { SubmitResult } from '@/lib/form/submit';

export const deleteDraftEvent = async (
    eventID: string,
): Promise<SubmitResult<null, never>> => {
    try {
        await apiClient.delete<void>(`/api/calendar/event/draft/${eventID}`);

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
                'イベントの削除に失敗しました。時間をおいて再度お試しください。',
            ),
        };
    }
};
