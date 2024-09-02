import React from 'react';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';

interface SortableItemProps {
    id: string;
    children: React.ReactNode;
}

const SortableItem: React.FC<SortableItemProps> = ({ id, children }) => {
    const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id });

    const style = {
        transform: CSS.Transform.toString(transform),
        transition: isDragging ? 'none' : transition,
    };

    return (
        <div ref={setNodeRef}
            style={style}
            className="p-4 mb-3 text-white bg-indigo-400 border-indigo-300 shadow-md hover:bg-indigo-700 transition duration-300 ease-in-out cursor-grab"
            {...attributes}
            {...listeners}
        >
            {children}
        </div>
    );
};

export default SortableItem;