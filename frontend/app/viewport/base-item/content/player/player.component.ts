import {
  Component,
  Input,
  Signal,
  // Inject,
  WritableSignal,
  ElementRef,
  ViewChild,
  ViewEncapsulation,
  ChangeDetectorRef,
} from '@angular/core';
import { signal } from '@angular/core';
// import { DOCUMENT } from '@angular/common';
import { NgIf, PercentPipe } from '@angular/common';
import { OnInit, OnDestroy, AfterViewInit } from '@angular/core';

// import { MessageService } from 'primeng/api';
// import { DialogService, DynamicDialogRef } from 'primeng/dynamicdialog';
import { ButtonModule } from 'primeng/button';
import { ProgressBar } from 'primeng/progressbar';
import { Skeleton } from 'primeng/skeleton';
import { ScrollPanelModule } from 'primeng/scrollpanel';
import { CardModule } from 'primeng/card';
import { Message } from 'primeng/message';
import { PanelModule } from 'primeng/panel';

import { clamp as ldClamp } from 'lodash-es';
import { Subscription, Subject, debounceTime, timeout, lastValueFrom, delay, of, first } from 'rxjs';

import { IFileListItem, ISocketMessage, IDubFileInfo } from '@T/storage';
import { EnvService } from '@/env.service';
import { StorageService } from '@/storage.service';
import { FileService, TFileList } from '@/file.service';
import { FileURLPipe } from '@/fileurl.pipe';
import { FileSizePipe } from '@/filesize.pipe';
import { MyPercentPipe } from '@/mypercent.pipe';
// import parseRange from 'range-parser';

@Component({
  selector: 'app-player',
  standalone: true,
  imports: [
    NgIf,
    ButtonModule,
    Skeleton,
    ScrollPanelModule,
    CardModule,
    ProgressBar,
    Message,
    PanelModule,
    FileURLPipe,
    FileSizePipe,
    PercentPipe,
    MyPercentPipe,
  ],
  templateUrl: './player.component.html',
  styleUrl: './player.component.scss',
  encapsulation: ViewEncapsulation.None,
})
export class CPlayer implements OnInit, OnDestroy, AfterViewInit {
  @ViewChild('videoEl', {
    /*static: true*/
  })
  videoEl!: ElementRef;

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
  };

  player!: any;

  @Input() file!: IFileListItem;

  baseURL: Signal<string> = signal<string>(
    <string>this.fileService.baseURL(),
  );

  f: WritableSignal<IFileListItem> = signal<IFileListItem>(this.file);
  loaded: WritableSignal<boolean> = signal<boolean>(false);

  destroy$: Subscription = new Subscription();

  constructor(
    private envService: EnvService,
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
            of('')
              .pipe(delay(1), first())
              .subscribe(() => this.initPlayer());
          }
        }),
      );
      this.fileService.getNativeSocket(f).next({ event: 9 });
    })(this.file);

    this.destroy$.add(
      this.fileService.changed$.subscribe((f: IFileListItem) => {
        this.f.update((v) => v);
      }),
    );
  }

  ngOnDestroy(): void {
    if (this.player) {
      this.player.dispose();
    }
  }

  ngAfterViewInit(): void { }

  handleItemSelect(e: Event, f: IFileListItem) {
    if (!f) {
      return;
    }

    this.destroy$.add(
      this.storageService
        .fetchSource(f.sessionId, f.path)
        // .pipe(tap((v) => { console.log(v) }))
        .subscribe(),
    );
  }

  async initAudio(dub: IDubFileInfo) {
    let context;
    if (!context) {
      // @ts-ignore
      context = new (window.AudioContext || window.webkitAudioContext)();
    }

    // console.log(dub)
    // 0-22920625/22920626
    // number_of_bytes:22721005
    // number_of_frames:61157

    const res$ = this.fileService
      .fetchFile(`${this.f().nodeId}/${this.f().fileKey}.${dub.lang_code}${dub.stream_index}.mp4`, 0)
      .pipe(timeout({ first: 10_000, each: 10_000 }));

    const arrayBuffer = await lastValueFrom(res$);
    // console.log(arrayBuffer);

    const audioBuffer = await context.decodeAudioData(arrayBuffer);

    // Create a GainNode for volume control
    const gain: GainNode = context.createGain();
    gain.gain.value = 0.76; // Default volume (1 = 100%)
    gain.connect(context.destination);

    // Create a buffer source
    const initSource = () => {
      const source: AudioBufferSourceNode = context.createBufferSource();
      source.buffer = audioBuffer;
      source.connect(gain);

      return source;
    };

    let source = initSource();
    let isPlaying = false;

    return {
      // Function to change volume
      async start(time: number) {
        if (!isPlaying) {
          source.start();
          isPlaying = true;
        } else {
          if (context.state === 'suspended') {
            await context.resume();
          }

          // Stop current playback
          source.stop();
          source.disconnect();

          // Restart from seeked position
          source = initSource();
          source.start(0, time || 0);

          // interval(10000).pipe(
          //   takeWhile(() => currentStart < totalSize),
          //   switchMap(() => { }),
          // ).subscribe()
        }
      },
      async pause() {
        if (context && context.state === 'running') {
          await context.suspend();
        }
      },
      async close() {
        context.close();
      },
      setVolume(value: number) {
        gain.gain.value = ldClamp(value, 0.0, 1.0);
      },
    };
  }

  async initPlayer() {
    if (this.player) {
      return;
    }

    // Fetch and decode the audio file
    const dubs = this.f().dubs;
    if (!dubs || dubs.length == 0) {
      return;
    }

    const audio = await this.initAudio(dubs[0]);
    const audioEventInit = () => {
      if (!audio) {
        return;
      }

      this.player.on('playing', async () => {
        await audio.start(this.getPlayerCurrentTime());
      });

      this.player.on('pause', async () => {
        await audio.pause();
      });

      this.player.on('volumechange', async () => {
        audio.setVolume(this.getPlayerCurrentVolume());
      });

      this.player.on('seeked', async () => {
        await audio.start(this.getPlayerCurrentTime());
      });

      this.player.on('dispose', async () => {
        await audio.close();
      });
    };

    // @ts-ignore
    this.player = videojs(
      this.videoEl.nativeElement,
      {
        autoplay: false,
        controls: true,
        fluid: true,
        responsive: true,
        userActions: {
          hotkeys: (e: KeyboardEvent) => {
            switch (e.key) {
              case 'ArrowLeft': {
                const destTime = this.getPlayerCurrentTime()-5;
                this.player.currentTime(destTime) && audio.start(destTime);
                break;
              }
              case 'ArrowRight': {
                const destTime = this.getPlayerCurrentTime()+5;
                this.player.currentTime(destTime) && audio.start(destTime);
                break;
              }
              case 'ArrowUp': {
                const destVol = this.getPlayerCurrentVolume()+0.1;
                this.player.volume(destVol) && audio.setVolume(destVol);
                break;
              }
              case 'ArrowDown': {
                const destVol = this.getPlayerCurrentVolume()-0.1;
                this.player.volume(destVol) && audio.setVolume(destVol);
                break;
              }
              case ' ': case 'Spacebar': {
                if (!this.player.paused()) {
                  this.player.pause();
                  audio.pause();
                } else {
                  this.player.play();
                  audio.start(this.getPlayerCurrentTime());
                }
                break;
              }
              default: break;
            }
          }
        }
      },
      audioEventInit,
    );
  }

  getPlayerCurrentVolume() {
    let val;
    if (this.player.ended()) {
      val = this.player.volume();
    } else {
      val = this.player.scrubbing()
        ? this.player.getCache().volume
        : this.player.volume();
    }
    return val;
  }

  getPlayerCurrentTime() {
    let val;
    if (this.player.ended()) {
      val = this.player.duration();
    } else {
      val = this.player.scrubbing()
        ? this.player.getCache().currentTime
        : this.player.currentTime();
    }
    return val;
  }
}
