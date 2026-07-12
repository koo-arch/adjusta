'use client'
import React from 'react';
import { Provider } from 'jotai';
import { useHydrateAtoms } from 'jotai/utils';
import {
    descriptionAtomFamily,
    locationAtomFamily,
    selectedDatesAtomFamily,
    titleAtomFamily,
} from '@/features/events/store/calendar';
import EventForm from '@/features/events/components/form/EventForm';
import { useEventDraftSubmit } from '@/features/events/hooks/useEventDraftSubmit';

const draftFormScope = 'draft';

const EventDraftContent: React.FC = () => {
    useHydrateAtoms([
        [titleAtomFamily(draftFormScope), ''],
        [descriptionAtomFamily(draftFormScope), ''],
        [locationAtomFamily(draftFormScope), ''],
        [selectedDatesAtomFamily(draftFormScope), []],
    ]);

    const { handleSubmit, isPending } = useEventDraftSubmit(draftFormScope);

    return (
        <form onSubmit={handleSubmit}>
            <EventForm
                formType="draft"
                formScope={draftFormScope}
                submitLabel="登録する"
                isSubmitting={isPending}
            />
        </form>
    )
}

// Provider を張るコンポーネント自身は新しい store に繋がれない(Context は子にのみ届く)ため、
// store を使う部分を EventDraftContent に分離している。Provider はフォーム状態をページ単位で隔離する
const EventDraft: React.FC = () => {
    return (
        <Provider>
            <EventDraftContent />
        </Provider>
    );
};

export default EventDraft;
