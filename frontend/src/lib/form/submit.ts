import type { FormErrors } from './errors';

export interface SubmitSuccess<T> {
    ok: true;
    data: T;
}

export interface ValidationSubmitFailure<FieldKey extends string> {
    ok: false;
    type: 'validation';
    errors: Partial<Record<FieldKey, string>>;
}

export interface RequestSubmitFailure<FieldKey extends string> {
    ok: false;
    type: 'request';
    errors: FormErrors<FieldKey>;
}

export type SubmitResult<T, FieldKey extends string> =
    | SubmitSuccess<T>
    | ValidationSubmitFailure<FieldKey>
    | RequestSubmitFailure<FieldKey>;
