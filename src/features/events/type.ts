import type { SendSelectedDate, SendProposedDate } from '@/atoms/calendar';

export interface EventDraftForm {
    title: string;
    description: string;
    allDay: boolean;
    location: string;
    url: string;
    selected_dates: SendSelectedDate[];
}

export interface EventUpdateForm {
    id: string;
    title?: string;
    description?: string;
    allDay?: boolean;
    location?: string;
    url?: string;
    proposed_dates?: SendProposedDate[];
}