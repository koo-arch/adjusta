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
    google_event_id: string;
    slug: string;
    proposed_dates: ProposedDate[];
}

export interface SearchParams {
    title?: string;
    location?: string;
    startTime?: string;
    endTime?: string;
    status?: "confirmed" | "pending" | "rejected";
}

export interface UpcomingEvent extends Omit<EventDraftDetail, 'proposed_dates'> {
    start: Date;
    end: Date;
}

export interface NeedsActionDraft extends Omit<UpcomingEvent, 'confirmed_date_id' | 'google_event_id'> {
    needs_attention: boolean;
}