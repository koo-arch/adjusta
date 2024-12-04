'use client'
import React from 'react';
import { selectedDatesAtom, selectedEventsAtom } from '@/atoms/calendar';
import SelectableCalendar from '@/features/calendar/SelectableCalendar';


const CalendarForm: React.FC = () => {
    return (
        <div>
            <h2 className="text-lg font-bold mb-2">カレンダー</h2>
            <p className="text-sm text-gray-500 mb-4">カレンダー上をクリックすることで日程選択ができます</p>
            <SelectableCalendar
                dateAtom={selectedDatesAtom}
                eventAtom={selectedEventsAtom}
            />
        </div>
    )
}

export default CalendarForm;