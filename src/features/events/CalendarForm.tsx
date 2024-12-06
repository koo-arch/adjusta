'use client'
import React from 'react';
import { useParams } from 'next/navigation';
import SelectableCalendar from '@/features/calendar/SelectableCalendar';
import type { EventDraftDetail } from '@/hooks/event/type';
import { 
    selectedDatesAtom,
    selectedEventsAtomFamily,
    proposedDatesAtom,
    proposedEventsAtomFamily,
} from '@/atoms/calendar';

interface CalendarFormProps {
    formType: 'draft' | 'edit';
    editingEvent?: EventDraftDetail
}

const CalendarForm: React.FC<CalendarFormProps> = ({ formType, editingEvent }) => {
    const { id } = useParams<{ id?: string}>();
    const dateAtom = formType === 'draft' ? selectedDatesAtom : proposedDatesAtom;
    const eventAtom = formType === 'draft' ? selectedEventsAtomFamily(id) : proposedEventsAtomFamily(id);

    return (
        <div>
            <h2 className="text-lg font-bold mb-2">カレンダー</h2>
            <p className="text-sm text-gray-500 mb-4">カレンダー上をクリックすることで日程選択ができます</p>
            <SelectableCalendar
                dateAtom={dateAtom}
                eventAtom={eventAtom}
                editingEvent={editingEvent}
            />
        </div>
    )
}

export default CalendarForm;