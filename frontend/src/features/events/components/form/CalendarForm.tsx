'use client'
import React from 'react';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';
import {
    proposedDatesAtomFamily,
    proposedEventsAtomFamily,
    selectedDatesAtomFamily,
    selectedEventsAtomFamily,
} from '@/features/events/store/calendar';
import SelectableCalendar from '@/features/calendar/SelectableCalendar';
import type { EventDraftDetail } from '@/features/events/types';
import { clearEditedEventFieldStateAtomFamily } from '@/features/events/store/errors';

type DraftCalendarFormProps = {
    formType: 'draft';
    formScope: string;
    editingEvent?: EventDraftDetail;
};

type EditCalendarFormProps = {
    formType: 'edit';
    formScope: string;
    editingEvent?: EventDraftDetail;
};

type CalendarFormProps = DraftCalendarFormProps | EditCalendarFormProps;

const CalendarDescription = () => (
    <>
        <h2 className="text-lg font-bold mb-2">カレンダー</h2>
        <p className="text-sm text-gray-500 mb-4">カレンダー上をクリックすることで日程選択ができます</p>
    </>
);

const DraftCalendarSection: React.FC<DraftCalendarFormProps> = ({ formScope, editingEvent }) => {
    const [dates, setDates] = useAtom(selectedDatesAtomFamily(formScope));
    const selectedEvents = useAtomValue(selectedEventsAtomFamily(formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(formScope));

    return (
        <div>
            <CalendarDescription />
            <SelectableCalendar
                dates={dates}
                onDatesChange={(value) => {
                    setDates(value);
                    clearEditedFieldState('selected_dates');
                }}
                selectedEvents={selectedEvents}
                editingEvent={editingEvent}
            />
        </div>
    );
};

const EditCalendarSection: React.FC<EditCalendarFormProps> = ({ formScope, editingEvent }) => {
    const [dates, setDates] = useAtom(proposedDatesAtomFamily(formScope));
    const selectedEvents = useAtomValue(proposedEventsAtomFamily(formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(formScope));

    return (
        <div>
            <CalendarDescription />
            <SelectableCalendar
                dates={dates}
                onDatesChange={(value) => {
                    setDates(value);
                    clearEditedFieldState('proposed_dates');
                }}
                selectedEvents={selectedEvents}
                editingEvent={editingEvent}
            />
        </div>
    );
};

const CalendarForm: React.FC<CalendarFormProps> = (props) => {
    if (props.formType === 'draft') {
        return <DraftCalendarSection {...props} />;
    }

    return <EditCalendarSection {...props} />;
};

export default CalendarForm;
