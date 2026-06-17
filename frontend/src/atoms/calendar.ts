import { atomFamily, atomWithStorage } from "jotai/utils";
import { atom } from "jotai";
import { validateUUID } from "@/lib/validation/uuid";
import { CalendarEvent } from "@/features/calendar/type";
import { SendSelectedDate } from "@/features/events/zod";
import type { EventProposedDate } from "@/hooks/event/type";

export interface SelectedDate {
    id: string;
    start: Date;
    end: Date;
}

export interface SelectedEvent extends SelectedDate {
    title: string;
    origin: "google" | "local";
}

// 既存の日付データを保存するatom
export const selectedDatesAtom = atomWithStorage<SelectedDate[]>("selectedDates", []);

export const titleAtomFamily = atomFamily((eventID?: string) => atom(""));

// 日付に基づいてイベントを生成するatom
export const selectedEventsAtomFamily = atomFamily((eventID?: string) => {
    return atom<SelectedEvent[]>((get) => {
        const selectedDates = get(selectedDatesAtom);
        let title = get(titleAtomFamily(eventID));

        if (!title) {
            title = "新しいイベント";
        }

        // selectedDatesAtomに基づいて初期のselectedEventsを生成
        return selectedDates.map((date, index) => ({
            ...date,
            title: `${title} ${index + 1}`,
            origin: "local",
        }));
    })
});

// SelectedDateを送信するために調整したatom
export const sendSelectedDatesAtom= atom<SendSelectedDate[]>(
    (get) => {
        const selectedDates = get(selectedDatesAtom);
        return selectedDates.map((date, index) => ({
            ...date,
            id: validateUUID(date.id) ? date.id : null,
            priority: index + 1,
        }))
    }
);

type ProposedDateMetadata = Pick<EventProposedDate, "google_event_id" | "status" | "sync_status" | "last_synced_at" | "last_sync_error">;

export interface ProposedDate extends SelectedDate, Partial<ProposedDateMetadata> {
    priority: number;
}

export interface ProposedEvent extends SelectedEvent {
    title: string;
}

export interface SendProposedDate {
    id: string | null;
    start: Date;
    end: Date;
    priority: number;
}


// 編集する候補日程を保存するatom
export const proposedDatesAtom = atom<ProposedDate[]>([]);

// proposedDatesAtomに基づいてイベントを生成するatom
export const proposedEventsAtomFamily = atomFamily((eventID?: string) => {
    return atom<ProposedEvent[]>((get) => {
        const proposedDates = get(proposedDatesAtom);
        let title = get(titleAtomFamily(eventID));

        if (!title) {
            title = "新しいイベント";
        }

        return proposedDates.map((date, index) => ({
            ...date,
            title: `${title} ${index + 1}`,
            origin: "local",
            priority: index + 1,
        }));
    })
});

// ProposedDateを送信するために調整したatom
export const sendProposedDatesAtom = atom<SendProposedDate[]>(
    (get) => {
        const proposedDates = get(proposedDatesAtom);
        return proposedDates.map((date, index) => ({
            id: validateUUID(date.id) ? date.id : null,
            start: date.start,
            end: date.end,
            priority: index + 1,
        }));
    }
);

export const allEventsAtom = atom<CalendarEvent[]>([]);
