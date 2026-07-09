import { apiClient } from '@/lib/api/client';
import type { CalendarSetting, CalendarSettingUpdate } from '@/features/auth/types';

export const updateCalendarSetting = async (id: string, payload: CalendarSettingUpdate) => {
    const response = await apiClient.patch<CalendarSetting, CalendarSettingUpdate>(
        `/api/user-calendars/${id}`,
        payload,
    );
    return response.data;
};
