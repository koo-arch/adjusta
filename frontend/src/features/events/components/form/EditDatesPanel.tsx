'use client'
import React from 'react';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';
import {
    locationAtomFamily,
    proposedDatesAtomFamily,
    titleAtomFamily,
} from '@/features/events/store/calendar';
import { formStepAtomFamily } from '@/features/events/store/formStep';
import { buildDefaultCandidateSpan, type ProposedDate } from '@/features/events/store/dates';
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

    const handleAdd = () => {
        const newDate: ProposedDate = {
            id: new Date().getTime().toString(),
            ...buildDefaultCandidateSpan(),
            priority: dates.length + 1,
        };
        handleDatesChange((prev) => [...prev, newDate]);
        return newDate.id;
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
