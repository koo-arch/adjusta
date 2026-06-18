'use client'
import { useState } from 'react';
import { useSetAtom } from 'jotai';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { clearEventFormErrorStateAtomFamily, setServerEventFormErrorsAtomFamily } from '@/features/events/store/errors';
import { toast } from 'react-toastify';
import type { EventFormErrors, EventDraftForm } from '@/features/events/schema';
import { buildFormErrorsFromAPIError } from '@/lib/form/errors';
import { createDraftEvent } from '@/features/events/draft/api/createDraftEvent';
import {
    buildDraftEventListQueryKey,
    buildNeedsActionDraftsQueryKey,
} from '@/features/events/queryKeys';

export const useCreateDraftMutation = (formScope: string) => {
    const queryClient = useQueryClient();
    const [createdDraftID, setCreatedDraftID] = useState<string | null>(null);
    const clearErrorState = useSetAtom(clearEventFormErrorStateAtomFamily(formScope));
    const setServerErrors = useSetAtom(setServerEventFormErrorsAtomFamily(formScope));

    const mutation = useMutation({
        mutationFn: createDraftEvent,
        onMutate: () => {
            clearErrorState();
            setCreatedDraftID(null);
        },
        onSuccess: async (result) => {
            await Promise.all([
                queryClient.invalidateQueries({ queryKey: buildDraftEventListQueryKey() }),
                queryClient.invalidateQueries({ queryKey: buildNeedsActionDraftsQueryKey() }),
            ]);
            clearErrorState();
            setCreatedDraftID(result.id);
            toast.success('イベントを作成しました');
        },
        onError: (error) => {
            setServerErrors(buildFormErrorsFromAPIError<keyof EventFormErrors>(error, 'イベントの作成に失敗しました。時間をおいて再度お試しください。'));
            setCreatedDraftID(null);
        },
    });

    const submit = async (draft: EventDraftForm): Promise<string | null> => {
        try {
            const result = await mutation.mutateAsync(draft);
            return result.id;
        } catch {
            return null;
        }
    };

    const reset = () => {
        clearErrorState();
        setCreatedDraftID(null);
        mutation.reset();
    };

    return {
        createdDraftID,
        submit,
        reset,
        isPending: mutation.isPending,
    };
};
