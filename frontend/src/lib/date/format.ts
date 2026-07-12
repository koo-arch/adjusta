import { format, isSameDay } from 'date-fns';
import { ja } from 'date-fns/locale';

export const formatJaDate = (date: Date) => {
    return format(date, 'M月d日(E) H:mm', { locale: ja });
}

export const formatDate = (date: Date) => {
    return format(date, 'yyyy-MM-dd');
}

// 日付までが同じ場合、終了時は時刻のみを表示する
export const formatJaDateSpan = (start: Date, end: Date) => {

    if (isSameDay(start, end)) {
        return `${formatJaDate(start)} - ${format(end, 'H:mm')}`;
    }
    return `${formatJaDate(start)} - ${formatJaDate(end)}`;
}

export const formatJaDateOnly = (date: Date) => {
    return format(date, 'M月d日(E)', { locale: ja });
}

// 時刻のみの範囲表示。日跨ぎは終了側に日付を含める
export const formatJaTimeSpan = (start: Date, end: Date) => {
    if (isSameDay(start, end)) {
        return `${format(start, 'H:mm')} - ${format(end, 'H:mm')}`;
    }
    return `${format(start, 'H:mm')} - ${formatJaDate(end)}`;
}