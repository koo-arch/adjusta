import { Meta, StoryObj } from '@storybook/nextjs';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './card';

const meta: Meta<typeof Card> = {
    title: 'UI/Card',
    component: Card,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ['autodocs'],
};

export default meta;

type Story = StoryObj<typeof Card>;

export const Default: Story = {
    render: () => (
        <Card className="max-w-md">
            <CardHeader>
                <CardTitle>カードタイトル</CardTitle>
                <CardDescription>カードの補足説明が入ります。</CardDescription>
            </CardHeader>
            <CardContent>
                <p>カードの本文コンテンツが入ります。</p>
            </CardContent>
        </Card>
    ),
};
