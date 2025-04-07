import { atom } from 'jotai';
import { atomFamily } from 'jotai/utils';
import { atomWithQuery } from 'jotai-tanstack-query';
import axios from '@/lib/axios/public';
import { EventDraftDetail } from '@/hooks/event/type';
import { authAtom } from '../auth';

export const eventDetailIdAtom = atom<string | null>(null);

export const fetchEventDetailAtomFamily = atomFamily((slug?: string) =>
    atomWithQuery((get) => ({
        queryKey: ['eventDetail', slug],
        queryFn: async () => {
            const { data } = await axios.get<EventDraftDetail>(`/api/calendar/event/draft/${slug}`);
            return data;
        },
        enabled: !!slug && get(authAtom),
    }))
);