'use client'
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-toastify';
import { deleteDraftEvent } from '@/features/events/detail/api/deleteDraftEvent';
import {
    buildDraftEventListQueryKey,
    buildEventDetailQueryKey,
    buildNeedsActionDraftsQueryKey,
    buildUpcomingEventsQueryKey,
} from '@/features/events/queryKeys';

export const useDeleteDraftMutation = (eventID: string) => {
    const queryClient = useQueryClient();

    const mutation = useMutation({
        mutationFn: async () => {
            await deleteDraftEvent(eventID);
        },
        onSuccess: async () => {
            await Promise.all([
                queryClient.invalidateQueries({ queryKey: buildEventDetailQueryKey(eventID) }),
                queryClient.invalidateQueries({ queryKey: buildDraftEventListQueryKey() }),
                queryClient.invalidateQueries({ queryKey: buildNeedsActionDraftsQueryKey() }),
                queryClient.invalidateQueries({ queryKey: buildUpcomingEventsQueryKey() }),
            ]);
            toast.success('イベントを削除しました');
        },
        onError: () => {
            toast.error('イベントの削除に失敗しました。時間をおいて再度お試しください。');
        },
    });

    const submit = async (): Promise<boolean> => {
        try {
            await mutation.mutateAsync();
            return true;
        } catch {
            return false;
        }
    };

    return {
        submit,
        isPending: mutation.isPending,
    };
};
