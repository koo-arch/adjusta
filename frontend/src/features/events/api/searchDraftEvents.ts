import { apiClient } from '@/lib/api/client';
import type { EventDraftListResponse, SearchParams } from '@/features/events/types';

export const searchDraftEvents = async (params: SearchParams): Promise<EventDraftListResponse> => {
    const response = await apiClient.get<EventDraftListResponse>('/api/event/draft/search', {
        query: params,
    });
    return response.data;
};
