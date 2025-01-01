import { fetchEventDetailAtomFamily } from './queries/event';
import { atomWithDefault, atomFamily } from 'jotai/utils';

export const isConfirmedAtomFamily = atomFamily((id?: string) => {
    return atomWithDefault((get) => {
        const { data } = get(fetchEventDetailAtomFamily(id));
        return data?.status === 'confirmed';
    })
});