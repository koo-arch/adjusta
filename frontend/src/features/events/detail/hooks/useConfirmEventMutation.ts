'use client'
import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-toastify';
import type { ConfirmForm, ConfirmFormErrors } from '@/features/events/detail/schema';
import { buildFormErrorsFromAPIError, emptyFormErrors, type FormErrors } from '@/lib/form/errors';
import { confirmEvent } from '@/features/events/detail/api/confirmEvent';
import {
    buildDraftEventListQueryKey,
    buildEventDetailQueryKey,
    buildNeedsActionDraftsQueryKey,
    buildUpcomingEventsQueryKey,
} from '@/features/events/queryKeys';

type ConfirmMutationFieldKey = keyof ConfirmFormErrors;

const createEmptyConfirmMutationErrors = (): FormErrors<ConfirmMutationFieldKey> =>
    emptyFormErrors<ConfirmMutationFieldKey>();

export const useConfirmEventMutation = (eventID: string) => {
    const queryClient = useQueryClient();
    const [errors, setErrors] = useState<FormErrors<ConfirmMutationFieldKey>>(createEmptyConfirmMutationErrors);
    const [confirmed, setConfirmed] = useState(false);

    const mutation = useMutation({
        mutationFn: async (payload: ConfirmForm) => {
            await confirmEvent(eventID, payload);
        },
        onMutate: () => {
            setErrors(createEmptyConfirmMutationErrors());
            setConfirmed(false);
        },
        onSuccess: async () => {
            await Promise.all([
                queryClient.invalidateQueries({ queryKey: buildEventDetailQueryKey(eventID) }),
                queryClient.invalidateQueries({ queryKey: buildDraftEventListQueryKey() }),
                queryClient.invalidateQueries({ queryKey: buildNeedsActionDraftsQueryKey() }),
                queryClient.invalidateQueries({ queryKey: buildUpcomingEventsQueryKey() }),
            ]);
            setConfirmed(true);
            toast.success('日程を確定しました');
        },
        onError: (error) => {
            setErrors(buildFormErrorsFromAPIError<ConfirmMutationFieldKey>(error, '日程の確定に失敗しました。時間をおいて再度お試しください。'));
            setConfirmed(false);
        },
    });

    const submit = async (draft: ConfirmForm): Promise<boolean> => {
        try {
            await mutation.mutateAsync(draft);
            return true;
        } catch {
            return false;
        }
    };

    const reset = () => {
        setErrors(createEmptyConfirmMutationErrors());
        setConfirmed(false);
        mutation.reset();
    };

    return {
        errors,
        confirmed,
        submit,
        reset,
        isPending: mutation.isPending,
    };
};
