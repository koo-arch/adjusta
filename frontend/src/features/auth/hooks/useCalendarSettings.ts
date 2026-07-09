'use client'
import { useQuery } from '@tanstack/react-query';
import { fetchCalendarSettings } from '@/features/auth/api/fetchCalendarSettings';
import { buildCalendarSettingsQueryKey } from '@/features/auth/queryKeys';

export const useCalendarSettings = () => {
    const { data, isLoading, error } = useQuery({
        queryKey: buildCalendarSettingsQueryKey(),
        queryFn: fetchCalendarSettings,
    });

    return {
        calendarSettings: data,
        isLoading,
        error,
    };
};
