import { apiClient } from '@/lib/api/client';
import type { EventDraftDetail, EventProposedDate } from '@/features/events/types';

type EventProposedDateResponse = Omit<EventProposedDate, 'start' | 'end'> & {
    start: string;
    end: string;
};

type EventDraftDetailResponse = Omit<EventDraftDetail, 'proposed_dates'> & {
    proposed_dates: EventProposedDateResponse[];
};

export const fetchEventDetail = async (eventID: string) => {
    const response = await apiClient.get<EventDraftDetailResponse>(`/api/calendar/event/draft/${eventID}`);

    return {
        ...response.data,
        proposed_dates: response.data.proposed_dates.map((date) => ({
            ...date,
            start: new Date(date.start),
            end: new Date(date.end),
        })),
    } satisfies EventDraftDetail;
};
