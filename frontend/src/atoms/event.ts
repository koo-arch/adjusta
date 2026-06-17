import { atomFamily, atomWithReset } from 'jotai/utils';

export const isConfirmedAtomFamily = atomFamily((formScope: string) => atomWithReset(false));
