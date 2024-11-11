export type CalendarEvent = {
    id: string;
    title: string;
    start: Date;
    end: Date;
    origin: "google" | "local";
}