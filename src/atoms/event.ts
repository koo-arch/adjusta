import { fetchEventDetailAtomFamily } from './queries/event';
import { atomWithDefault, atomFamily } from 'jotai/utils';

export const isConfirmedAtomFamily = atomFamily((slug?: string) => {
    return atomWithDefault((get) => {
        const { data } = get(fetchEventDetailAtomFamily(slug));
        return data?.status === 'confirmed';
    })
});