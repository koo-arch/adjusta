export interface GoogleEvent {
    id : string;
    summary : string;
    description : string;
    location : string;
    color : string;
    start : Date;
    end : Date;
}

export interface WarningCalendars {
    failed_calendars: string[];
}

export interface GoogleCalendarResponse {
    events: GoogleEvent[];
    warning: WarningCalendars;
}