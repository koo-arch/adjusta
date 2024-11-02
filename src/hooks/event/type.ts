import type { ProposedDate } from '@/atoms/calendar';

export interface EventDraftDetail {
    id: string;
    title: string;
    description: string;
    allDay: boolean;
    location: string;
    url: string;
    status: string;
    confirmed_date_id: string | null;
    proposed_dates: ProposedDate[];
}