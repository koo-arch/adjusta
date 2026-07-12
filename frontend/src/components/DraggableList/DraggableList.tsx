'use client'
import React from 'react';
import { DndContext, closestCenter, MouseSensor, TouchSensor, KeyboardSensor, useSensor, useSensors, type DragEndEvent } from '@dnd-kit/core';
import { arrayMove, SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import SortableItem from './SortableItem';

interface DraggableListProps<T> {
    items: T[];
    onReorder: (newItems: T[]) => void;
    renderItem: (item: T, index: number) => React.ReactNode;
    getKey: (item: T) => string;
    // ドラッグを無効にする行(例: インライン編集中の行)
    disabledIds?: string[];
}

const DraggableList = <T extends unknown>({
    items,
    onReorder,
    renderItem,
    getKey,
    disabledIds,
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

    // マウスは 5px 移動で即ドラッグ、タッチは長押し(250ms)で発動しスクロールと衝突させない
    const mouseSensor = useSensor(MouseSensor, {
        activationConstraint: {
            distance: 5,
        },
    });
    const touchSensor = useSensor(TouchSensor, {
        activationConstraint: {
            delay: 250,
            tolerance: 8,
        },
    });
    const keyboardSensor = useSensor(KeyboardSensor);
    const sensors = useSensors(mouseSensor, touchSensor, keyboardSensor);


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
