'use client'
import React from 'react';
import { Provider, useAtomValue } from 'jotai';
import { useHydrateAtoms } from 'jotai/utils';
import { useRouter } from 'next/navigation';
import {
    descriptionAtomFamily,
    locationAtomFamily,
    selectedDatesAtomFamily,
    sendSelectedDatesAtomFamily,
    titleAtomFamily,
} from '@/features/events/store/calendar';
import EventForm from '@/features/events/components/form/EventForm';
import { useCreateDraftMutation } from '@/features/events/hooks/useCreateDraftMutation';
import type { EventDraftForm } from '@/features/events/schema';

const draftFormScope = 'draft';

const EventDraftContent: React.FC = () => {
    const router = useRouter();
    const createDraftMutation = useCreateDraftMutation(draftFormScope);

    useHydrateAtoms([
        [titleAtomFamily(draftFormScope), ''],
        [descriptionAtomFamily(draftFormScope), ''],
        [locationAtomFamily(draftFormScope), ''],
        [selectedDatesAtomFamily(draftFormScope), []],
    ]);

    const title = useAtomValue(titleAtomFamily(draftFormScope));
    const description = useAtomValue(descriptionAtomFamily(draftFormScope));
    const location = useAtomValue(locationAtomFamily(draftFormScope));
    const selectedDates = useAtomValue(sendSelectedDatesAtomFamily(draftFormScope));

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

    return (
        <div>
            <form onSubmit={handleSubmit}>
                <EventForm
                    formType="draft"
                    formScope={draftFormScope}
                    isSubmitting={createDraftMutation.isPending}
                />
            </form>
        </div>
    )
}

const EventDraft: React.FC = () => {
    return (
        <Provider>
            <EventDraftContent />
        </Provider>
    );
};

export default EventDraft;
