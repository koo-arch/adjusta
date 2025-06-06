import { atomFamily, atomWithStorage } from "jotai/utils";
import { atom } from "jotai";
import { atomWithDefault } from "jotai/utils";
import { validateUUID } from "@/lib/validation/uuid";
import { CalendarEvent } from "@/features/calendar/type";
import { SendSelectedDate } from "@/features/events/zod";
import { fetchEventDetailAtomFamily } from "./queries/event";

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

export const titleAtomFamily = atomFamily((slug?: string) => {
    return atomWithDefault((get) => {
        const { data } = get(fetchEventDetailAtomFamily(slug));
        return data?.title || "";
    })
})

// 日付に基づいてイベントを生成するatom
export const selectedEventsAtomFamily = atomFamily((slug?: string) => {
    return atom<SelectedEvent[]>((get) => {
        const selectedDates = get(selectedDatesAtom);
        let title = get(titleAtomFamily(slug));

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

export interface ProposedDate extends SelectedDate {
    priority: number;
}

export interface ProposedEvent extends SelectedEvent {
    title: string;
}

export interface SendProposedDate extends Omit<ProposedDate, "id"> {
    id: string | null;
}


// 編集する候補日程を保存するatom
export const proposedDatesAtom = atom<ProposedDate[]>([]);

// proposedDatesAtomに基づいてイベントを生成するatom
export const proposedEventsAtomFamily = atomFamily((slug?: string) => {
    return atom<ProposedEvent[]>((get) => {
        const proposedDates = get(proposedDatesAtom);
        let title = get(titleAtomFamily(slug));

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
            ...date,
            id: validateUUID(date.id) ? date.id : null,
            priority: index + 1,
        }));
    }
);

export const allEventsAtom = atom<CalendarEvent[]>([]);