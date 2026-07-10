export type EventStatus = 'draft' | 'active' | 'confirmed' | 'cancelled';

export type ProposedDateStatus = 'active' | 'confirmed' | 'not_selected' | 'cancelled';

export type SyncStatus = 'not_synced' | 'pending_sync' | 'synced' | 'sync_failed';

export type EventSortBy = 'created_at' | 'updated_at' | 'title' | 'status';

export type SortOrder = 'asc' | 'desc';

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
    google_event_id?: string;
    last_synced_at?: string;
    last_sync_error?: string;
    proposed_dates: EventProposedDate[];
}

export interface SearchParams {
    title?: string;
    location?: string;
    description?: string;
    status?: EventStatus;
    start_time_gte?: string;
    start_time_lte?: string;
    end_time_gte?: string;
    end_time_lte?: string;
    sort_by?: EventSortBy;
    sort_order?: SortOrder;
    page?: number;
    per_page?: number;
}

export interface Pagination {
    page: number;
    per_page: number;
    total_items: number;
    total_pages: number;
}

export interface EventDraftListResponse {
    items: EventDraftDetail[];
    pagination: Pagination;
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
    start: Date;
    end: Date;
    needs_attention: boolean;
}
