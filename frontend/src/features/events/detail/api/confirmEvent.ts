import { apiClient } from '@/lib/api/client';
import type { ConfirmForm } from '@/features/events/detail/schema';

export const confirmEvent = async (eventID: string, payload: ConfirmForm) => {
    await apiClient.patch<void, ConfirmForm>(`/api/calendar/event/confirm/${eventID}`, payload);
};
