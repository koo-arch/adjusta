import { Meta, StoryObj } from '@storybook/react';
import Card from './Card';

const meta: Meta<typeof Card> = {
    title: 'Components/Card',
    component: Card,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ['autodocs'],
    argTypes: {
        variant: {
            options: ['shadow', 'outlined'],
            control: { type: 'radio' },
        },
    },
};

export default meta;

type Story = StoryObj<typeof Card>;

const cardContent = (
    <div>
        <h1 className="text-xl font-bold">Card title</h1>
        <p className="text-sm text-gray-500">Card content</p>
    </div>
);

export const Default: Story = {
    args: {
        children: cardContent,
        onClick: () => {console.log('Card clicked')},
    },
};

export const WithActions: Story = {
    args: {
        children: cardContent,
        actions: <button>Action</button>,
    },
};