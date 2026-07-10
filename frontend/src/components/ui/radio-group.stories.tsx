import { Meta, StoryObj } from '@storybook/nextjs';
import { RadioGroup, RadioGroupItem } from './radio-group';

const meta: Meta<typeof RadioGroup> = {
    title: 'UI/RadioGroup',
    component: RadioGroup,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
};

export default meta;

type Story = StoryObj<typeof RadioGroup>;

export const Default: Story = {
    render: () => (
        <RadioGroup defaultValue="a">
            <div className="flex items-center gap-2">
                <RadioGroupItem value="a" id="option-a" />
                <label htmlFor="option-a">選択肢 A</label>
            </div>
            <div className="flex items-center gap-2">
                <RadioGroupItem value="b" id="option-b" />
                <label htmlFor="option-b">選択肢 B</label>
            </div>
            <div className="flex items-center gap-2">
                <RadioGroupItem value="c" id="option-c" disabled />
                <label htmlFor="option-c">選択肢 C(無効)</label>
            </div>
        </RadioGroup>
    ),
};
