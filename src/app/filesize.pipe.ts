import { Pipe, PipeTransform } from '@angular/core';
import { filesize } from "filesize";

@Pipe({
    name: 'filesize',
    standalone: true,
})
export class FileSizePipe implements PipeTransform {
    transform(value: number): string {
        return filesize(value, { standard: "jedec" });
    }
}
