import { Meta, StoryObj } from '@storybook/nextjs';
import { fn } from 'storybook/test';
import EventsToolbar from './EventsToolbar';

const STATUS_TABS = [
    { value: 'all', label: 'すべて' },
    { value: 'active', label: '調整中' },
    { value: 'confirmed', label: '確定' },
    { value: 'draft', label: '下書き' },
    { value: 'cancelled', label: 'キャンセル' },
];

const meta: Meta<typeof EventsToolbar> = {
    title: 'Events/EventsToolbar',
    component: EventsToolbar,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    args: {
        tabs: STATUS_TABS,
        activeTab: 'all',
        searchValue: '',
        onTabChange: fn(),
        onSearch: fn(),
    },
    tags: ['autodocs'],
};

export default meta;

type Story = StoryObj<typeof EventsToolbar>;

export const Desktop: Story = {};

export const MobileStacked: Story = {
    globals: {
        viewport: { value: 'mobile1', isRotated: false },
    },
};

export const ManyFiltersOverflow: Story = {
    args: {
        tabs: [
            ...STATUS_TABS,
            { value: 'archived', label: 'アーカイブ' },
            { value: 'recurring', label: '定期' },
            { value: 'shared', label: '共有' },
            { value: 'past', label: '過去' },
        ],
    },
    globals: {
        viewport: { value: 'mobile1', isRotated: false },
    },
};
