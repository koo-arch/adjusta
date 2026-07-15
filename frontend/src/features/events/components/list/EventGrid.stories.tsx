import { Meta, StoryObj } from '@storybook/nextjs';
import EventGrid from './EventGrid';
import type { EventDraftDetail, EventProposedDate, EventStatus } from '@/features/events/types';

const makeDate = (id: string, start: string, end: string, status: EventProposedDate['status'] = 'active'): EventProposedDate => ({
    id,
    start: new Date(start),
    end: new Date(end),
    priority: 1,
    status,
    sync_status: 'not_synced',
});

const makeEvent = (
    id: string,
    title: string,
    status: EventStatus,
    proposedDates: EventProposedDate[] = [],
    confirmedDateId: string | null = null,
): EventDraftDetail => ({
    id,
    title,
    description: '',
    location: '',
    status,
    sync_status: 'not_synced',
    confirmed_date_id: confirmedDateId,
    proposed_dates: proposedDates,
});

const events: EventDraftDetail[] = [
    makeEvent('1', 'チーム定例ミーティング', 'active', [
        makeDate('d1', '2026-07-21T10:00:00', '2026-07-21T11:00:00'),
        makeDate('d2', '2026-07-22T14:00:00', '2026-07-22T15:00:00'),
        makeDate('d3', '2026-07-23T09:00:00', '2026-07-23T10:00:00'),
    ]),
    makeEvent('2', '歓迎会', 'confirmed', [
        makeDate('d4', '2026-07-24T19:00:00', '2026-07-24T21:00:00', 'confirmed'),
    ], 'd4'),
    makeEvent('3', '合宿の日程調整', 'draft'),
    makeEvent('4', 'デザインレビュー', 'active', [
        makeDate('d5', '2026-07-27T13:00:00', '2026-07-27T14:00:00'),
    ]),
    makeEvent('5', '打ち上げ', 'cancelled'),
    makeEvent('6', '四半期キックオフ', 'confirmed', [
        makeDate('d6', '2026-08-03T10:00:00', '2026-08-03T12:00:00', 'confirmed'),
    ], 'd6'),
];

const meta: Meta<typeof EventGrid> = {
    title: 'Events/EventGrid',
    component: EventGrid,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ['autodocs'],
};

export default meta;

type Story = StoryObj<typeof EventGrid>;

export const OneItemWithPlaceholder: Story = {
    args: {
        events: events.slice(0, 1),
        showCreatePlaceholder: true,
    },
};

export const Full: Story = {
    args: {
        events,
        showCreatePlaceholder: false,
    },
};

export const Empty: Story = {
    args: {
        events: [],
        showCreatePlaceholder: true,
    },
};
