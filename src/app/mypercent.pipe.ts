import { Pipe, PipeTransform } from '@angular/core';
import { parseInt as ldParseInt } from 'lodash-es';

@Pipe({
    name: 'mypercent',
    standalone: true,
})
export class MyPercentPipe implements PipeTransform {
    transform(value: string, min: number): string {
        var valueNumber: number = ldParseInt(value) || 0
        return (valueNumber >= min) ? String(valueNumber / 100) : '';
    }
}
