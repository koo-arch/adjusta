'use client'
import React from 'react';
import { useAtom } from 'jotai';
import type { SelectedDate, ProposedDate } from '@/atoms/calendar';
import DraggableList from '@/components/DraggableList';
import WrapText from '@/components/WrapText';
import { formatJaDateSpan } from '@/lib/date/format';
import IconButton from '@/components/IconButton';
import { TrashIcon } from '@heroicons/react/20/solid';

interface DraggableEventListProps {
    atom: any;
    enableTopHighlight?: boolean;
}

const DraggableEventList = <T extends SelectedDate | ProposedDate>({
    atom,
    enableTopHighlight
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
            renderItem={(date: T, index) => (
                <div className="flex justify-between items-center">
                    <span>
                        {index + 1}.
                    </span>
                    <WrapText 
                        text={formatJaDateSpan(date.start, date.end)}
                        maxLength={23}
                        marker='-'
                    />
                    <IconButton
                        type="button"
                        onClick={() => {
                            handleDelete(date.id);
                            console.log('delete');
                        }}
                        className="ml-1"
                        iconSize="md"
                        iconColor="clear"
                    >
                        <TrashIcon />
                    </IconButton>
                </div>
            )}
            getKey={(date: T) => date.id}
            enableTopHighlight={enableTopHighlight}
        />
    )
}

export default DraggableEventList;