import { useState } from 'react';
import { Meta, StoryObj } from '@storybook/nextjs';
import { PaginationControls } from './PaginationControls';

const meta: Meta<typeof PaginationControls> = {
    title: 'Common/PaginationControls',
    component: PaginationControls,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ['autodocs'],
};

export default meta;

type Story = StoryObj<typeof PaginationControls>;

const InteractiveTemplate = ({ total, limit }: { total: number; limit: number }) => {
    const [page, setPage] = useState(1);
    return <PaginationControls page={page} total={total} limit={limit} onPageChange={setPage} />;
};

export const Default: Story = {
    render: () => <InteractiveTemplate total={95} limit={20} />,
};

export const FewPages: Story = {
    render: () => <InteractiveTemplate total={45} limit={20} />,
};

export const SinglePage: Story = {
    render: () => <InteractiveTemplate total={8} limit={20} />,
};
