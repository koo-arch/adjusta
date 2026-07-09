import { apiClient } from '@/lib/api/client';
import type { CalendarSetting } from '@/features/auth/types';

export const fetchCalendarSettings = async () => {
    const response = await apiClient.get<CalendarSetting[]>('/api/user-calendars');
    return response.data;
};
