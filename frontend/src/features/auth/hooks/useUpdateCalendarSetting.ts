'use client'
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { updateCalendarSetting } from '@/features/auth/api/updateCalendarSetting';
import { buildCalendarSettingsQueryKey } from '@/features/auth/queryKeys';
import type { CalendarSettingUpdate } from '@/features/auth/types';

export const useUpdateCalendarSetting = () => {
    const queryClient = useQueryClient();

    const mutation = useMutation({
        mutationFn: ({ id, payload }: { id: string; payload: CalendarSettingUpdate }) =>
            updateCalendarSetting(id, payload),
        onSuccess: async () => {
            await queryClient.invalidateQueries({ queryKey: buildCalendarSettingsQueryKey() });
        },
    });

    return {
        update: mutation.mutateAsync,
        isPending: mutation.isPending,
        error: mutation.error,
    };
};
