import { atom } from 'jotai';
import { atomFamily } from 'jotai/utils';

export const isConfirmedAtomFamily = atomFamily((eventID?: string) => {
    return atom(false);
})
