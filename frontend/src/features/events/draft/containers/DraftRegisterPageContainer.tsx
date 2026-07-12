import React from 'react';
import Link from 'next/link';
import EventDraft from '@/features/events/draft/components/EventDraft';
import { ChevronLeft } from 'lucide-react';

const DraftRegisterPageContainer = () => {
    return (
        <main className="mx-auto max-w-screen-2xl space-y-4 px-4 py-8 md:px-8">
            <Link
                href="/events"
                className="inline-flex items-center gap-1 text-sm text-muted-foreground transition-colors hover:text-foreground"
            >
                <ChevronLeft className="size-4" />
                イベント一覧へ
            </Link>
            <h1 className="text-2xl font-bold leading-snug tracking-normal text-gray-900">イベント作成</h1>
            <EventDraft />
        </main>
    );
};

export default DraftRegisterPageContainer;
