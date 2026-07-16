import React from 'react';
import { Skeleton } from '@/components/ui/skeleton';

const EventFormSkeleton = () => (
    <div className="grid grid-cols-1 gap-8 md:grid-cols-10 md:gap-6">
        <div className="space-y-4 md:col-span-4">
            <Skeleton className="h-6 w-24" />
            <Skeleton className="h-9 w-full" />
            <Skeleton className="h-9 w-full" />
            <Skeleton className="h-24 w-full" />
        </div>
        <div className="md:col-span-6">
            <Skeleton className="h-96 w-full" />
        </div>
    </div>
);

export default EventFormSkeleton;
