'use client'
import React from 'react';
import { useAtom } from 'jotai';
import type { SelectedDate, ProposedDate } from '@/atoms/calendar';
import DraggableList from '@/components/DraggableList';
import { formatJaDate } from '@/lib/date/format';
import IconButton from '@/components/IconButton';
import { TrashIcon } from '@heroicons/react/20/solid';

interface DraggableEventListProps {
    atom: any;
}

const DraggableEventList = <T extends SelectedDate | ProposedDate>({
    atom,
 }: DraggableEventListProps) => {
    const [dates, setDates] = useAtom<T[]>(atom);

    const handleReorder = (newDates: T[]) => {
        setDates(newDates);
    }

    const handleDelete = (id: string) => {
        setDates(dates.filter(date => date.id !== id));
    }

    return (
        <DraggableList
            items={dates}
            onReorder={handleReorder}
            renderItem={(date: T) => (
                <div className="flex justify-between items-center">
                    <span>
                        {formatJaDate(date.start)} 〜 {formatJaDate(date.end)}
                    </span>
                    <IconButton
                        onClick={() => handleDelete(date.id)}
                        className="ml-2"
                        iconSize="md"
                        iconColor="clear"
                    >
                        <TrashIcon />
                    </IconButton>
                </div>
            )}
            getKey={(date: T) => date.id}
        />
    )
}

export default DraggableEventList;