import React, { Suspense } from 'react';
import EventDetailPageContainer from '@/features/events/detail/containers/EventDetailPageContainer';
import { EventDetailSkeleton } from '@/features/events/detail/components/EventDetail';

interface EventDetailPageProps {
    // Next.js 15+ では params は Promise
    params: Promise<{
        id: string;
    }>;
}

const EventDetailPageContent = async ({ params }: EventDetailPageProps) => {
    const { id } = await params;
    return <EventDetailPageContainer eventID={id} />;
};

const EventDetailPage = ({ params }: EventDetailPageProps) => (
    <Suspense
        fallback={
            <main className="mx-auto max-w-screen-md space-y-4 px-4 py-8">
                <EventDetailSkeleton />
            </main>
        }
    >
        <EventDetailPageContent params={params} />
    </Suspense>
);

export default EventDetailPage;
