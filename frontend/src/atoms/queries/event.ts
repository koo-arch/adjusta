import { atom } from 'jotai';
import { atomFamily } from 'jotai/utils';
import { atomWithQuery } from 'jotai-tanstack-query';
import axios from '@/lib/axios/public';
import { EventDraftDetail } from '@/hooks/event/type';

export const eventDetailIdAtom = atom<string | null>(null);

export const fetchEventDetailAtomFamily = atomFamily((eventID?: string) =>
    atomWithQuery((get) => ({
        queryKey: ['eventDetail', eventID],
        queryFn: async () => {
            const { data } = await axios.get<EventDraftDetail>(`/api/calendar/event/draft/${eventID}`);
            return data;
        },
        enabled: !!eventID,
    }))
);
