'use client'
import React from 'react';
import { Provider, useAtomValue } from 'jotai';
import { useHydrateAtoms } from 'jotai/utils';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import {
    descriptionAtomFamily,
    locationAtomFamily,
    proposedDatesAtomFamily,
    sendProposedDatesAtomFamily,
    titleAtomFamily,
} from '@/features/events/store/calendar';
import { useUpdateDraftMutation } from '@/features/events/edit/hooks/useUpdateDraftMutation';
import { useFetchEventDetail } from '@/features/events/hooks/useFetchEventDetail';
import { isNotFoundAPIError } from '@/lib/api/errors';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import EventForm from '@/features/events/components/form/EventForm';
import type { EventDraftDetail } from '@/features/events/types';
import type { EventUpdateForm } from '@/features/events/schema';

interface LoadedEventEditProps {
    eventID: string;
    eventDetail: EventDraftDetail;
}

const EventEditFormContent: React.FC<LoadedEventEditProps> = ({ eventID, eventDetail }) => {
    const router = useRouter();
    const updateDraftMutation = useUpdateDraftMutation(eventID);

    useHydrateAtoms([
        [titleAtomFamily(eventID), eventDetail.title],
        [descriptionAtomFamily(eventID), eventDetail.description],
        [locationAtomFamily(eventID), eventDetail.location],
        [proposedDatesAtomFamily(eventID), eventDetail.proposed_dates],
    ]);

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

    return (
        <form onSubmit={handleSubmit}>
            <EventForm
                formType="edit"
                formScope={eventID}
                submitLabel="保存する"
                isSubmitting={updateDraftMutation.isPending}
                eventDetail={eventDetail}
            />
        </form>
    );
};

const LoadedEventEdit: React.FC<LoadedEventEditProps> = ({ eventID, eventDetail }) => {
    return (
        <Provider>
            <EventEditFormContent eventID={eventID} eventDetail={eventDetail} />
        </Provider>
    );
};

export const EventFormSkeleton = () => (
    <div className="grid grid-cols-1 gap-8 md:grid-cols-10 md:gap-6">
        <div className="space-y-4 md:col-span-4">
            <Skeleton className="h-6 w-24" />
            <Skeleton className="h-9 w-full" />
            <Skeleton className="h-9 w-full" />
            <Skeleton className="h-24 w-full" />
        </div>
        <div className="md:col-span-6">
            <Skeleton className="h-96 w-full" />
        </div>
    </div>
);

interface EventEditProps {
    eventID: string;
}

const EventEdit: React.FC<EventEditProps> = ({ eventID }) => {
    const { eventDetail, isPending, error, refetch } = useFetchEventDetail(eventID);

    if (isPending) {
        return <EventFormSkeleton />;
    }

    if (isNotFoundAPIError(error)) {
        return (
            <div className="flex flex-col items-center gap-4 py-16 text-center">
                <p className="text-sm text-muted-foreground">イベントが見つかりませんでした。</p>
                <Button variant="outline" asChild>
                    <Link href="/events">イベント一覧へ戻る</Link>
                </Button>
            </div>
        );
    }

    if (error || !eventDetail) {
        return (
            <div className="flex flex-col items-center gap-4 py-16 text-center">
                <p className="text-sm text-muted-foreground">
                    イベントの取得に失敗しました。時間をおいて再度お試しください。
                </p>
                <Button variant="outline" onClick={() => refetch()}>
                    再試行
                </Button>
            </div>
        );
    }

    return <LoadedEventEdit key={eventDetail.id} eventID={eventID} eventDetail={eventDetail} />;
}

export default EventEdit;
