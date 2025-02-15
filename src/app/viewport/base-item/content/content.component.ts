import { Component, WritableSignal } from '@angular/core';
import { inject, signal } from '@angular/core';
import { OnInit } from '@angular/core';
import { NgFor, NgClass } from '@angular/common';

import { ActivatedRoute } from '@angular/router';

import { MessageService } from 'primeng/api';
import { DialogService, DynamicDialogRef } from 'primeng/dynamicdialog';
import { ButtonModule } from 'primeng/button';
import { ProgressBar } from 'primeng/progressbar';
import { Skeleton } from 'primeng/skeleton';
import { ScrollPanelModule } from 'primeng/scrollpanel';

import { switchMap, tap } from 'rxjs/operators';

import { FileSizePipe } from '@/filesize.pipe';
import { IFileListItem } from '@/storage.d';
import { StorageService } from '@/storage.service';
// import { CPlayer } from './player/player.component';

@Component({
  selector: 'app-content',
  standalone: true,
  imports: [
    NgFor, NgClass,
    ButtonModule, Skeleton, ScrollPanelModule,
    // CPlayer, ProgressBar,
    FileSizePipe,
  ],
  templateUrl: './content.component.html',
  styleUrl: './content.component.scss',
  providers: [DialogService, MessageService],
})
export class CContent implements OnInit {

  private readonly route = inject(ActivatedRoute);

  files: WritableSignal<IFileListItem[]> = signal<IFileListItem[]>([]);
  // files$ = toObservable(this.files);
  selected: WritableSignal<IFileListItem> = signal<IFileListItem>(this.files()[0]);

  // data = [
  //   // {
  //   //   id: '1000',
  //   //   name: 'Bamboo Watch',
  //   //   description: 'Product Description',
  //   // },
  // ];

  constructor(
    private storage: StorageService,
    private dialogService: DialogService,
    private ref: DynamicDialogRef,
  ) { }

  ngOnInit() {
    this.storage
      .getData()
      .pipe(
        switchMap(({ root, active, paths, nodes }) => {
          return this.storage.createSession(active)
        }),
        // tap((v) => { console.log(v) }),
      )
      .subscribe((files: IFileListItem[]) => {
        this.selected.set(files[0])
        this.files.set(files)
      });
  }

  handleItemSelect(e: MouseEvent, file: IFileListItem) {
    this.storage.fetchSource(file.sessionId, file.path)
      .pipe(
        tap((v) => { console.log(v) }),
      )
      .subscribe()

    this.selected.set(file)
  }

  closeDialog(data: any) {
    this.ref.close(data);
  }
}
