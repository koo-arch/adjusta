import { apiClient } from '@/lib/api/client';
import type { EventDraftListResponse, EventListParams } from '@/features/events/types';

export const fetchDraftEventList = async (params: EventListParams = {}) => {
    const response = await apiClient.get<EventDraftListResponse>('/api/calendar/event/draft/list', {
        query: params,
    });
    return response.data;
};
