import { Meta, StoryObj } from '@storybook/nextjs';
import CreateEventFab from './CreateEventFab';

const meta: Meta<typeof CreateEventFab> = {
    title: 'Events/CreateEventFab',
    component: CreateEventFab,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    globals: {
        // FAB は md:hidden のためモバイルビューポートで表示する
        viewport: { value: 'mobile1', isRotated: false },
    },
    decorators: [
        (Story) => (
            // transform を持つ祖先を作り、fixed 配置をフレーム内に閉じ込める
            <div className="relative min-h-[400px]" style={{ transform: 'translateZ(0)' }}>
                <Story />
            </div>
        ),
    ],
    tags: ['autodocs'],
};

export default meta;

type Story = StoryObj<typeof CreateEventFab>;

export const Default: Story = {};

export const FocusRing: Story = {
    play: async ({ canvasElement }) => {
        canvasElement.querySelector('a')?.focus();
    },
};
