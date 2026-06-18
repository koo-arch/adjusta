import { APIClientError, apiClient } from '@/lib/api/client';
import { buildFormErrorsFromAPIError } from '@/lib/form/errors';
import type { SubmitResult } from '@/lib/form/submit';
import { buildZodFieldErrors } from '@/lib/validation/zod';
import {
    EventUpdateFormSchema,
    type EventFormErrors,
    type EventUpdateForm,
} from '@/features/events/schema';

type UpdateDraftEventFieldKey = keyof EventFormErrors;

export const updateDraftEvent = async (
    eventID: string,
    payload: EventUpdateForm,
): Promise<SubmitResult<null, UpdateDraftEventFieldKey>> => {
    const validated = EventUpdateFormSchema.safeParse(payload);
    if (!validated.success) {
        return {
            ok: false,
            type: 'validation',
            errors: buildZodFieldErrors<UpdateDraftEventFieldKey>(validated.error),
        };
    }

    try {
        await apiClient.put<void, EventUpdateForm>(
            `/api/calendar/event/draft/${eventID}`,
            validated.data,
        );

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
            errors: buildFormErrorsFromAPIError<UpdateDraftEventFieldKey>(
                error,
                'イベントの更新に失敗しました。時間をおいて再度お試しください。',
            ),
        };
    }
};
