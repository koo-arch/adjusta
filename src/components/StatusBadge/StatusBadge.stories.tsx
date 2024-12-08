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
        label: {
            control: 'text',
        },
        circleSize: {
            options: ['sm', 'md', 'lg'],
            control: { type: 'radio' },
        },
        circleColor: {
            options: ['gray', 'red', 'green', 'blue', 'yellow', 'indigo', 'purple', 'pink'],
            control: { type: 'radio' },
        },
        textSize: {
            options: ['sm', 'md', 'lg', 'xl'],
            control: { type: 'radio' },
        },
        textColor: {
            options: ['gray', 'red', 'green', 'blue', 'yellow', 'indigo', 'purple', 'pink'],
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
        circleColor: 'green',
        circleSize: 'lg',
        textColor: 'green',
        textSize: 'xl',
    },
};