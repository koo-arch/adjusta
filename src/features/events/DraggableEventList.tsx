'use client'
import React from 'react';
import { useAtom } from 'jotai';
import { selectedDatesAtom, type SelectedDate } from '@/atoms/calendar';
import DraggableList from '@/components/DraggableList';
import { formatJaDate } from '@/lib/date/format';

const DraggableEventList = () => {
    const [selectedDates, setSelectedDates] = useAtom(selectedDatesAtom);

    const handleReorder = (newDates: SelectedDate[]) => {
        setSelectedDates(newDates);
    }

    return (
        <DraggableList
            items={selectedDates}
            onReorder={handleReorder}
            renderItem={(date) => (
                <div>
                    {formatJaDate(date.start)} - {formatJaDate(date.end)}
                </div>
            )}
            getKey={(date) => date.id}
        />
    )
}

export default DraggableEventList;