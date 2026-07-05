'use client'
import React from 'react';
import DraggableList from '@/components/DraggableList';
import WrapText from '@/components/WrapText';
import { formatJaDateSpan } from '@/lib/date/format';
import IconButton from '@/components/IconButton';
import { TrashIcon } from '@heroicons/react/20/solid';
import type { ProposedDate, SelectedDate } from '@/features/events/store/dates';

interface DraggableDateListProps<T extends SelectedDate | ProposedDate> {
    dates: T[];
    onDatesChange: React.Dispatch<React.SetStateAction<T[]>>;
    enableTopHighlight?: boolean;
}

const DraggableDateList = <T extends SelectedDate | ProposedDate>({
    dates,
    onDatesChange,
    enableTopHighlight
 }: DraggableDateListProps<T>) => {
    const handleReorder = (newDates: T[]) => {
        onDatesChange(newDates);
    }

    const handleDelete = (id: string) => {
        onDatesChange(dates.filter(date => date.id !== id));
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

export default DraggableDateList;
