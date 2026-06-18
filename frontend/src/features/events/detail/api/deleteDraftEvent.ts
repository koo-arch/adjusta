import { apiClient } from '@/lib/api/client';

export const deleteDraftEvent = async (eventID: string) => {
    await apiClient.delete<void>(`/api/calendar/event/draft/${eventID}`);
};
