import { Component, WritableSignal, ChangeDetectorRef } from '@angular/core';
import { inject, signal } from '@angular/core';
import { OnInit, OnDestroy } from '@angular/core';
import { NgFor, NgClass } from '@angular/common';
import { KeyValuePipe } from '@angular/common';
import { toObservable } from '@angular/core/rxjs-interop';

import { MessageService } from 'primeng/api';
import { DialogService, DynamicDialogRef } from 'primeng/dynamicdialog';
import { ButtonModule } from 'primeng/button';
import { ProgressBar } from 'primeng/progressbar';
import { Skeleton } from 'primeng/skeleton';
import { ScrollPanelModule } from 'primeng/scrollpanel';
import { CardModule } from 'primeng/card';
import { Message } from 'primeng/message';
import { PanelModule } from 'primeng/panel';

import { Subscription } from 'rxjs';
// import { WebSocketSubject } from 'rxjs/webSocket';
// import { switchMap, tap } from 'rxjs/operators';
// import { map as ldMap } from 'lodash-es';

import { FileSizePipe } from '@/filesize.pipe';
import { IFileListItem } from '@/storage.d';
import { StorageService } from '@/storage.service';
import { FileService, TFileList } from '@/file.service';
import { CPlayer } from './player/player.component';

@Component({
  selector: 'app-content',
  standalone: true,
  imports: [
    NgFor, NgClass,
    ButtonModule, ScrollPanelModule, CardModule, PanelModule,
    CPlayer,
    KeyValuePipe,
  ],
  templateUrl: './content.component.html',
  styleUrl: './content.component.scss',
  providers: [DialogService, MessageService],
})
export class CContent implements OnInit, OnDestroy {

  files: WritableSignal<TFileList> = signal<TFileList>({});
  files$ = toObservable(this.files);

  destroy$: Subscription = new Subscription();

  constructor(
    private storageService: StorageService,
    private fileService: FileService,
    // private dialogService: DialogService,
    private ref: DynamicDialogRef,
  ) { }

  ngOnInit() {
    this.destroy$.add(
      this.fileService.loadFiles().subscribe((files: TFileList) => {
        this.files.set(files);
      }),
    )

    this.destroy$.add(() => Object.keys(this.files())
      .forEach(this.fileService.disconnect.bind(this.fileService)))
  }

  ngOnDestroy(): void {
    if (this.destroy$) { this.destroy$.unsubscribe() }
  }

  closeDialog(data: any) {
    this.ref.close(data);
  }
}
