'use client'
import { useAtomValue } from 'jotai';
import { useRouter } from 'next/navigation';
import {
    descriptionAtomFamily,
    locationAtomFamily,
    sendSelectedDatesAtomFamily,
    titleAtomFamily,
} from '@/features/events/store/calendar';
import { useCreateDraftMutation } from '@/features/events/hooks/useCreateDraftMutation';
import type { EventDraftForm } from '@/features/events/schema';

// フォーム atom から作成 payload を組み立てて送信し、成功時は作成したイベントの詳細へ遷移する
export const useEventDraftSubmit = (formScope: string) => {
    const router = useRouter();
    const createDraftMutation = useCreateDraftMutation(formScope);

    const title = useAtomValue(titleAtomFamily(formScope));
    const description = useAtomValue(descriptionAtomFamily(formScope));
    const location = useAtomValue(locationAtomFamily(formScope));
    const selectedDates = useAtomValue(sendSelectedDatesAtomFamily(formScope));

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        const payload: EventDraftForm = {
            form_type: 'draft' as const,
            title,
            description,
            location,
            selected_dates: selectedDates,
        };

        const createdDraftID = await createDraftMutation.submit(payload);
        if (createdDraftID) {
            router.push(`/events/${createdDraftID}`);
        }
    };

    return {
        handleSubmit,
        isPending: createDraftMutation.isPending,
    };
};
