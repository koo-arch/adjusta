import { apiClient } from '@/lib/api/client';
import type { UpcomingEvent } from '@/features/events/types';

export const fetchUpcomingEvents = async () => {
    const response = await apiClient.get<UpcomingEvent[]>('/api/event/confirmed/upcoming');
    return response.data;
};
