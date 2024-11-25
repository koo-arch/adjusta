import { format } from 'date-fns';
import { ja } from 'date-fns/locale';

export const formatJaDate = (date: Date) => {
    return format(date, 'M月d日(E) H:mm', { locale: ja });
}

export const formatDate = (date: Date) => {
    return format(date, 'yyyy-MM-dd');
}

// 日付までが同じ場合、終了時は時刻のみを表示する
export const formatJaDateSpan = (start: Date, end: Date) => {
    if (start.toDateString === end.toDateString) {
        return `${formatJaDate(start)} 〜 ${format(end, 'H:mm')}`;
    }
    return `${formatJaDate(start)} 〜 ${formatJaDate(end)}`;
}