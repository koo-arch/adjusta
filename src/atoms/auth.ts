import { atomWithStorage } from 'jotai/utils';

export const authAtom = atomWithStorage<boolean>('isAuthenticated', false);