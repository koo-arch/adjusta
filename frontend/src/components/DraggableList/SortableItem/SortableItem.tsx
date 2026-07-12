import React from 'react';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { cn } from '@/lib/utils';
import { GripVertical } from 'lucide-react';

interface SortableItemProps {
    id: string;
    children: React.ReactNode;
    disabled?: boolean;
}

// 行全体がドラッグ対象(マウスは即、タッチは長押しで発動)。
// グリップアイコンは「掴める」ことを示す目印として置く
const SortableItem: React.FC<SortableItemProps> = ({ id, children, disabled }) => {
    const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id, disabled });

    const style = {
        transform: CSS.Transform.toString(transform),
        transition: isDragging ? 'none' : transition,
    };

    return (
        <div
            ref={setNodeRef}
            style={style}
            className={cn(
                'mb-2 flex items-center gap-2 rounded-md border border-border bg-card px-3 py-2 shadow-sm',
                !disabled && 'cursor-grab',
                isDragging && 'z-10 cursor-grabbing opacity-80 shadow-md',
            )}
            {...attributes}
            {...(disabled ? {} : listeners)}
        >
            {!disabled && (
                <GripVertical aria-hidden className="size-4 shrink-0 text-muted-foreground" />
            )}
            <div className="min-w-0 flex-1">{children}</div>
        </div>
    );
};

export default SortableItem;
