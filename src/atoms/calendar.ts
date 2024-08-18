import { atomWithStorage } from "jotai/utils";
import { atom } from "jotai";

interface SelectedDate {
    id: string;
    start: Date;
    end: Date;
}

export interface SelectedEvent extends SelectedDate {
    title: string;
    origin: string;
}

// 既存の日付データを保存するatom
export const selectedDatesAtom = atomWithStorage<SelectedDate[]>("selectedDates", []);

export const titleAtom = atomWithStorage<string>("title", "選択されたイベント");

// 日付に基づいてイベントを生成するatom
export const selectedEventsAtom = atom<SelectedEvent[]>(
    (get) => {
        const selectedDates = get(selectedDatesAtom);
        const title = get(titleAtom);

        // selectedDatesAtomに基づいて初期のselectedEventsを生成
        return selectedDates.map((date, index) => ({
            ...date,
            title: `${title} ${index + 1}`,
            origin: "local",
        }));
    }
);