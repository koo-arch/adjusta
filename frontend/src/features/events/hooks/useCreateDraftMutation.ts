'use client'
import { useSetAtom } from 'jotai';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
    clearEventFormErrorStateAtomFamily,
    setClientEventFormErrorsAtomFamily,
    setServerEventFormErrorsAtomFamily,
} from '@/features/events/store/errors';
import { toast } from 'react-toastify';
import type { EventDraftForm } from '@/features/events/schema';
import { createDraftEvent } from '@/features/events/draft/api/createDraftEvent';
import {
    buildDraftEventSearchQueryKey,
    buildNeedsActionDraftsQueryKey,
} from '@/features/events/queryKeys';

export const useCreateDraftMutation = (formScope: string) => {
    const queryClient = useQueryClient();
    const clearErrorState = useSetAtom(clearEventFormErrorStateAtomFamily(formScope));
    const setClientErrors = useSetAtom(setClientEventFormErrorsAtomFamily(formScope));
    const setServerErrors = useSetAtom(setServerEventFormErrorsAtomFamily(formScope));

    const mutation = useMutation({
        mutationFn: createDraftEvent,
        onMutate: () => {
            clearErrorState();
        },
        onSuccess: async (result) => {
            if (!result.ok) {
                if (result.type === 'validation') {
                    setClientErrors(result.errors);
                    return;
                }

                setServerErrors(result.errors);
                return;
            }

            await Promise.all([
                queryClient.invalidateQueries({ queryKey: buildDraftEventSearchQueryKey() }),
                queryClient.invalidateQueries({ queryKey: buildNeedsActionDraftsQueryKey() }),
            ]);
            clearErrorState();
            toast.success('イベントを作成しました');
        },
        onError: () => {
            setServerErrors({
                formErrors: ['イベントの作成に失敗しました。時間をおいて再度お試しください。'],
                fieldErrors: {},
            });
        },
    });

    const submit = async (draft: EventDraftForm): Promise<string | null> => {
        const result = await mutation.mutateAsync(draft);
        return result.ok ? result.data.id : null;
    };

    const reset = () => {
        clearErrorState();
        mutation.reset();
    };

    return {
        submit,
        reset,
        isPending: mutation.isPending,
    };
};
