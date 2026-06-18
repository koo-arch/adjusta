import { atom } from 'jotai';
import { atomFamily, atomWithReset } from 'jotai/utils';
import type { CalendarEvent } from '@/features/calendar/type';
import {
    buildLocalCalendarEvents,
    buildSendProposedDates,
    buildSendSelectedDates,
    type ProposedDate,
    type SelectedDate,
} from '@/features/events/store/dates';

export const allEventsAtom = atom<CalendarEvent[]>([]);

export const titleAtomFamily = atomFamily((formScope: string) => atomWithReset(''));

export const descriptionAtomFamily = atomFamily((formScope: string) => atomWithReset(''));

export const locationAtomFamily = atomFamily((formScope: string) => atomWithReset(''));

export const selectedDatesAtomFamily = atomFamily((formScope: string) => atomWithReset<SelectedDate[]>([]));

export const proposedDatesAtomFamily = atomFamily((formScope: string) => atomWithReset<ProposedDate[]>([]));

export const selectedEventsAtomFamily = atomFamily((formScope: string) =>
    atom((get) => buildLocalCalendarEvents(
        get(selectedDatesAtomFamily(formScope)),
        get(titleAtomFamily(formScope)),
    )),
);

export const proposedEventsAtomFamily = atomFamily((formScope: string) =>
    atom((get) => buildLocalCalendarEvents(
        get(proposedDatesAtomFamily(formScope)),
        get(titleAtomFamily(formScope)),
    )),
);

export const sendSelectedDatesAtomFamily = atomFamily((formScope: string) =>
    atom((get) => buildSendSelectedDates(get(selectedDatesAtomFamily(formScope)))),
);

export const sendProposedDatesAtomFamily = atomFamily((formScope: string) =>
    atom((get) => buildSendProposedDates(get(proposedDatesAtomFamily(formScope)))),
);
