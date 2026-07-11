'use client'
import React from 'react';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';
import {
    locationAtomFamily,
    proposedDatesAtomFamily,
    titleAtomFamily,
} from '@/features/events/store/calendar';
import { formStepAtomFamily } from '@/features/events/store/formStep';
import type { ProposedDate } from '@/features/events/store/dates';
import DatesPanelView from '@/features/events/components/form/DatesPanelView';
import { clearEditedEventFieldStateAtomFamily, mergedEventFormErrorsAtomFamily } from '@/features/events/store/errors';

interface EditDatesPanelProps {
    formScope: string;
}

const EditDatesPanel: React.FC<EditDatesPanelProps> = ({ formScope }) => {
    const [dates, setDates] = useAtom(proposedDatesAtomFamily(formScope));
    const title = useAtomValue(titleAtomFamily(formScope));
    const location = useAtomValue(locationAtomFamily(formScope));
    const errors = useAtomValue(mergedEventFormErrorsAtomFamily(formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(formScope));
    const setStep = useSetAtom(formStepAtomFamily(formScope));

    const handleDatesChange: React.Dispatch<React.SetStateAction<ProposedDate[]>> = (value) => {
        setDates(value);
        clearEditedFieldState('proposed_dates');
    };

    const handleAdd = ({ start, end }: { start: Date; end: Date }) => {
        handleDatesChange((prev) => {
            const newDate: ProposedDate = {
                id: new Date().getTime().toString(),
                start,
                end,
                priority: prev.length + 1,
            };
            return [...prev, newDate];
        });
    };

    return (
        <DatesPanelView
            title={title}
            location={location}
            dates={dates}
            onDatesChange={handleDatesChange}
            onAdd={handleAdd}
            onEditBasic={() => setStep('basic')}
            error={errors.proposed_dates}
        />
    );
};

export default EditDatesPanel;
