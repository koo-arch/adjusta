'use client'
import React from 'react';
import Link from 'next/link';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Skeleton } from '@/components/ui/skeleton';
import { Switch } from '@/components/ui/switch';
import { useAccounts } from '@/features/auth/hooks/useAccounts';
import { useCalendarSettings } from '@/features/auth/hooks/useCalendarSettings';
import { useUpdateCalendarSetting } from '@/features/auth/hooks/useUpdateCalendarSetting';
import { isGoogleReauthorizationRequiredError } from '@/lib/api/errors';
import type { CalendarSetting, UserCalendarRole } from '@/features/auth/types';

const roleBadges: Record<UserCalendarRole, { label: string; className: string }> = {
    primary: {
        label: 'メイン',
        className: 'border-transparent bg-primary text-primary-foreground hover:bg-primary',
    },
    adjusta_candidate: {
        label: 'Adjusta 候補用',
        className: 'border-primary/40 bg-transparent text-primary hover:bg-transparent',
    },
    reference: {
        label: '参照',
        className: 'border-transparent bg-secondary text-muted-foreground hover:bg-secondary',
    },
};

// Adjusta 専用カレンダー(自動管理)は一覧の末尾にまとめ、それ以外は名前順で安定させる
const sortSettings = (settings: CalendarSetting[]): CalendarSetting[] =>
    [...settings].sort((a, b) => {
        const aCandidate = a.role === 'adjusta_candidate' ? 1 : 0;
        const bCandidate = b.role === 'adjusta_candidate' ? 1 : 0;
        if (aCandidate !== bCandidate) {
            return aCandidate - bCandidate;
        }
        return a.summary.localeCompare(b.summary, 'ja');
    });

interface CalendarSettingRowProps {
    setting: CalendarSetting;
    disabled: boolean;
    onToggleVisible: (checked: boolean) => void;
    onToggleSync: (checked: boolean) => void;
}

const CalendarSettingRow: React.FC<CalendarSettingRowProps> = ({
    setting,
    disabled,
    onToggleVisible,
    onToggleSync,
}) => {
    const isCandidate = setting.role === 'adjusta_candidate';
    const badge = roleBadges[setting.role];

    return (
        <li className="flex flex-col gap-3 py-4 sm:flex-row sm:items-start sm:justify-between">
            <div className="flex min-w-0 items-start gap-3">
                <RadioGroupItem
                    value={setting.id}
                    id={`primary-${setting.id}`}
                    disabled={disabled || isCandidate}
                    aria-label={`「${setting.summary}」を確定予定の登録先にする`}
                    className="mt-1 shrink-0"
                />
                <div className="min-w-0">
                    <div className="flex flex-wrap items-center gap-2">
                        <label
                            htmlFor={`primary-${setting.id}`}
                            className="break-words font-medium"
                        >
                            {setting.summary}
                        </label>
                        <Badge className={badge.className}>{badge.label}</Badge>
                    </div>
                    {setting.timezone && (
                        <p className="mt-0.5 text-sm text-muted-foreground">{setting.timezone}</p>
                    )}
                    {isCandidate && (
                        <p className="mt-1 text-sm text-muted-foreground">
                            Adjusta が候補日程の仮予定用に自動作成・管理するカレンダーです。名称の変更・削除や、確定予定の登録先への指定はできません。
                        </p>
                    )}
                </div>
            </div>
            <div className="flex shrink-0 flex-col gap-3 pl-7 sm:pl-0">
                <div className="flex items-center justify-between gap-3 sm:justify-end">
                    <label htmlFor={`visible-${setting.id}`} className="text-sm">
                        アプリに表示
                    </label>
                    <Switch
                        id={`visible-${setting.id}`}
                        checked={setting.is_visible}
                        disabled={disabled}
                        onCheckedChange={onToggleVisible}
                    />
                </div>
                {isCandidate && (
                    <div className="sm:max-w-64">
                        <div className="flex items-center justify-between gap-3 sm:justify-end">
                            <label htmlFor={`sync-${setting.id}`} className="text-sm">
                                候補日程を同期
                            </label>
                            <Switch
                                id={`sync-${setting.id}`}
                                checked={setting.sync_proposed_dates}
                                disabled={disabled}
                                onCheckedChange={onToggleSync}
                            />
                        </div>
                        <p className="mt-1 text-xs text-muted-foreground">
                            {setting.sync_proposed_dates
                                ? 'オフにすると候補日程の同期を停止します(専用カレンダーは削除されません)。'
                                : 'オンにすると候補日程の同期を再開します(専用カレンダーがない場合は再作成されます)。'}
                        </p>
                    </div>
                )}
            </div>
        </li>
    );
};

const CalendarSettingsSection = () => {
    const { calendarSettings, isLoading, error, refetch } = useCalendarSettings();
    const { error: accountError } = useAccounts();
    const { update } = useUpdateCalendarSetting();

    // 再認可が必要な間は設定変更を無効化し、Google 連携セクションの再認可導線を優先する(screen-design 5.8)
    const reauthRequired = isGoogleReauthorizationRequiredError(accountError);

    const settings = calendarSettings ? sortSettings(calendarSettings) : [];
    const primaryId = settings.find((setting) => setting.role === 'primary')?.id ?? '';

    return (
        <Card>
            <CardHeader>
                <CardTitle>カレンダー設定</CardTitle>
                <CardDescription>
                    アプリに表示するカレンダーと、確定した予定の登録先(メイン)を設定します。登録先の変更は今後作成するイベントに適用されます(既存イベントの登録先は変わりません)。
                </CardDescription>
            </CardHeader>
            <CardContent>
                {reauthRequired && (
                    <p className="mb-4 rounded-md bg-yellow-500/10 p-3 text-sm text-foreground">
                        Google アカウントの再認可が必要なため、カレンダー設定は変更できません。上の「Google
                        連携」から再認可してください。
                    </p>
                )}
                {isLoading ? (
                    <div className="space-y-4">
                        <Skeleton className="h-14 w-full" />
                        <Skeleton className="h-14 w-full" />
                        <Skeleton className="h-14 w-full" />
                    </div>
                ) : error ? (
                    <div className="space-y-3 py-4 text-center">
                        <p className="text-sm text-muted-foreground">
                            カレンダー設定を取得できませんでした。
                        </p>
                        <Button variant="outline" onClick={() => void refetch()}>
                            再試行
                        </Button>
                    </div>
                ) : settings.length === 0 ? (
                    <div className="space-y-3 py-4 text-center">
                        <p className="text-sm text-muted-foreground">
                            カレンダーがまだ同期されていません。ダッシュボードを開くと Google
                            カレンダーから自動で同期されます。
                        </p>
                        <Button asChild variant="outline">
                            <Link href="/dashboard">ダッシュボードを開く</Link>
                        </Button>
                    </div>
                ) : (
                    <RadioGroup
                        value={primaryId}
                        disabled={reauthRequired}
                        onValueChange={(id) => update({ id, payload: { role: 'primary' } })}
                        className="block"
                    >
                        <ul className="divide-y divide-border">
                            {settings.map((setting) => (
                                <CalendarSettingRow
                                    key={setting.id}
                                    setting={setting}
                                    disabled={reauthRequired}
                                    onToggleVisible={(checked) =>
                                        update({ id: setting.id, payload: { is_visible: checked } })
                                    }
                                    onToggleSync={(checked) =>
                                        update({
                                            id: setting.id,
                                            payload: { sync_proposed_dates: checked },
                                        })
                                    }
                                />
                            ))}
                        </ul>
                    </RadioGroup>
                )}
            </CardContent>
        </Card>
    );
};

export default CalendarSettingsSection;
