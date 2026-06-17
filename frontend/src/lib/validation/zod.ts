import type { ZodError } from 'zod';

export const buildZodFieldErrors = <T extends string>(error: ZodError): Partial<Record<T, string>> => {
    const fieldErrors: Partial<Record<T, string>> = {};

    for (const issue of error.issues) {
        const path = issue.path.join('.') as T;
        if (!path || fieldErrors[path]) {
            continue;
        }

        fieldErrors[path] = issue.message;
    }

    return fieldErrors;
};
