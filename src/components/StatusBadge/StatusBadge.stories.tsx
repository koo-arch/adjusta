import { Meta, StoryObj } from '@storybook/react';
import StatusBadge from './StatusBadge';

const meta: Meta<typeof StatusBadge> = {
    title: 'Components/StatusBadge',
    component: StatusBadge,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ['autodocs'],
    argTypes: {
        color: {
            options: ['gray', 'red', 'green', 'blue', 'yellow', 'indigo', 'purple', 'pink'],
            control: { type: 'radio' },
        },
        size: {
            options: ['sm', 'md', 'lg'],
            control: { type: 'radio' },
        },
    },
};

export default meta;

type Story = StoryObj<typeof StatusBadge>;

export const Default: Story = {
    args: {
        label: 'Status',
    },
};

export const CustomColor: Story = {
    args: {
        label: 'Status',
        color: 'green',
    },
};