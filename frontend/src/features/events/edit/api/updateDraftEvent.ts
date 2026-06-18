import { apiClient } from '@/lib/api/client';
import type { EventUpdateForm } from '@/features/events/schema';

export const updateDraftEvent = async (eventID: string, payload: EventUpdateForm) => {
    await apiClient.put<void, EventUpdateForm>(`/api/calendar/event/draft/${eventID}`, payload);
};
