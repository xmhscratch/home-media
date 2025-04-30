import {
  Component,
  Input,
  Signal,
  WritableSignal,
  ElementRef,
  ViewChild,
  ViewEncapsulation,
  ChangeDetectorRef,
} from '@angular/core';
import { signal } from '@angular/core';
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
import { Subscription, Subject, debounceTime, delay, of, first } from 'rxjs';

import { IFileListItem, ISocketMessage } from '../../../../../types/storage';
import { EnvService } from '@/env.service';
import { StorageService } from '@/storage.service';
import { FileService, TFileList } from '@/file.service';
import { FileSizePipe } from '@/filesize.pipe';
import { MyPercentPipe } from '@/mypercent.pipe';
// import parseRange from 'range-parser';

const BUFFER_CHUNK_SIZE = 1024 * 100;

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
    <string>this.envService.get('endpoint.file'),
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
            console.log(this.f());
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

  async initAudio() {
    let context;
    if (!context) {
      // @ts-ignore
      context = new (window.AudioContext || window.webkitAudioContext)();
    }

    // @for (dub of f().dubs; track $index) {
    //   <source
    //     src="http://192.168.56.55:4150/{{ f().nodeId }}/{{
    //       f().fileKey
    //     }}.{{ dub.lang_code }}{{ dub.stream_index }}.mp4"
    //     type="audio/ogg"
    //   />
    // }

    // Fetch and decode the audio file
    const dubs = this.f().dubs;
    if (!dubs || dubs.length == 0) {
      return;
    }

    // interval(10000).pipe(
    //   takeWhile(() => currentStart < totalSize),
    //   switchMap(() => { }),
    // ).subscribe()

    const fileUrl = `${this.baseURL()}${this.f().nodeId}/${this.f().fileKey}.${dubs[0].lang_code}${dubs[0].stream_index}.mp4`;
    // this.fileService.fetchFile(fileUrl, 0, 1024 - 1).subscribe()
    const response = await fetch(fileUrl);
    const arrayBuffer = await response.arrayBuffer();
    console.log(arrayBuffer);
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

    const audio = await this.initAudio();
    const audioEventInit = () => {
      if (!audio) {
        return;
      }

      // @ts-ignore
      this.player.on('playing', async () => {
        await audio.start(this.getPlayerCurrentTime());
      });

      // @ts-ignore
      this.player.on('pause', async () => {
        await audio.pause();
      });

      // @ts-ignore
      this.player.on('volumechange', async () => {
        audio.setVolume(this.getPlayerCurrentVolume());
      });

      // @ts-ignore
      this.player.on('seeked', async () => {
        await audio.start(this.getPlayerCurrentTime());
      });

      // @ts-ignore
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
