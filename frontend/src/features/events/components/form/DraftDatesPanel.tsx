'use client'
import React from 'react';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';
import {
    locationAtomFamily,
    selectedDatesAtomFamily,
    titleAtomFamily,
} from '@/features/events/store/calendar';
import { formStepAtomFamily } from '@/features/events/store/formStep';
import type { SelectedDate } from '@/features/events/store/dates';
import DatesPanelView from '@/features/events/components/form/DatesPanelView';
import { clearEditedEventFieldStateAtomFamily, mergedEventFormErrorsAtomFamily } from '@/features/events/store/errors';

interface DraftDatesPanelProps {
    formScope: string;
}

const DraftDatesPanel: React.FC<DraftDatesPanelProps> = ({ formScope }) => {
    const [dates, setDates] = useAtom(selectedDatesAtomFamily(formScope));
    const title = useAtomValue(titleAtomFamily(formScope));
    const location = useAtomValue(locationAtomFamily(formScope));
    const errors = useAtomValue(mergedEventFormErrorsAtomFamily(formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(formScope));
    const setStep = useSetAtom(formStepAtomFamily(formScope));

    const handleDatesChange: React.Dispatch<React.SetStateAction<SelectedDate[]>> = (value) => {
        setDates(value);
        clearEditedFieldState('selected_dates');
    };

    const handleAdd = ({ start, end }: { start: Date; end: Date }) => {
        const newDate: SelectedDate = { id: new Date().getTime().toString(), start, end };
        handleDatesChange((prev) => [...prev, newDate]);
    };

    return (
        <DatesPanelView
            title={title}
            location={location}
            dates={dates}
            onDatesChange={handleDatesChange}
            onAdd={handleAdd}
            onEditBasic={() => setStep('basic')}
            error={errors.selected_dates}
        />
    );
};

export default DraftDatesPanel;
