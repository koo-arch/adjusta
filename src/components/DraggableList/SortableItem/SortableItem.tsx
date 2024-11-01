import React from 'react';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';

interface SortableItemProps {
    id: string;
    children: React.ReactNode;
    index?: number;
    enableTopHighlight?: boolean;
}

const SortableItem: React.FC<SortableItemProps> = ({ id, children, index, enableTopHighlight }) => {
    const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id });

    const style = {
        transform: CSS.Transform.toString(transform),
        transition: isDragging ? 'none' : transition,
    };

    const baseStyle = `p-4 mb-3 text-white shadow-md transition duration-300 ease-in-out cursor-grab`
    const highlightStyle = enableTopHighlight && index === 0 
        ? 'bg-orange-400 border-orange-300 hover:bg-orange-700' 
        : 'bg-indigo-400 border-indigo-300 hover:bg-indigo-700';
    

    return (
        <div ref={setNodeRef}
            style={style}
            className={`${baseStyle} ${highlightStyle}`}
            {...attributes}
            {...listeners}
        >
            {children}
        </div>
    );
};

export default SortableItem;