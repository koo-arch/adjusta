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
import { useCandidateSyncSetting } from '@/features/auth/hooks/useCandidateSyncSetting';
import { useUpdateCalendarSetting } from '@/features/auth/hooks/useUpdateCalendarSetting';
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

// 現在の API は accessRole を返さないため、Google 公式の祝日カレンダー ID を判定する。
const isGoogleHolidayCalendar = (setting: CalendarSetting) =>
    setting.google_calendar_id.includes('#holiday@group.v.calendar.google.com');

const canBePrimary = (setting: CalendarSetting) =>
    setting.role !== 'adjusta_candidate' && !isGoogleHolidayCalendar(setting);

interface CalendarSettingRowProps {
    setting: CalendarSetting;
    disabled: boolean;
    onToggleVisible: (checked: boolean) => void;
}

const CalendarSettingRow: React.FC<CalendarSettingRowProps> = ({
    setting,
    disabled,
    onToggleVisible,
}) => {
    const isCandidate = setting.role === 'adjusta_candidate';
    const badge = roleBadges[setting.role];

    return (
        <li className="flex flex-col gap-3 py-4 sm:flex-row sm:items-start sm:justify-between">
            <div className="flex min-w-0 items-start gap-3">
                <div className="min-w-0">
                    <div className="flex flex-wrap items-center gap-2">
                        <p className="break-words font-medium">{setting.summary}</p>
                        <Badge className={badge.className}>{badge.label}</Badge>
                    </div>
                    {setting.timezone && (
                        <p className="mt-0.5 text-sm text-muted-foreground">{setting.timezone}</p>
                    )}
                    {isCandidate && (
                        <p className="mt-1 text-sm text-muted-foreground">
                            候補日程の仮予定用に自動管理します。
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
            </div>
        </li>
    );
};

const CalendarSettingsSection = () => {
    const { calendarSettings, isLoading, error, refetch } = useCalendarSettings();
    const { connectionState } = useAccounts();
    const { update } = useUpdateCalendarSetting();
    const candidateSync = useCandidateSyncSetting();

    // 再認可が必要な間は設定変更を無効化し、Google 連携セクションの再認可導線を優先する(screen-design 5.8)
    const reauthRequired = connectionState.kind === 'reauthorization_required';

    const settings = calendarSettings ? sortSettings(calendarSettings) : [];
    const primaryId = settings.find((setting) => setting.role === 'primary')?.id ?? '';
    const primaryChoices = settings.filter(canBePrimary);

    return (
        <Card>
            <CardHeader>
                <CardTitle>カレンダー設定</CardTitle>
                <CardDescription>
                    候補日程の同期、確定予定の登録先、表示するカレンダーを設定します。
                </CardDescription>
            </CardHeader>
            <CardContent>
                {reauthRequired && (
                    <p className="mb-4 rounded-md bg-yellow-500/10 p-3 text-sm text-foreground">
                        Google アカウントの再認可が必要なため、カレンダー設定は変更できません。上の「Google
                        連携」から再認可してください。
                    </p>
                )}
                <section className="border-b border-border py-4 first:pt-0">
                    <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                    <div>
                        <h3 className="font-medium">候補日程の同期</h3>
                        <p className="mt-1 text-sm text-muted-foreground">
                            専用カレンダーに仮予定を表示します。
                        </p>
                    </div>
                    {candidateSync.isLoading ? (
                        <Skeleton className="h-6 w-11 shrink-0" />
                    ) : (
                        <Switch
                            aria-label="候補日程の Google Calendar 同期"
                            checked={candidateSync.setting?.enabled ?? false}
                            disabled={reauthRequired || candidateSync.isUpdating}
                            onCheckedChange={candidateSync.setEnabled}
                        />
                    )}
                    {candidateSync.error && (
                        <p className="text-sm text-destructive">同期設定を取得できませんでした。</p>
                    )}
                    </div>
                </section>
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
                    <>
                    <section className="border-b border-border py-4">
                        <h3 className="font-medium">確定予定の登録先</h3>
                        <p className="mt-1 text-sm text-muted-foreground">今後確定する予定を追加するカレンダーです。</p>
                    <RadioGroup
                        value={primaryId}
                        disabled={reauthRequired}
                        onValueChange={(id) => update({ id, payload: { role: 'primary' } })}
                        className="mt-3 space-y-2"
                    >
                        {primaryChoices.map((setting) => (
                            <label key={setting.id} htmlFor={`primary-${setting.id}`} className="flex cursor-pointer items-center gap-3 rounded-md border border-border px-3 py-2 text-sm">
                                <RadioGroupItem value={setting.id} id={`primary-${setting.id}`} />
                                <span>{setting.summary}</span>
                            </label>
                        ))}
                    </RadioGroup>
                    </section>
                    <section className="pt-4">
                        <h3 className="font-medium">表示するカレンダー</h3>
                        <p className="mt-1 text-sm text-muted-foreground">予定確認に表示するカレンダーを選びます。</p>
                        <ul className="mt-2 divide-y divide-border">
                            {settings.map((setting) => (
                                <CalendarSettingRow
                                    key={setting.id}
                                    setting={setting}
                                    disabled={reauthRequired}
                                    onToggleVisible={(checked) =>
                                        update({ id: setting.id, payload: { is_visible: checked } })
                                    }
                                />
                            ))}
                        </ul>
                    </section>
                    </>
                )}
            </CardContent>
        </Card>
    );
};

export default CalendarSettingsSection;
