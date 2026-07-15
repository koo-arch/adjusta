import type { Meta, StoryObj } from '@storybook/nextjs';
import StatusBadge from './StatusBadge';

const meta = {
    title: 'Components/Common/StatusBadge',
    component: StatusBadge,
    tags: ['autodocs'],
    args: {
        label: '調整中',
        color: 'yellow',
        dotSize: 'md',
        textSize: 'sm',
    },
} satisfies Meta<typeof StatusBadge>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Failed: Story = {
    args: {
        label: '同期失敗',
        color: 'red',
    },
};
