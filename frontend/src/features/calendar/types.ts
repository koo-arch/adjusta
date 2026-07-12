export interface CalendarEventBase {
    id: string;
    title: string;
    start: Date;
    end: Date;
    location?: string;
    description?: string;
    // FullCalendar が解釈するイベント個別の色。未指定は --fc-event-* の既定(Primary)
    backgroundColor?: string;
    borderColor?: string;
    textColor?: string;
}

export interface GoogleCalendarEvent extends CalendarEventBase {
    origin: 'google';
    local_event_id: null;
}

export interface LocalCalendarEvent extends CalendarEventBase {
    origin: 'local';
    local_event_id: string;
}

export type CalendarEvent = GoogleCalendarEvent | LocalCalendarEvent;

export interface GoogleEvent {
    id: string;
    summary: string;
    description: string;
    location: string;
    color: string;
    start: Date;
    end: Date;
}

export interface WarningCalendars {
    failed_calendars: string[];
}

export interface GoogleCalendarResponse {
    events: GoogleEvent[];
    warning: WarningCalendars;
}
