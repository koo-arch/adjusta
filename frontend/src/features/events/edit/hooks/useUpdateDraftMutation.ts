'use client'
import { useSetAtom } from 'jotai';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
    clearEventFormErrorStateAtomFamily,
    setClientEventFormErrorsAtomFamily,
    setServerEventFormErrorsAtomFamily,
} from '@/features/events/store/errors';
import { toast } from 'sonner';
import type { EventUpdateForm } from '@/features/events/schema';
import { updateDraftEvent } from '@/features/events/edit/api/updateDraftEvent';
import {
    buildDraftEventSearchQueryKey,
    buildEventDetailQueryKey,
    buildNeedsActionDraftsQueryKey,
    buildUpcomingEventsQueryKey,
} from '@/features/events/queryKeys';

export const useUpdateDraftMutation = (eventID: string) => {
    const queryClient = useQueryClient();
    const clearErrorState = useSetAtom(clearEventFormErrorStateAtomFamily(eventID));
    const setClientErrors = useSetAtom(setClientEventFormErrorsAtomFamily(eventID));
    const setServerErrors = useSetAtom(setServerEventFormErrorsAtomFamily(eventID));

    const mutation = useMutation({
        mutationFn: async (payload: EventUpdateForm) => updateDraftEvent(eventID, payload),
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
                queryClient.invalidateQueries({ queryKey: buildEventDetailQueryKey(eventID) }),
                queryClient.invalidateQueries({ queryKey: buildDraftEventSearchQueryKey() }),
                queryClient.invalidateQueries({ queryKey: buildNeedsActionDraftsQueryKey() }),
                queryClient.invalidateQueries({ queryKey: buildUpcomingEventsQueryKey() }),
            ]);
            clearErrorState();
            toast.success('イベントを更新しました');
        },
        onError: () => {
            setServerErrors({
                formErrors: ['イベントの更新に失敗しました。時間をおいて再度お試しください。'],
                fieldErrors: {},
            });
        },
    });

    const submit = async (draft: EventUpdateForm): Promise<boolean> => {
        const result = await mutation.mutateAsync(draft);
        return result.ok;
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
