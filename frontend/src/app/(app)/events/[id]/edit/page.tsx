import React, { Suspense } from 'react';
import EventEditPageContainer from '@/features/events/edit/containers/EventEditPageContainer';
import EventFormSkeleton from '@/features/events/components/form/EventFormSkeleton';

interface EventEditPageProps {
    // Next.js 15+ では params は Promise
    params: Promise<{
        id: string;
    }>;
}

const EventEditPageContent = async ({ params }: EventEditPageProps) => {
    const { id } = await params;
    return <EventEditPageContainer eventID={id} />;
}

const EventEditPage = ({ params }: EventEditPageProps) => (
    <Suspense
        fallback={
            <main className="mx-auto max-w-screen-2xl space-y-4 px-4 py-8 md:px-8">
                <h1 className="text-2xl font-bold leading-snug tracking-normal text-gray-900">
                    イベント編集
                </h1>
                <EventFormSkeleton />
            </main>
        }
    >
        <EventEditPageContent params={params} />
    </Suspense>
);

export default EventEditPage;
