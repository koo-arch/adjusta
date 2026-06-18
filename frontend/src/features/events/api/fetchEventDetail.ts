import { apiClient } from '@/lib/api/client';
import type { EventDraftDetail } from '@/features/events/types';

export const fetchEventDetail = async (eventID: string) => {
    const response = await apiClient.get<EventDraftDetail>(`/api/calendar/event/draft/${eventID}`);
    return response.data;
};
