export type CalendarEvent = {
    id: string;
    title: string;
    start: Date;
    end: Date;
    location?: string;
    description?: string;
    origin: "google" | "local";
    local_event_id?: string;
}