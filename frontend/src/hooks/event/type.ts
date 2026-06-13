export type EventStatus = 'draft' | 'active' | 'confirmed' | 'cancelled';

export type ProposedDateStatus = 'active' | 'confirmed' | 'not_selected' | 'cancelled';

export type SyncStatus = 'not_synced' | 'pending_sync' | 'synced' | 'sync_failed';

export interface EventProposedDate {
    id: string;
    google_event_id?: string;
    start: Date;
    end: Date;
    priority: number;
    status: ProposedDateStatus;
    sync_status: SyncStatus;
    last_synced_at?: string;
    last_sync_error?: string;
}

export interface EventDraftDetail {
    id: string;
    title: string;
    description: string;
    location: string;
    status: EventStatus;
    sync_status: SyncStatus;
    confirmed_date_id: string | null;
    confirmed_google_event_id?: string;
    google_event_id: string;
    last_synced_at?: string;
    last_sync_error?: string;
    slug: string;
    proposed_dates: EventProposedDate[];
}

export interface SearchParams {
    title?: string;
    location?: string;
    startTime?: string;
    endTime?: string;
    status?: EventStatus;
}

export interface UpcomingEvent extends Omit<EventDraftDetail, 'proposed_dates'> {
    start: Date;
    end: Date;
}

export interface NeedsActionDraft {
    id: string;
    title: string;
    location: string;
    description: string;
    status: EventStatus;
    slug: string;
    start: Date;
    end: Date;
    needs_attention: boolean;
}
