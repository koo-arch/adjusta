import { apiClient } from '@/lib/api/client';
import type { GoogleCalendarResponse } from '@/features/calendar/types';

export const fetchGoogleCalendarEvents = async () => {
    const response = await apiClient.get<GoogleCalendarResponse>('/api/calendar/list');
    return response.data;
};
