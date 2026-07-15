'use client'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { fetchCandidateSyncSetting } from '@/features/auth/api/fetchCandidateSyncSetting';
import { updateCandidateSyncSetting } from '@/features/auth/api/updateCandidateSyncSetting';
import {
    buildCalendarSettingsQueryKey,
    buildCandidateSyncSettingQueryKey,
} from '@/features/auth/queryKeys';
import type { CandidateSyncSetting } from '@/features/auth/types';

export const useCandidateSyncSetting = () => {
    const queryClient = useQueryClient();
    const queryKey = buildCandidateSyncSettingQueryKey();
    const query = useQuery({ queryKey, queryFn: fetchCandidateSyncSetting });
    const mutation = useMutation({
        mutationFn: updateCandidateSyncSetting,
        onMutate: async (enabled) => {
            await queryClient.cancelQueries({ queryKey });
            const previous = queryClient.getQueryData<CandidateSyncSetting>(queryKey);
            queryClient.setQueryData<CandidateSyncSetting>(queryKey, (current) => ({
                enabled,
                calendar: current?.calendar ?? null,
            }));
            return { previous };
        },
        onError: (_error, _enabled, context) => {
            if (context?.previous) queryClient.setQueryData(queryKey, context.previous);
            toast.error('候補日程の同期設定を更新できませんでした');
        },
        onSuccess: (setting) => queryClient.setQueryData(queryKey, setting),
        onSettled: async () => {
            await Promise.all([
                queryClient.invalidateQueries({ queryKey }),
                queryClient.invalidateQueries({ queryKey: buildCalendarSettingsQueryKey() }),
            ]);
        },
    });

    return {
        setting: query.data,
        isLoading: query.isLoading,
        error: query.error,
        setEnabled: mutation.mutate,
        isUpdating: mutation.isPending,
    };
};
