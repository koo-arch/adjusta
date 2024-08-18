import { Meta, StoryObj } from '@storybook/react';
import PopupMenu from './PopupMenu';

const meta: Meta<typeof PopupMenu> = {
    title: 'Components/PopupMenu',
    component: PopupMenu,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ['autodocs'],
    argTypes: {
        position: {
            control: 'object',
            description: 'PopupMenuの表示位置',
        },
    },
};

export default meta;

type Story = StoryObj<typeof PopupMenu>;

export const Default: Story = {
    args: {
        items: [
            {
                label: 'Item 1',
                onClick: () => console.log('Item 1 clicked'),
            },
            {
                label: 'Item 2',
                onClick: () => console.log('Item 2 clicked'),
            },
            {
                label: 'Item 3',
                onClick: () => console.log('Item 3 clicked'),
            },

        ],
        position: { top: 50, left: 1000 },
    }
};