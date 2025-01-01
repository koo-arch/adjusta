import React from 'react';
import type { DayCellContentArg, DayHeaderContentArg, SlotLabelMountArg } from '@fullcalendar/core';
import { format } from 'date-fns';
import { CircleNumber, CircleToday } from './style';

export const renderDayCell = (e: DayCellContentArg) => {
    const { date, dayNumberText, isToday } = e
    const replaceDayNumberText = dayNumberText.replace('日', '')

    return dayNumberText && isToday ? (
        <CircleNumber>{ replaceDayNumberText } </CircleNumber>
    ) : dayNumberText === '1日' ? (
        <>{ format(date, 'M月d日') } </>
    ) : (
        <>{ replaceDayNumberText } </>
    )
}

export const renderDayHeader = (e: DayHeaderContentArg) => {
    const { text, isToday, view } = e
    if (view.type === 'dayGridMonth') {
        return text
    }

    if (isToday) {
        return (
            <CircleToday>{ text } </CircleToday>
        )
    }
    return text
}

export const renderSlotLabel = (e: SlotLabelMountArg) => {
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