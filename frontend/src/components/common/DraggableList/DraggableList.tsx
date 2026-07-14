'use client'

import React from 'react';
import {
    DndContext,
    KeyboardSensor,
    MouseSensor,
    TouchSensor,
    closestCenter,
    type DragEndEvent,
    useSensor,
    useSensors,
} from '@dnd-kit/core';
import { CSS } from '@dnd-kit/utilities';
import {
    SortableContext,
    arrayMove,
    useSortable,
    verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { GripVertical } from 'lucide-react';
import { cn } from '@/lib/utils';

interface SortableItemProps {
    id: string;
    children: React.ReactNode;
    disabled?: boolean;
}

const SortableItem: React.FC<SortableItemProps> = ({ id, children, disabled }) => {
    const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id, disabled });

    return (
        <div
            ref={setNodeRef}
            style={{
                transform: CSS.Transform.toString(transform),
                transition: isDragging ? 'none' : transition,
            }}
            className={cn(
                'mb-2 flex items-center gap-2 rounded-md border border-border bg-card px-3 py-2 shadow-sm',
                !disabled && 'cursor-grab',
                isDragging && 'z-10 cursor-grabbing opacity-80 shadow-md',
            )}
            {...(disabled ? {} : attributes)}
            {...(disabled ? {} : listeners)}
        >
            {!disabled && <GripVertical aria-hidden className="size-4 shrink-0 text-muted-foreground" />}
            <div className="min-w-0 flex-1">{children}</div>
        </div>
    );
};

interface DraggableListProps<T> {
    items: T[];
    onReorder: (newItems: T[]) => void;
    renderItem: (item: T, index: number) => React.ReactNode;
    getKey: (item: T) => string;
    disabledIds?: string[];
}

const DraggableList = <T,>({
    items,
    onReorder,
    renderItem,
    getKey,
    disabledIds,
}: DraggableListProps<T>) => {
    const handleDragEnd = ({ active, over }: DragEndEvent) => {
        if (!over || active.id === over.id) return;

        const oldIndex = items.findIndex((item) => getKey(item) === active.id);
        const newIndex = items.findIndex((item) => getKey(item) === over.id);
        if (oldIndex >= 0 && newIndex >= 0) {
            onReorder(arrayMove(items, oldIndex, newIndex));
        }
    };

    const sensors = useSensors(
        useSensor(MouseSensor, { activationConstraint: { distance: 5 } }),
        useSensor(TouchSensor, { activationConstraint: { delay: 250, tolerance: 8 } }),
        useSensor(KeyboardSensor),
    );

    return (
        <DndContext collisionDetection={closestCenter} onDragEnd={handleDragEnd} sensors={sensors}>
            <SortableContext items={items.map(getKey)} strategy={verticalListSortingStrategy}>
                {items.map((item, index) => (
                    <SortableItem
                        key={getKey(item)}
                        id={getKey(item)}
                        disabled={disabledIds?.includes(getKey(item))}
                    >
                        {renderItem(item, index)}
                    </SortableItem>
                ))}
            </SortableContext>
        </DndContext>
    );
};

export default DraggableList;
