import React from 'react';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { cn } from '@/lib/utils';
import { GripVertical } from 'lucide-react';

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
        <div
            ref={setNodeRef}
            style={style}
            className={cn(
                'mb-2 flex items-center gap-2 rounded-md border border-border bg-card px-3 py-2 shadow-sm',
                isDragging && 'z-10 opacity-80 shadow-md',
            )}
        >
            {/* ドラッグはハンドルに限定する(行内のボタン操作・モバイルのスクロールと衝突させない) */}
            <button
                type="button"
                aria-label="ドラッグして並べ替え"
                className="shrink-0 cursor-grab touch-none rounded-sm p-1 text-muted-foreground hover:text-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                {...attributes}
                {...listeners}
            >
                <GripVertical className="size-4" />
            </button>
            <div className="min-w-0 flex-1">{children}</div>
        </div>
    );
};

export default SortableItem;
