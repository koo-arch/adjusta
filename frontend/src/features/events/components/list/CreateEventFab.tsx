import React from 'react';
import Link from 'next/link';
import { buttonVariants } from '@/components/ui/button';
import { CREATE_EVENT_HREF, CREATE_EVENT_LABEL } from '@/features/events/constants';
import { cn } from '@/lib/utils';
import { Plus } from 'lucide-react';

// モバイル専用の作成 FAB。md 以上ではツールバーの primary ボタンが担う
const CreateEventFab: React.FC = () => (
    <Link
        href={CREATE_EVENT_HREF}
        aria-label={CREATE_EVENT_LABEL}
        className={cn(
            buttonVariants(),
            'fixed right-4 bottom-[calc(1rem+env(safe-area-inset-bottom))] z-40 md:hidden',
            // CVA base の [&_svg]:size-4 より優先させるためリンク側で上書きする
            'size-14 rounded-full p-0 shadow-lg [&_svg]:size-6',
        )}
    >
        <Plus aria-hidden="true" />
    </Link>
);

export default CreateEventFab;
