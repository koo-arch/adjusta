import { Meta, StoryObj } from '@storybook/nextjs';
import { Button } from './button';

const meta: Meta<typeof Button> = {
    title: 'UI/Button',
    component: Button,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ['autodocs'],
    argTypes: {
        variant: {
            options: ['default', 'destructive', 'outline', 'secondary', 'ghost', 'link'],
            control: { type: 'select' },
        },
        size: {
            options: ['default', 'sm', 'lg', 'icon'],
            control: { type: 'radio' },
        },
    },
};

export default meta;

type Story = StoryObj<typeof Button>;

export const Default: Story = {
    args: {
        children: 'ボタン',
    },
};

export const Outline: Story = {
    args: {
        variant: 'outline',
        children: 'ボタン',
    },
};

export const Destructive: Story = {
    args: {
        variant: 'destructive',
        children: '削除',
    },
};

export const Disabled: Story = {
    args: {
        disabled: true,
        children: 'ボタン',
    },
};
