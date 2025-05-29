import { Injectable } from '@angular/core';
import { Pipe, PipeTransform } from '@angular/core';
// import { parseInt as ldParseInt } from 'lodash-es';
import { FileService } from '@/file.service';

@Pipe({
  name: 'fileurl',
  standalone: true,
})
@Injectable({
  providedIn: 'root',
})
export class FileURLPipe implements PipeTransform {
  constructor(private fileService: FileService) {}

  transform(fileURL: string, fileExt: string): URL {
    fileExt = fileExt || 'bin';
    return new URL(`${fileURL}.${fileExt}`, <string>this.fileService.baseURL());
  }
}
