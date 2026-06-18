import { apiClient } from '@/lib/api/client';
import type { EventDraftDetail } from '@/features/events/types';

export const fetchDraftEventList = async () => {
    const response = await apiClient.get<EventDraftDetail[]>('/api/calendar/event/draft/list');
    return response.data;
};
