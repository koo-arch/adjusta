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
        mutationFn: async () => deleteDraftEvent(eventID),
        onSuccess: async (result) => {
            if (!result.ok) {
                if (result.type === 'request') {
                    toast.error(result.errors.formErrors[0] ?? 'イベントの削除に失敗しました。時間をおいて再度お試しください。');
                    return;
                }

                toast.error('イベントの削除に失敗しました。時間をおいて再度お試しください。');
                return;
            }

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
        const result = await mutation.mutateAsync();
        return result.ok;
    };

    return {
        submit,
        isPending: mutation.isPending,
    };
};
