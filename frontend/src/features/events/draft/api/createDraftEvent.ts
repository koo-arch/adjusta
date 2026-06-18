import { apiClient } from '@/lib/api/client';
import type { EventDraftForm } from '@/features/events/schema';

export const createDraftEvent = async (payload: EventDraftForm) => {
    const response = await apiClient.post<{ id: string }, EventDraftForm>('/api/calendar/event/draft', payload);
    return response.data;
};
