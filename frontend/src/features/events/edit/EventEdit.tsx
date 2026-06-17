'use client'
import React from 'react';
import { Provider, useAtomValue, useSetAtom } from 'jotai';
import { useHydrateAtoms } from 'jotai/utils';
import Link from 'next/link';
import { useParams } from 'next/navigation';
import {
    descriptionAtomFamily,
    locationAtomFamily,
    proposedDatesAtomFamily,
    sendProposedDatesAtomFamily,
    titleAtomFamily,
} from '@/atoms/calendar';
import { isConfirmedAtomFamily } from '@/atoms/event';
import { setClientEventFormErrorsAtomFamily } from '../form-meta-atoms';
import { useUpdateDraftMutation } from '@/hooks/event/useUpdateDraftMutation';
import { useFetchEventDetail } from '@/hooks/event/useFetchEventDetail';
import { buildZodFieldErrors } from '@/lib/validation/zod';
import EventForm from '../EventForm';
import type { EventDraftDetail } from '@/hooks/event/type';
import {
    EventUpdateFormSchema,
    type EventFormErrors,
} from '../zod';

interface LoadedEventEditProps {
    eventID: string;
    eventDetail: EventDraftDetail;
}

const EventEditFormContent: React.FC<LoadedEventEditProps> = ({ eventID, eventDetail }) => {
    const updateDraftMutation = useUpdateDraftMutation(eventID);
    const setClientErrors = useSetAtom(setClientEventFormErrorsAtomFamily(eventID));

    useHydrateAtoms([
        [titleAtomFamily(eventID), eventDetail.title],
        [descriptionAtomFamily(eventID), eventDetail.description],
        [locationAtomFamily(eventID), eventDetail.location],
        [proposedDatesAtomFamily(eventID), eventDetail.proposed_dates],
        [isConfirmedAtomFamily(eventID), eventDetail.status === 'confirmed'],
    ]);

    const title = useAtomValue(titleAtomFamily(eventID));
    const description = useAtomValue(descriptionAtomFamily(eventID));
    const location = useAtomValue(locationAtomFamily(eventID));
    const proposedDates = useAtomValue(sendProposedDatesAtomFamily(eventID));
    const isConfirmed = useAtomValue(isConfirmedAtomFamily(eventID));

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        const payload = {
            id: eventDetail.id,
            form_type: 'edit' as const,
            title,
            description,
            location,
            status: isConfirmed ? 'confirmed' as const : 'active' as const,
            proposed_dates: proposedDates,
        };

        const result = EventUpdateFormSchema.safeParse(payload);
        if (!result.success) {
            setClientErrors(buildZodFieldErrors<keyof EventFormErrors>(result.error));
            return;
        }

        setClientErrors({});
        await updateDraftMutation.submit(result.data);
    };

    return (
        <form onSubmit={handleSubmit}>
            <EventForm
                formType="edit"
                formScope={eventID}
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

const EventEdit = () => {
    const params = useParams<{ id: string }>();
    const { eventDetail, isLoading, error } = useFetchEventDetail(params.id);

    if (isLoading) {
        return <p>Loading...</p>;
    }

    if (error || !eventDetail) {
        return (
            <div className="py-8 text-center">
                <p className="mb-4 text-sm text-gray-500">イベントが見つかりませんでした。</p>
                <Link href="/schedule/draft" className="text-sm text-indigo-600 hover:underline">
                    イベント一覧へ戻る
                </Link>
            </div>
        );
    }

    return <LoadedEventEdit key={eventDetail.id} eventID={params.id} eventDetail={eventDetail} />;
}

export default EventEdit;
