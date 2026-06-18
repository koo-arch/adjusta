import type { CalendarEvent } from '@/features/calendar/type';
import type { EventProposedDate } from '@/features/events/types';
import { validateUUID } from '@/lib/validation/uuid';
import type { SendSelectedDate } from '../schema';

export interface SelectedDate {
    id: string;
    start: Date;
    end: Date;
}

type ProposedDateMetadata = Pick<EventProposedDate, 'google_event_id' | 'status' | 'sync_status' | 'last_synced_at' | 'last_sync_error'>;

export interface ProposedDate extends SelectedDate, Partial<ProposedDateMetadata> {
    priority: number;
}

export interface SendProposedDate {
    id: string | null;
    start: Date;
    end: Date;
    priority: number;
}

export const buildLocalCalendarEvents = <T extends SelectedDate>(dates: T[], title: string): CalendarEvent[] => {
    const eventTitle = title || '新しいイベント';

    return dates.map((date, index) => ({
        ...date,
        title: `${eventTitle} ${index + 1}`,
        origin: 'local' as const,
        local_event_id: '',
    }));
};

export const buildSendSelectedDates = (selectedDates: SelectedDate[]): SendSelectedDate[] =>
    selectedDates.map((date, index) => ({
        ...date,
        id: validateUUID(date.id) ? date.id : null,
        priority: index + 1,
    }));

export const buildSendProposedDates = (proposedDates: ProposedDate[]): SendProposedDate[] =>
    proposedDates.map((date, index) => ({
        id: validateUUID(date.id) ? date.id : null,
        start: date.start,
        end: date.end,
        priority: index + 1,
    }));
