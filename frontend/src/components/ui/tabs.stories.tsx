import { Meta, StoryObj } from '@storybook/nextjs';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './tabs';

const meta: Meta<typeof Tabs> = {
    title: 'UI/Tabs',
    component: Tabs,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ['autodocs'],
};

export default meta;

type Story = StoryObj<typeof Tabs>;

export const Default: Story = {
    render: () => (
        <Tabs defaultValue="all">
            <TabsList>
                <TabsTrigger value="all">すべて</TabsTrigger>
                <TabsTrigger value="active">調整中</TabsTrigger>
                <TabsTrigger value="confirmed">確定</TabsTrigger>
            </TabsList>
            <TabsContent value="all">すべてのコンテンツ</TabsContent>
            <TabsContent value="active">調整中のコンテンツ</TabsContent>
            <TabsContent value="confirmed">確定のコンテンツ</TabsContent>
        </Tabs>
    ),
};
