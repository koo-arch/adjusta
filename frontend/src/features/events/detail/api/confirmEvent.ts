import { APIClientError, apiClient } from '@/lib/api/client';
import { buildFormErrorsFromAPIError } from '@/lib/form/errors';
import type { SubmitResult } from '@/lib/form/submit';
import { buildZodFieldErrors } from '@/lib/validation/zod';
import {
    ConfirmFormSchema,
    type ConfirmForm,
    type ConfirmFormErrors,
} from '@/features/events/detail/schema';

export interface ConfirmEventInput {
    confirm_date: {
        id: string | null;
        google_event_id?: string;
        start: Date | null;
        end: Date | null;
        priority: number;
    };
    selectionMode: 'dropdown' | 'manual';
}

type ConfirmEventFieldKey = keyof ConfirmFormErrors;

export const confirmEvent = async (
    eventID: string,
    payload: ConfirmEventInput,
): Promise<SubmitResult<null, ConfirmEventFieldKey>> => {
    if (payload.selectionMode === 'dropdown' && !payload.confirm_date.id) {
        return {
            ok: false,
            type: 'validation',
            errors: {
                confirm_date: '日程を選択してください',
            },
        };
    }

    const candidate = {
        confirm_date: {
            id: payload.confirm_date.id,
            google_event_id: payload.confirm_date.google_event_id,
            start: payload.confirm_date.start,
            end: payload.confirm_date.end,
            priority: payload.confirm_date.priority,
        },
    };

    const validated = ConfirmFormSchema.safeParse(candidate);
    if (!validated.success) {
        return {
            ok: false,
            type: 'validation',
            errors: buildZodFieldErrors<ConfirmEventFieldKey>(validated.error),
        };
    }

    try {
        await apiClient.patch<void, ConfirmForm>(
            `/api/calendar/event/confirm/${eventID}`,
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
            errors: buildFormErrorsFromAPIError<ConfirmEventFieldKey>(
                error,
                '日程の確定に失敗しました。時間をおいて再度お試しください。',
            ),
        };
    }
};
