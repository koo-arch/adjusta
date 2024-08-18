'use client'
import React, { useState, useRef } from 'react';
import { useAtom, useAtomValue } from 'jotai';
import { EventClickArg } from '@fullcalendar/core';
import { selectedDatesAtom, selectedEventsAtom } from '@/atoms/calendar';
import Calendar from './Calendar';
import { EventImpl } from '@fullcalendar/core/internal';
import PopupMenu from '@/components/PopupMenu';

const SelectableCalendar: React.FC = () => {
    const [selectedDates, setSelectedDates] = useAtom(selectedDatesAtom);
    const selectedEvents = useAtomValue(selectedEventsAtom);
    const [clickedEvent, setClickedEvent] = useState<EventImpl>();
    const [popupPosition, setPopupPosition] = useState({ top: 0, left: 0 });
    const buttonRef = useRef<HTMLButtonElement>(null);

    const handleDateSelect = (selectInfo: { start: Date, end: Date }) => {
        const newDate = {
            id : new Date().getTime().toString(),
            start: selectInfo.start,
            end: selectInfo.end,
        }
        setSelectedDates([...selectedDates, newDate]);
    }

    const handleEventClick = (e: EventClickArg) => {
        if (buttonRef.current) {
            buttonRef.current.click();
        }
        setClickedEvent(e.event);
        setPopupPosition({ top: e.jsEvent.clientY, left: e.jsEvent.clientX });
    }

    const handleDeleteEvent = () => {
        if (clickedEvent) {
            clickedEvent.remove();
            setSelectedDates((prev) => prev.filter((date) => date.id !== clickedEvent.id));
        }
    }

    return (
        <div>
            <Calendar 
                initialView="timeGridWeek"
                headerToolbar={{
                    left: 'prev,next today',
                    center: 'title',
                    right: '',
                }}
                select={handleDateSelect}
                selectedEvents={selectedEvents}
                handleEventClick={handleEventClick}
            />
            <PopupMenu
                items={[
                    { label: '削除', onClick: handleDeleteEvent },
                ]}
                position={popupPosition}
                buttonRef={buttonRef}
            />
        </div>
    )
}

export default SelectableCalendar;