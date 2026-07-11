'use client'
import React from 'react';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';
import { selectedDatesAtomFamily, selectedEventsAtomFamily } from '@/features/events/store/calendar';
import type { SelectedDate } from '@/features/events/store/dates';
import SelectableCalendar from '@/features/calendar/components/SelectableCalendar';
import { clearEditedEventFieldStateAtomFamily } from '@/features/events/store/errors';
import type { EventDraftDetail } from '@/features/events/types';

interface DraftCalendarPaneProps {
    formScope: string;
    editingEvent?: EventDraftDetail;
}

const DraftCalendarPane: React.FC<DraftCalendarPaneProps> = ({ formScope, editingEvent }) => {
    const [dates, setDates] = useAtom(selectedDatesAtomFamily(formScope));
    const selectedEvents = useAtomValue(selectedEventsAtomFamily(formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(formScope));

    const handleDatesChange: React.Dispatch<React.SetStateAction<SelectedDate[]>> = (value) => {
        setDates(value);
        clearEditedFieldState('selected_dates');
    };

    return (
        <SelectableCalendar
            dates={dates}
            onDatesChange={handleDatesChange}
            selectedEvents={selectedEvents}
            editingEvent={editingEvent}
        />
    );
};

export default DraftCalendarPane;
