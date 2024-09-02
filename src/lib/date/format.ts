import { format } from 'date-fns';
import { ja } from 'date-fns/locale';

export const formatJaDate = (date: Date) => {
    return format(date, 'M月d日(E) HH:mm', { locale: ja });
}

export const formatDate = (date: Date) => {
    return format(date, 'yyyy-MM-dd');
}