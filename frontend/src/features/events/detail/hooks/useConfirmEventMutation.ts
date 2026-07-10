'use client'
import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-toastify';
import type { ConfirmFormErrors } from '@/features/events/detail/schema';
import { emptyFormErrors, type FormErrors } from '@/lib/form/errors';
import { confirmEvent, type ConfirmEventInput } from '@/features/events/detail/api/confirmEvent';
import {
    buildDraftEventSearchQueryKey,
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

    const mutation = useMutation({
        mutationFn: async (payload: ConfirmEventInput) => confirmEvent(eventID, payload),
        onMutate: () => {
            setErrors(createEmptyConfirmMutationErrors());
        },
        onSuccess: async (result) => {
            if (!result.ok) {
                if (result.type === 'validation') {
                    setErrors({
                        formErrors: [],
                        fieldErrors: result.errors,
                    });
                    return;
                }

                setErrors(result.errors);
                return;
            }

            await Promise.all([
                queryClient.invalidateQueries({ queryKey: buildEventDetailQueryKey(eventID) }),
                queryClient.invalidateQueries({ queryKey: buildDraftEventSearchQueryKey() }),
                queryClient.invalidateQueries({ queryKey: buildNeedsActionDraftsQueryKey() }),
                queryClient.invalidateQueries({ queryKey: buildUpcomingEventsQueryKey() }),
            ]);
            toast.success('日程を確定しました');
        },
        onError: () => {
            setErrors({
                formErrors: ['日程の確定に失敗しました。時間をおいて再度お試しください。'],
                fieldErrors: {},
            });
        },
    });

    const submit = async (draft: ConfirmEventInput): Promise<boolean> => {
        const result = await mutation.mutateAsync(draft);
        return result.ok;
    };

    const reset = () => {
        setErrors(createEmptyConfirmMutationErrors());
        mutation.reset();
    };

    return {
        errors,
        submit,
        reset,
        isPending: mutation.isPending,
    };
};
