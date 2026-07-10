'use client'
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-toastify';
import { updateCalendarSetting } from '@/features/auth/api/updateCalendarSetting';
import { buildCalendarSettingsQueryKey } from '@/features/auth/queryKeys';
import type { CalendarSetting, CalendarSettingUpdate } from '@/features/auth/types';

const applyUpdate = (
    settings: CalendarSetting[],
    id: string,
    payload: CalendarSettingUpdate,
): CalendarSetting[] =>
    settings.map((setting) => {
        if (setting.id === id) {
            return {
                ...setting,
                ...(payload.role !== undefined ? { role: payload.role } : {}),
                ...(payload.is_visible !== undefined ? { is_visible: payload.is_visible } : {}),
                ...(payload.sync_proposed_dates !== undefined
                    ? { sync_proposed_dates: payload.sync_proposed_dates }
                    : {}),
            };
        }
        // 登録先(primary)の付け替え時はサーバ側で旧 primary が reference に降格される
        if (payload.role === 'primary' && setting.role === 'primary') {
            return { ...setting, role: 'reference' as const };
        }
        return setting;
    });

export const useUpdateCalendarSetting = () => {
    const queryClient = useQueryClient();
    const queryKey = buildCalendarSettingsQueryKey();

    const mutation = useMutation({
        mutationFn: ({ id, payload }: { id: string; payload: CalendarSettingUpdate }) =>
            updateCalendarSetting(id, payload),
        // トグル・付け替えは楽観更新とし、失敗時はロールバック + エラートースト(screen-design 5.8)
        onMutate: async ({ id, payload }) => {
            await queryClient.cancelQueries({ queryKey });
            const previous = queryClient.getQueryData<CalendarSetting[]>(queryKey);
            if (previous) {
                queryClient.setQueryData(queryKey, applyUpdate(previous, id, payload));
            }
            return { previous };
        },
        onError: (_error, _variables, context) => {
            if (context?.previous) {
                queryClient.setQueryData(queryKey, context.previous);
            }
            toast.error('カレンダー設定の更新に失敗しました');
        },
        onSettled: async () => {
            await queryClient.invalidateQueries({ queryKey });
        },
    });

    return {
        update: mutation.mutate,
        isPending: mutation.isPending,
        error: mutation.error,
    };
};
