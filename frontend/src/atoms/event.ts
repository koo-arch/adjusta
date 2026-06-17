import { fetchEventDetailAtomFamily } from './queries/event';
import { atomWithDefault, atomFamily } from 'jotai/utils';

export const isConfirmedAtomFamily = atomFamily((eventID?: string) => {
    return atomWithDefault((get) => {
        const { data } = get(fetchEventDetailAtomFamily(eventID));
        return data?.status === 'confirmed';
    })
});
