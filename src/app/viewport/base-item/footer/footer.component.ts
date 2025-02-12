import { Component, WritableSignal, OnInit } from '@angular/core';
import { signal } from '@angular/core';
import { FormsModule } from '@angular/forms';

import { DynamicDialogRef } from 'primeng/dynamicdialog';
import { ButtonModule } from 'primeng/button';
import { Select, SelectChangeEvent } from 'primeng/select';

import { switchMap, tap } from 'rxjs/operators';

import { IFileListItem } from '@/storage.d'
import { StorageService } from '@/storage.service'

@Component({
  selector: 'app-footer',
  standalone: true,
  imports: [
    ButtonModule, FormsModule, Select,
  ],
  templateUrl: './footer.component.html',
  styleUrl: './footer.component.scss',
})
export class CFooter implements OnInit {

  files: WritableSignal<IFileListItem[]> = signal<IFileListItem[]>([]);
  // files$ = toObservable(this.files);
  selected!: IFileListItem;

  constructor(
    private storage: StorageService,
    private ref: DynamicDialogRef,
  ) { }

  ngOnInit(): void {
    this.storage
      .getData()
      .pipe(
        switchMap(({ root, active, paths, nodes }) => {
          return this.storage.createSession(active)
        }),
        tap((v) => { console.log(v) }),
      )
      .subscribe((files: IFileListItem[]) => {
        this.selected = files[0]
        this.files.set(files)
      });
  }

  handleItemSelect(e: SelectChangeEvent) {
    console.log(e.value)
  }

  closeDialog(data: any) {
    this.ref.close(data);
  }
}
