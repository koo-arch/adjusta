import { atomWithStorage } from "jotai/utils";
import { atom } from "jotai";

export interface SelectedDate {
    id: string;
    start: Date;
    end: Date;
}

export interface SelectedEvent extends SelectedDate {
    title: string;
    origin: string;
}

export interface PrioritizedSelectedDate extends SelectedDate {
    priority: number;
}

// 既存の日付データを保存するatom
export const selectedDatesAtom = atomWithStorage<SelectedDate[]>("selectedDates", []);

export const titleAtom = atomWithStorage<string>("title", "");

// 日付に基づいてイベントを生成するatom
export const selectedEventsAtom = atom<SelectedEvent[]>(
    (get) => {
        const selectedDates = get(selectedDatesAtom);
        let title = get(titleAtom);

        if (!title) {
            title = "新しいイベント";
        }

        // selectedDatesAtomに基づいて初期のselectedEventsを生成
        return selectedDates.map((date, index) => ({
            ...date,
            title: `${title} ${index + 1}`,
            origin: "local",
        }));
    }
);

export const prioritizedSelectedDatesAtom = atom<PrioritizedSelectedDate[]>(
    (get) => {
        const selectedDates = get(selectedDatesAtom);
        return selectedDates.map((date, index) => ({
            ...date,
            priority: index + 1,
        }))
    }
);