'use client'
import React from 'react';
import Link from 'next/link';
import EventSearchForm from './EventSearchForm';
import { Button } from '@/components/ui/button';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { CREATE_EVENT_HREF, CREATE_EVENT_LABEL } from '@/features/events/constants';
import { Plus } from 'lucide-react';

export interface ToolbarTab {
    value: string;
    label: string;
}

interface EventsToolbarProps {
    tabs: ToolbarTab[];
    activeTab: string;
    onTabChange: (value: string) => void;
    searchValue: string;
    onSearch: (value: string) => void;
}

// フィルタ・検索・作成を1行に集約するツールバー(frontend/DESIGN.md 2026-07-15)。
// md 未満では縦積みになり、作成は FAB 側が担う
const EventsToolbar: React.FC<EventsToolbarProps> = ({
    tabs,
    activeTab,
    onTabChange,
    searchValue,
    onSearch,
}) => (
    <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <Tabs value={activeTab} onValueChange={onTabChange}>
            {/* 画面端までブリードさせ、フィルタに続きがあることを示す */}
            <div className="overflow-x-auto -mx-4 px-4 sm:-mx-6 sm:px-6 md:mx-0 md:px-0">
                <TabsList>
                    {tabs.map((tab) => (
                        <TabsTrigger key={tab.value} value={tab.value}>
                            {tab.label}
                        </TabsTrigger>
                    ))}
                </TabsList>
            </div>
        </Tabs>
        <div className="flex items-center gap-2">
            <EventSearchForm defaultValue={searchValue} onSearch={onSearch} />
            <Button asChild className="hidden shrink-0 md:inline-flex">
                <Link href={CREATE_EVENT_HREF}>
                    <Plus aria-hidden="true" />
                    {CREATE_EVENT_LABEL}
                </Link>
            </Button>
        </div>
    </div>
);

export default EventsToolbar;
