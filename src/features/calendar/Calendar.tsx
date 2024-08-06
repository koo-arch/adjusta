'use client'
import React from 'react';
import { format } from 'date-fns';
import { StyleWrapper, CircleNumber, CircleToday } from './style';
import FullCalendar from '@fullcalendar/react';
import { DayCellContentArg, DayHeaderContentArg, SlotLabelMountArg } from '@fullcalendar/core';
import dayGridPlugin from '@fullcalendar/daygrid';
import interactionPlugin from '@fullcalendar/interaction';
import timeGridPlugin from '@fullcalendar/timegrid';
import jaLocale from '@fullcalendar/core/locales/ja';
import momentPlugin from '@fullcalendar/moment';
import { useFetchEvent } from '@/hooks/calendar/useFetchEvent';


const Calendar = () => {
    const { events, isLoading, error } = useFetchEvent();

    if (isLoading) return <div>Loading...</div>;

    const eventList = events?.flatMap((account) => {
        return account.events.map((event) => {
            return {
                title: event.summary,
                start: event.start,
                end: event.end,
            }
        })
    })

    const renderDayCell = (e: DayCellContentArg) => {
        const { date, dayNumberText, isToday } = e
        const replaceDayNumberText = dayNumberText.replace('日', '')

        return dayNumberText && isToday ? (
            <CircleNumber>{replaceDayNumberText}</CircleNumber>
        ) : dayNumberText === '1日' ? (
            <>{format(date, 'M月d日')}</>
        ) : (
            <>{replaceDayNumberText}</>
        )
    }

    const renderDayHeader = (e: DayHeaderContentArg) => {
        const { text, isToday, view } = e
        if (view.type === 'dayGridMonth') {
            return text
        }

        if (isToday) {
            return (
                <CircleToday>{text}</CircleToday>
            )
        }
        return text
    }

    const renderSlotLabel = (e: SlotLabelMountArg) => {
        const { date, view } = e

        if (view.type === 'dayGridMonth') {
            return
        }

        let hhmm = format(date, 'HH:mm')
        if (hhmm[0] === '0') {
            hhmm = hhmm.slice(1)
        }

        return hhmm
    }
    return (
    <StyleWrapper>
        <FullCalendar
            plugins={[dayGridPlugin, timeGridPlugin, interactionPlugin, momentPlugin]}
            initialView="dayGridMonth"
            headerToolbar={{
                left: 'prev,next today',
                center: 'title',
                right: 'dayGridMonth,timeGridWeek,timeGridDay'
            }}
            businessHours={{ daysOfWeek: [1, 2, 3, 4, 5] }}
            dayCellContent={renderDayCell}
            dayHeaderContent={renderDayHeader}
            slotLabelContent={renderSlotLabel}
            aspectRatio={1.6}
            locale={jaLocale}
            events={eventList}
        />
    </StyleWrapper>
    );
}

export default Calendar;