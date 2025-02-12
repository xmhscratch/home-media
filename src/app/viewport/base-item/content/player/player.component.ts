import { Component, Input } from '@angular/core';
import { OnInit, OnDestroy } from '@angular/core';

import { ElementRef, ViewChild, ViewEncapsulation } from '@angular/core';
// import { default as videojs } from 'video.js';
// import { default as VideoPlayer } from 'video.js/dist/types/player.d';

type VJS_PRELOAD = "auto" | "metadata" | "none"

@Component({
  selector: 'app-player',
  standalone: true,
  imports: [],
  templateUrl: './player.component.html',
  styleUrl: './player.component.scss'
})
export class CPlayer implements OnInit, OnDestroy {
  @ViewChild('target', { static: true }) target!: ElementRef;

  @Input() options!: {
    fluid?: boolean,
    aspectRatio?: string,
    autoplay?: boolean,
    controls?: boolean,
    preload?: VJS_PRELOAD,
    width?: number,
    height?: number,
    autoSetup?: boolean,
    poster?: string,
    sources: {
      src: string,
      type: string,
    }[],
  };

  // player!: VideoPlayer;

  constructor(
    private elementRef: ElementRef,
  ) { }

  ngOnInit(): void {
    // @ts-ignore
    const player = videojs(this.target.nativeElement, {
      autoplay: false,
      controls: true,
      fluid: true,
    });
  }

  ngOnDestroy(): void {
    // if (this.player) {
    //   this.player.dispose();
    // }
  }
}
