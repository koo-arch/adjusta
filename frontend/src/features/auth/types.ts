export interface AuthUser {
    sub: string;
    name: string;
    email: string;
    picture: string;
}

export type UserCalendarRole = 'primary' | 'adjusta_candidate' | 'reference';

export interface CalendarSetting {
    id: string;
    calendar_id: string;
    google_calendar_id: string;
    summary: string;
    description?: string;
    timezone?: string;
    role: UserCalendarRole;
    is_visible: boolean;
    sync_proposed_dates: boolean;
}

export interface CalendarSettingUpdate {
    role?: UserCalendarRole;
    is_visible?: boolean;
    sync_proposed_dates?: boolean;
}

export interface CandidateSyncSetting {
    enabled: boolean;
    calendar: CalendarSetting | null;
}
