'use client'
import React from 'react';
import { DndContext, closestCenter, MouseSensor, KeyboardSensor, useSensor, useSensors, type DragEndEvent } from '@dnd-kit/core';
import { arrayMove, SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import SortableItem from './SortableItem';

interface DraggableListProps<T> {
    items: T[];
    onReorder: (newItems: T[]) => void;
    renderItem: (item: T) => React.ReactNode;
    getKey: (item: T) => string;
}

const DraggableList = <T extends unknown>({
    items,
    onReorder,
    renderItem,
    getKey,
}: DraggableListProps<T>) => {

    // どの要素がドラッグされているかを管理するための変数
    const handleDragEnd = (event: DragEndEvent) => {
        const { active, over } = event;

        // ドラッグされている要素がない場合、何もしない
        if (!over || active.id === over.id) return;

        const oldIndex = items.findIndex((item) => getKey(item) === active.id);
        const newIndex = items.findIndex((item) => getKey(item) === over.id);

        // ドラッグされた要素を新しい位置に移動
        if (oldIndex >= 0 && newIndex >= 0) {
            const newItems = arrayMove(items, oldIndex, newIndex);
            onReorder(newItems);
        }
    }

    const mouseSensor = useSensor(MouseSensor, {
        activationConstraint: {
            distance: 5,
        },
    });
    const keyboardSensor = useSensor(KeyboardSensor);
    const sensors = useSensors(mouseSensor, keyboardSensor);


    return (
        <DndContext collisionDetection={closestCenter} onDragEnd={handleDragEnd} sensors={sensors}>
            <SortableContext items={items.map(getKey)} strategy={verticalListSortingStrategy}>
                {items.map((item) => (
                    <SortableItem key={getKey(item)} id={getKey(item)}>
                        {renderItem(item)}
                    </SortableItem>
                ))}
            </SortableContext>
        </DndContext>
    );
};

export default DraggableList;