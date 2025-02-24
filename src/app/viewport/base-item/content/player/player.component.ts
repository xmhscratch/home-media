import { Component, Input, WritableSignal, ElementRef, ViewChild, ViewEncapsulation, ChangeDetectorRef } from '@angular/core';
import { signal } from '@angular/core';
import { NgIf } from '@angular/common';
import { OnInit, OnDestroy, AfterViewInit } from '@angular/core';

import { MessageService } from 'primeng/api';
import { DialogService, DynamicDialogRef } from 'primeng/dynamicdialog';
import { ButtonModule } from 'primeng/button';
import { ProgressBar } from 'primeng/progressbar';
import { Skeleton } from 'primeng/skeleton';
import { ScrollPanelModule } from 'primeng/scrollpanel';
import { CardModule } from 'primeng/card';
import { Message } from 'primeng/message';
import { PanelModule } from 'primeng/panel';

import { Subscription, Subject, debounceTime, delay, of, takeUntil } from 'rxjs';

import { IFileListItem, ISocketMessage } from '@/storage.d';
import { StorageService } from '@/storage.service';
import { FileService, TFileList } from '@/file.service';
import { FileSizePipe } from '@/filesize.pipe';
// import { default as videojs } from 'video.js';
// import { default as VideoPlayer } from 'video.js/dist/types/player.d';

type VJS_PRELOAD = "auto" | "metadata" | "none"

@Component({
  selector: 'app-player',
  standalone: true,
  imports: [
    NgIf,
    ButtonModule, Skeleton, ScrollPanelModule, CardModule, ProgressBar, Message, PanelModule,
    FileSizePipe,
  ],
  templateUrl: './player.component.html',
  styleUrl: './player.component.scss',
  encapsulation: ViewEncapsulation.None,
})
export class CPlayer implements OnInit, OnDestroy, AfterViewInit {

  @ViewChild('videoEl', { /*static: false*/ }) videoEl!: ElementRef;
  @ViewChild('audioEl', { /*static: false*/ }) audioEl!: ElementRef;

  // @Input() options!: {
  //   fluid?: boolean,
  //   aspectRatio?: string,
  //   autoplay?: boolean,
  //   controls?: boolean,
  //   preload?: VJS_PRELOAD,
  //   width?: number,
  //   height?: number,
  //   autoSetup?: boolean,
  //   poster?: string,
  //   sources: {
  //     src: string,
  //     type: string,
  //   }[],
  // };

  cardStyles = {
    shadow: 'none',
    bodyPadding: '0',
  }

  player!: any;

  @Input() file!: IFileListItem

  f: WritableSignal<IFileListItem> = signal<IFileListItem>(this.file);
  loaded: WritableSignal<boolean> = signal<boolean>(false);

  destroy$: Subscription = new Subscription();

  constructor(
    private storageService: StorageService,
    private fileService: FileService,
    private cdRef: ChangeDetectorRef,
  ) { }

  ngOnInit(): void {
    ((f: IFileListItem) => {
      this.destroy$.add(
        this.fileService.activateFile(f).subscribe((m: IFileListItem) => {
          this.f.set(m);
          this.loaded.set(true);

          if (m.sourceReady) {
            of('').pipe(delay(100)).subscribe(() => this.initPlayer())
          }
        })
      )
      this.fileService.getNativeSocket(f).next({ event: 9 });
    })(this.file);

    this.destroy$.add(
      this.fileService.changed$
        // .pipe(debounceTime(1000))
        .subscribe((f: IFileListItem) => {
          this.f.update((v) => (v));
          // this.cdRef.detectChanges();
        }),
    )

    // of('').pipe(
    //   // delay(350),
    //   takeUntil(this.loaded$),
    // ).subscribe(() => {
    //   console.log(this.f())
    //   // this.initPlayer();
    // });
  }

  ngOnDestroy(): void {
    this.audioEl.nativeElement.src = "";
    if (this.player) {
      this.player.dispose();
    }
  }

  ngAfterViewInit(): void { }

  handleItemSelect(e: Event, f: IFileListItem) {
    if (!f) { return }

    this.destroy$.add(
      this.storageService
        .fetchSource(f.sessionId, f.path)
        // .pipe(tap((v) => { console.log(v) }))
        .subscribe()
    )
  }

  initPlayer() {
    if (this.player) { return }

    // @ts-ignore
    this.player = videojs(this.videoEl.nativeElement, {
      autoplay: false,
      controls: true,
      fluid: true,
      responsive: true,
    });

    // @ts-ignore
    this.player.on('playing', () => {
      this.audioEl.nativeElement.play()
    });

    // @ts-ignore
    this.player.on('pause', () => {
      this.audioEl.nativeElement.pause()
    });

    // @ts-ignore
    this.player.on('volumechange', () => {
      let val;
      if (this.player.ended()) {
        val = this.player.volume();
      } else {
        val = (this.player.scrubbing()) ? this.player.getCache().volume : this.player.volume();
      }
      this.audioEl.nativeElement.volume = val;
    });

    // @ts-ignore
    this.player.on('seeked', () => {
      let val;
      if (this.player.ended()) {
        val = this.player.duration();
      } else {
        val = (this.player.scrubbing()) ? this.player.getCache().currentTime : this.player.currentTime();
      }
      this.audioEl.nativeElement.currentTime = val;
    });
  }
}
