import { atom } from 'jotai';

export const authErrorAtom = atom({
    isOpen: false,
    message: '',
})