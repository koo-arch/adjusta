import type { PrioritizedSelectedDate } from '@/atoms/calendar';

export interface EventDraftForm {
    title: string;
    description: string;
    allDay: boolean;
    location: string;
    url: string;
    selected_dates: PrioritizedSelectedDate[];
}

export interface EventUpdateForm {
    id: string;
    title?: string;
    description?: string;
    allDay?: boolean;
    location?: string;
    url?: string;
    selected_dates?: PrioritizedSelectedDate[];
}