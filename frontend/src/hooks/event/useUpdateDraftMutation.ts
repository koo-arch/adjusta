'use client'
import { useState } from 'react';
import { useSetAtom } from 'jotai';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { clearEventFormErrorStateAtomFamily, setServerEventFormErrorsAtomFamily } from '@/features/events/form-meta-atoms';
import { toast } from 'react-toastify';
import type { EventFormErrors, EventUpdateForm } from '@/features/events/zod';
import { apiClient } from '@/lib/api/client';
import { buildFormErrorsFromAPIError } from '@/lib/form/errors';
import { buildEventDetailQueryKey } from './useFetchEventDetail';
import { buildDraftEventListQueryKey } from './useFetchEventList';
import { buildNeedsActionDraftsQueryKey } from './useFetchNeedsActionDrafts';
import { buildUpcomingEventsQueryKey } from './useFetchUpcomingEvents';

export const useUpdateDraftMutation = (eventID: string) => {
    const queryClient = useQueryClient();
    const [updated, setUpdated] = useState(false);
    const clearErrorState = useSetAtom(clearEventFormErrorStateAtomFamily(eventID));
    const setServerErrors = useSetAtom(setServerEventFormErrorsAtomFamily(eventID));

    const mutation = useMutation({
        mutationFn: async (payload: EventUpdateForm) => {
            await apiClient.put<void, EventUpdateForm>(`/api/calendar/event/draft/${eventID}`, payload);
        },
        onMutate: () => {
            clearErrorState();
            setUpdated(false);
        },
        onSuccess: async () => {
            await Promise.all([
                queryClient.invalidateQueries({ queryKey: buildEventDetailQueryKey(eventID) }),
                queryClient.invalidateQueries({ queryKey: buildDraftEventListQueryKey() }),
                queryClient.invalidateQueries({ queryKey: buildNeedsActionDraftsQueryKey() }),
                queryClient.invalidateQueries({ queryKey: buildUpcomingEventsQueryKey() }),
            ]);
            clearErrorState();
            setUpdated(true);
            toast.success('イベントを更新しました');
        },
        onError: (error) => {
            setServerErrors(buildFormErrorsFromAPIError<keyof EventFormErrors>(error, 'イベントの更新に失敗しました。時間をおいて再度お試しください。'));
            setUpdated(false);
        },
    });

    const submit = async (draft: EventUpdateForm): Promise<boolean> => {
        try {
            await mutation.mutateAsync(draft);
            return true;
        } catch {
            return false;
        }
    };

    const reset = () => {
        clearErrorState();
        setUpdated(false);
        mutation.reset();
    };

    return {
        updated,
        submit,
        reset,
        isPending: mutation.isPending,
    };
};
