import { atomWithStorage } from 'jotai/utils';

export const authAtom = atomWithStorage<boolean>('isAuthenticated', false);

export const signInAtom = atomWithStorage<boolean>('isSignIn', false);