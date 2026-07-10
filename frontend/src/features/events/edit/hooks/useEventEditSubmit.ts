'use client'
import { useAtomValue } from 'jotai';
import { useRouter } from 'next/navigation';
import {
    descriptionAtomFamily,
    locationAtomFamily,
    sendProposedDatesAtomFamily,
    titleAtomFamily,
} from '@/features/events/store/calendar';
import { useUpdateDraftMutation } from '@/features/events/edit/hooks/useUpdateDraftMutation';
import type { EventDraftDetail } from '@/features/events/types';
import type { EventUpdateForm } from '@/features/events/schema';

// フォーム atom から更新 payload を組み立てて送信し、成功時は詳細へ遷移する
export const useEventEditSubmit = (eventID: string, eventDetail: EventDraftDetail) => {
    const router = useRouter();
    const updateDraftMutation = useUpdateDraftMutation(eventID);

    const title = useAtomValue(titleAtomFamily(eventID));
    const description = useAtomValue(descriptionAtomFamily(eventID));
    const location = useAtomValue(locationAtomFamily(eventID));
    const proposedDates = useAtomValue(sendProposedDatesAtomFamily(eventID));

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        const payload: EventUpdateForm = {
            id: eventDetail.id,
            form_type: 'edit' as const,
            title,
            description,
            location,
            // 確定操作は詳細画面に一本化(ui-review 3.4)。編集では現在のステータスを維持する
            status: eventDetail.status,
            proposed_dates: proposedDates,
        };

        const updated = await updateDraftMutation.submit(payload);
        if (updated) {
            // 保存後は詳細へ遷移する(作成フローと統一。ui-review P2 #6)
            router.push(`/events/${eventID}`);
        }
    };

    return {
        handleSubmit,
        isPending: updateDraftMutation.isPending,
    };
};
