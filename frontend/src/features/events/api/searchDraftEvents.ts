import { apiClient } from '@/lib/api/client';
import type { EventDraftDetail, SearchParams } from '@/features/events/types';

export const searchDraftEvents = async (params: SearchParams): Promise<EventDraftDetail[]> => {
    const response = await apiClient.get<EventDraftDetail[]>('/api/event/draft/search', {
        query: params,
    });
    return response.data;
};
