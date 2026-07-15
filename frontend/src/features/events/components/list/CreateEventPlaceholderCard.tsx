import React from 'react';
import Link from 'next/link';
import { CREATE_EVENT_HREF, CREATE_EVENT_LABEL } from '@/features/events/constants';
import { Plus } from 'lucide-react';

// 疎なグリッドの末尾に置く破線の作成導線(frontend/DESIGN.md 2026-07-15)
const CreateEventPlaceholderCard: React.FC = () => (
    <Link
        href={CREATE_EVENT_HREF}
        className="flex h-full min-h-28 items-center justify-center gap-1 rounded-lg border border-dashed border-input text-sm text-muted-foreground transition-colors hover:border-primary hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
    >
        <Plus className="size-4" aria-hidden="true" />
        {CREATE_EVENT_LABEL}
    </Link>
);

export default CreateEventPlaceholderCard;
