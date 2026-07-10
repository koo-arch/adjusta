'use client'
import React from 'react';
import { arrayMove } from '@dnd-kit/sortable';
import DraggableList from '@/components/DraggableList';
import { Button } from '@/components/ui/button';
import { formatJaDateSpan } from '@/lib/date/format';
import { ChevronDown, ChevronUp, X } from 'lucide-react';
import type { ProposedDate, SelectedDate } from '@/features/events/store/dates';

interface DraggableDateListProps<T extends SelectedDate | ProposedDate> {
    dates: T[];
    onDatesChange: React.Dispatch<React.SetStateAction<T[]>>;
}

const DraggableDateList = <T extends SelectedDate | ProposedDate>({
    dates,
    onDatesChange,
 }: DraggableDateListProps<T>) => {
    const handleReorder = (newDates: T[]) => {
        onDatesChange(newDates);
    }

    const handleDelete = (id: string) => {
        onDatesChange(dates.filter(date => date.id !== id));
    }

    // ドラッグできない環境(キーボード・支援技術・モバイル)向けの代替手段
    const handleMove = (index: number, direction: -1 | 1) => {
        onDatesChange(arrayMove(dates, index, index + direction));
    }

    return (
        <DraggableList
            items={dates}
            onReorder={handleReorder}
            renderItem={(date: T, index) => (
                <div className="flex items-center justify-between gap-2">
                    <div className="flex min-w-0 items-center gap-2">
                        <span className="shrink-0 text-sm font-medium text-muted-foreground">{index + 1}.</span>
                        <span className="text-sm">{formatJaDateSpan(date.start, date.end)}</span>
                    </div>
                    <div className="flex shrink-0 items-center">
                        <Button
                            type="button"
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 text-muted-foreground hover:text-foreground"
                            aria-label="優先順位を上げる"
                            disabled={index === 0}
                            onClick={() => handleMove(index, -1)}
                        >
                            <ChevronUp />
                        </Button>
                        <Button
                            type="button"
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 text-muted-foreground hover:text-foreground"
                            aria-label="優先順位を下げる"
                            disabled={index === dates.length - 1}
                            onClick={() => handleMove(index, 1)}
                        >
                            <ChevronDown />
                        </Button>
                        <Button
                            type="button"
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 text-muted-foreground hover:text-destructive"
                            aria-label="この日程を削除"
                            onClick={() => handleDelete(date.id)}
                        >
                            <X />
                        </Button>
                    </div>
                </div>
            )}
            getKey={(date: T) => date.id}
        />
    )
}

export default DraggableDateList;
