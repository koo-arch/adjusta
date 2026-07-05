import { APIClientError, apiClient } from '@/lib/api/client';
import { buildFormErrorsFromAPIError } from '@/lib/form/errors';
import type { SubmitResult } from '@/lib/form/submit';
import { buildZodFieldErrors } from '@/lib/validation/zod';
import {
    EventDraftFormSchema,
    type EventDraftForm,
    type EventFormErrors,
} from '@/features/events/schema';

type CreateDraftEventFieldKey = keyof EventFormErrors;

export const createDraftEvent = async (
    payload: EventDraftForm,
): Promise<SubmitResult<{ id: string }, CreateDraftEventFieldKey>> => {
    const validated = EventDraftFormSchema.safeParse(payload);
    if (!validated.success) {
        return {
            ok: false,
            type: 'validation',
            errors: buildZodFieldErrors<CreateDraftEventFieldKey>(validated.error),
        };
    }

    try {
        const response = await apiClient.post<{ id: string }, EventDraftForm>(
            '/api/calendar/event/draft',
            validated.data,
        );

        return {
            ok: true,
            data: response.data,
        };
    } catch (error) {
        if (!(error instanceof APIClientError)) {
            throw error;
        }

        return {
            ok: false,
            type: 'request',
            errors: buildFormErrorsFromAPIError<CreateDraftEventFieldKey>(
                error,
                'イベントの作成に失敗しました。時間をおいて再度お試しください。',
            ),
        };
    }
};
