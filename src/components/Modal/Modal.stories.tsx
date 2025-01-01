import { Meta, StoryObj } from '@storybook/react';
import Modal from './Modal';

const meta: Meta<typeof Modal> = {
    title: 'Components/Modal',
    component: Modal,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ['autodocs'],
};

export default meta;

type Story = StoryObj<typeof Modal>;

export const Default: Story = {
    args: {
        title: 'Modal title',
        description: 'Modal description',
        isOpen: true,
        onClose: () => {console.log('Modal closed')},
    },
};

export const WithActions: Story = {
    args: {
        title: 'Modal title',
        description: 'Modal description',
        isOpen: true,
        onClose: () => {console.log('Modal closed')},
        actions: <button>Action</button>,
    },
};