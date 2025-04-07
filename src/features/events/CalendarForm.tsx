'use client'
import React from 'react';
import { useFormContext } from 'react-hook-form';
import { useParams } from 'next/navigation';
import SelectableCalendar from '@/features/calendar/SelectableCalendar';
import type { EventDraftDetail } from '@/hooks/event/type';
import { 
    selectedDatesAtom,
    selectedEventsAtomFamily,
    proposedDatesAtom,
    proposedEventsAtomFamily,
} from '@/atoms/calendar';
import type { DiscriminatedEventForm } from './zod';

interface CalendarFormProps {
    editingEvent?: EventDraftDetail
}

const CalendarForm: React.FC<CalendarFormProps> = ({ editingEvent }) => {
    const { slug } = useParams<{ slug?: string}>();
    const { getValues } = useFormContext<DiscriminatedEventForm>();
    
    const formType = getValues('form_type');
    const dateAtom = formType === 'draft' ? selectedDatesAtom : proposedDatesAtom;
    const eventAtom = formType === 'draft' ? selectedEventsAtomFamily(slug) : proposedEventsAtomFamily(slug);

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