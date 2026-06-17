import { atom } from 'jotai';
import { atomFamily, atomWithReset } from 'jotai/utils';
import { emptyFormErrors, type FormErrors } from '@/lib/form/errors';
import type { EventFormEditedField, EventFormErrors } from './zod';

type EventServerErrors = FormErrors<keyof EventFormErrors>;

const createEmptyEventServerErrors = (): EventServerErrors =>
    emptyFormErrors<keyof EventFormErrors>();

// formScope は atom の値計算では直接使っていなくても、
// 「draft 作成画面用」と「event ごとの edit 画面用」で
// 別インスタンスの atom を持つための key として使う。
//
// 例:
// - clientEventFormErrorsAtomFamily('draft')
// - clientEventFormErrorsAtomFamily('event-123')
//
// これらは初期値が同じでも、別々の error state として扱われる。

// フロント側のバリデーション結果を保持する。
// 例: zod.safeParse で作った title / proposed_dates のエラー
export const clientEventFormErrorsAtomFamily = atomFamily((formScope: string) =>
    atomWithReset<EventFormErrors>({}),
);

// API 送信後にサーバーから返ってきたエラーを保持する。
// formErrors はフォーム全体のエラー文、fieldErrors は項目ごとのエラー文
export const serverEventFormErrorsAtomFamily = atomFamily((formScope: string) =>
    atomWithReset<EventServerErrors>(createEmptyEventServerErrors()),
);

// 画面表示用の統合エラー。
// server 側の fieldErrors をベースにして、client 側のエラーで上書きする。
export const mergedEventFormErrorsAtomFamily = atomFamily((formScope: string) =>
    atom((get) => ({
        ...get(serverEventFormErrorsAtomFamily(formScope)).fieldErrors,
        ...get(clientEventFormErrorsAtomFamily(formScope)),
    })),
);

// 項目単位ではない、フォーム全体メッセージだけを取り出す。
export const eventFormMessagesAtomFamily = atomFamily((formScope: string) =>
    atom((get) => get(serverEventFormErrorsAtomFamily(formScope)).formErrors),
);

// client 側エラーをまとめて差し替えるための write-only atom。
export const setClientEventFormErrorsAtomFamily = atomFamily((formScope: string) =>
    atom(
        null,
        (_get, set, errors: EventFormErrors) => {
            set(clientEventFormErrorsAtomFamily(formScope), errors);
        },
    ),
);

// server 側エラーをまとめて差し替えるための write-only atom。
export const setServerEventFormErrorsAtomFamily = atomFamily((formScope: string) =>
    atom(
        null,
        (_get, set, errors: EventServerErrors) => {
            set(serverEventFormErrorsAtomFamily(formScope), errors);
        },
    ),
);

// フォーム送信前や reset 時に、client/server 両方の error state を初期化する。
export const clearEventFormErrorStateAtomFamily = atomFamily((formScope: string) =>
    atom(
        null,
        (_get, set) => {
            set(clientEventFormErrorsAtomFamily(formScope), {});
            set(serverEventFormErrorsAtomFamily(formScope), createEmptyEventServerErrors());
        },
    ),
);

// ユーザーが特定の項目を編集し始めたときに、その項目に紐づくエラーだけを消す。
// あわせて formErrors も消し、再入力中に古いエラーが残り続けないようにする。
export const clearEditedEventFieldStateAtomFamily = atomFamily((formScope: string) =>
    atom(
        null,
        (get, set, field: EventFormEditedField) => {
            const nextServerErrors = get(serverEventFormErrorsAtomFamily(formScope));

            if (field !== 'confirmed') {
                set(clientEventFormErrorsAtomFamily(formScope), {
                    ...get(clientEventFormErrorsAtomFamily(formScope)),
                    [field]: undefined,
                });
            }

            set(serverEventFormErrorsAtomFamily(formScope), {
                formErrors: [],
                fieldErrors: field === 'confirmed'
                    ? nextServerErrors.fieldErrors
                    : {
                        ...nextServerErrors.fieldErrors,
                        [field]: undefined,
                    },
            });
        },
    ),
);
