import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';

const EventBasicFormSchema = z.object({
    title: z.string().min(1, { message: 'タイトルは必須です'}).max(100, { message: 'タイトルは100文字以内で入力してください'}),
    description: z.string().max(500, { message: '詳細は500文字以内で入力してください'}),
    location: z.string().max(100, { message: '場所は100文字以内で入力してください'}),
    allDay: z.boolean().optional(),
    url: z.string().url().optional(),
});

const SendSelectedDateSchema = z.object({
    id: z.string().nullable(),
    start: z.date(),
    end: z.date(),
    priority: z.number(),
});

const EventDraftFormSchema = EventBasicFormSchema.merge(z.object({
    form_type: z.literal('draft'),
    selected_dates: z.array(SendSelectedDateSchema)
        .min(1, { message: '日程は1つ以上選択してください' })
        .max(10, { message: '日程は10個まで選択できます' }),
}));

const EventUpdateFormSchema = EventBasicFormSchema.merge(z.object({
    form_type: z.literal('edit'),
    id: z.string(),
    status: z.string().optional(),
    proposed_dates: z.array(SendSelectedDateSchema)
        .min(1, { message: '日程は1つ以上選択してください' })
        .max(10, { message: '日程は10個まで選択できます' }),
}));

const DiscriminatedEventFormSchema = z.discriminatedUnion('form_type', [EventDraftFormSchema, EventUpdateFormSchema]);

export type DiscriminatedEventForm = z.infer<typeof DiscriminatedEventFormSchema>;

export type SendSelectedDate = z.infer<typeof SendSelectedDateSchema>;

export const DiscriminatedEventFormResolver = zodResolver(DiscriminatedEventFormSchema);