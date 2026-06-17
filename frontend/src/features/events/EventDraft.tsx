'use client'
import React from 'react';
import { Provider, useAtomValue, useSetAtom } from 'jotai';
import { useHydrateAtoms } from 'jotai/utils';
import { useRouter } from 'next/navigation';
import {
    descriptionAtomFamily,
    locationAtomFamily,
    selectedDatesAtomFamily,
    sendSelectedDatesAtomFamily,
    titleAtomFamily,
} from '@/atoms/calendar';
import { setClientEventFormErrorsAtomFamily } from './form-meta-atoms';
import EventForm from './EventForm';
import { useCreateDraftMutation } from '@/hooks/event/useCreateDraftMutation';
import { buildZodFieldErrors } from '@/lib/validation/zod';
import {
    EventDraftFormSchema,
    type EventFormErrors,
} from './zod';

const draftFormScope = 'draft';

const EventDraftContent: React.FC = () => {
    const router = useRouter();
    const createDraftMutation = useCreateDraftMutation(draftFormScope);
    const setClientErrors = useSetAtom(setClientEventFormErrorsAtomFamily(draftFormScope));

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

        const payload = {
            form_type: 'draft' as const,
            title,
            description,
            location,
            selected_dates: selectedDates,
        };

        const result = EventDraftFormSchema.safeParse(payload);
        if (!result.success) {
            setClientErrors(buildZodFieldErrors<keyof EventFormErrors>(result.error));
            return;
        }

        setClientErrors({});
        const createdDraftID = await createDraftMutation.submit(result.data);
        if (createdDraftID) {
            router.push(`/schedule/draft/${createdDraftID}`);
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
