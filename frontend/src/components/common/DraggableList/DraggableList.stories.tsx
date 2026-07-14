import { useState } from 'react';
import type { Meta, StoryObj } from '@storybook/nextjs';
import DraggableList from './DraggableList';

type Item = { id: string; label: string };

const Example = () => {
    const [items, setItems] = useState<Item[]>([
        { id: '1', label: '第1候補' },
        { id: '2', label: '第2候補' },
        { id: '3', label: '第3候補' },
    ]);

    return (
        <DraggableList
            items={items}
            onReorder={setItems}
            getKey={(item) => item.id}
            renderItem={(item) => item.label}
        />
    );
};

const meta = {
    title: 'Components/Common/DraggableList',
    component: Example,
    tags: ['autodocs'],
} satisfies Meta<typeof Example>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
