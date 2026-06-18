import { z } from 'zod';
import dayjs from 'dayjs';

const buildDateSchema = (requiredMessage: string, invalidMessage: string) =>
    z.preprocess((value) => {
        if (value == null || value === '') {
            return undefined;
        }

        if (value instanceof Date) {
            return value;
        }

        if (typeof value === 'string' && value.length > 0) {
            const parsed = new Date(value);
            return Number.isNaN(parsed.getTime()) ? value : parsed;
        }

        return value;
    }, z.date({ required_error: requiredMessage, invalid_type_error: invalidMessage }));

export const ConfirmDateSchema = z.object({
    id: z.string().nullable(),
    google_event_id: z.string().optional(),
    start: buildDateSchema('開始日時は必須です', '開始日時は不正な形式です'),
    end: buildDateSchema('終了日時は必須です', '終了日時は不正な形式です'),
    priority: z.number(),
}).refine(
    (args) => {
        const { start, end } = args;
        const startDate = dayjs(start);
        const endDate = dayjs(end);
        return endDate.isAfter(startDate);
    },
    {
        message: '終了日時は開始日時より後に設定してください',
        path: ['end']
    }
)

export const ConfirmFormSchema = z.object({
    confirm_date: ConfirmDateSchema,
})

export type ConfirmForm = z.infer<typeof ConfirmFormSchema>;
export type ConfirmFormErrors = Partial<Record<'confirm_date' | 'confirm_date.start' | 'confirm_date.end', string>>;
