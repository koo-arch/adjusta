import { APIClientError } from '@/lib/api/client';

export interface FormErrors<T extends string> {
    formErrors: string[];
    fieldErrors: Partial<Record<T, string>>;
}

export const emptyFormErrors = <T extends string>(): FormErrors<T> => ({
    formErrors: [],
    fieldErrors: {},
});

const isRecord = (value: unknown): value is Record<string, unknown> =>
    typeof value === 'object' && value !== null;

const normalizeFieldErrors = <T extends string>(details: unknown): Partial<Record<T, string>> => {
    if (!isRecord(details)) {
        return {};
    }

    const fieldErrors: Partial<Record<T, string>> = {};

    for (const [key, rawValue] of Object.entries(details)) {
        if (typeof rawValue === 'string') {
            fieldErrors[key as T] = rawValue;
            continue;
        }

        if (Array.isArray(rawValue) && typeof rawValue[0] === 'string') {
            fieldErrors[key as T] = rawValue[0];
        }
    }

    return fieldErrors;
};

export const buildFormErrorsFromAPIError = <T extends string>(
    error: unknown,
    fallbackMessage: string,
): FormErrors<T> => {
    if (error instanceof APIClientError) {
        const fieldErrors = normalizeFieldErrors<T>(
            isRecord(error.data) ? error.data.details : undefined,
        );

        if (Object.keys(fieldErrors).length > 0) {
            return {
                formErrors: [],
                fieldErrors,
            };
        }

        return {
            formErrors: [error.message || fallbackMessage],
            fieldErrors: {},
        };
    }

    return {
        formErrors: [fallbackMessage],
        fieldErrors: {},
    };
};
