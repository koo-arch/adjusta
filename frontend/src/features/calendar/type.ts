interface CalendarEventBase {
    id: string;
    title: string;
    start: Date;
    end: Date;
    location?: string;
    description?: string;
}

interface GoogleCalendarEvent extends CalendarEventBase {
    origin: "google";
    local_event_id: null;
}

interface LocalCalendarEvent extends CalendarEventBase {
    origin: "local";
    local_event_id: string;
}

export type CalendarEvent = GoogleCalendarEvent | LocalCalendarEvent;
