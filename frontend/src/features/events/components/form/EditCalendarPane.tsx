'use client'
import React from 'react';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';
import { proposedDatesAtomFamily, proposedEventsAtomFamily } from '@/features/events/store/calendar';
import type { ProposedDate } from '@/features/events/store/dates';
import SelectableCalendar from '@/features/calendar/components/SelectableCalendar';
import { clearEditedEventFieldStateAtomFamily } from '@/features/events/store/errors';
import type { EventDraftDetail } from '@/features/events/types';

interface EditCalendarPaneProps {
    formScope: string;
    editingEvent?: EventDraftDetail;
}

const EditCalendarPane: React.FC<EditCalendarPaneProps> = ({ formScope, editingEvent }) => {
    const [dates, setDates] = useAtom(proposedDatesAtomFamily(formScope));
    const proposedEvents = useAtomValue(proposedEventsAtomFamily(formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(formScope));

    const handleDatesChange: React.Dispatch<React.SetStateAction<ProposedDate[]>> = (value) => {
        setDates(value);
        clearEditedFieldState('proposed_dates');
    };

    return (
        <SelectableCalendar
            dates={dates}
            onDatesChange={handleDatesChange}
            selectedEvents={proposedEvents}
            editingEvent={editingEvent}
        />
    );
};

export default EditCalendarPane;
