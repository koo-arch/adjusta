import { atom } from 'jotai';
import { atomFamily } from 'jotai/utils';

export type EventFormStep = 'basic' | 'dates';

// フォームのステップ遷移(基本情報 ⇄ 候補日程)。スコープ付き Provider 内で使うため
// フォーム離脱時は store ごと破棄される
export const formStepAtomFamily = atomFamily((formScope: string) => atom<EventFormStep>('basic'));
