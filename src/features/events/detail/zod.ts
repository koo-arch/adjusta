import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import dayjs from 'dayjs';

const ConfirmDateSchema = z.object({
    id: z.string().nullable(),
    google_event_id: z.string().optional(),
    start: z.date().or(z.string().transform((val) => new Date(val))), // string を Date に変換
    end: z.date().or(z.string().transform((val) => new Date(val))), // string を Date に変換
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

const ConfirmFormSchema = z.object({
    confirm_date: ConfirmDateSchema,
})

export type ConfirmForm = z.infer<typeof ConfirmFormSchema>;

export const ConfirmFormResolver = zodResolver(ConfirmFormSchema);