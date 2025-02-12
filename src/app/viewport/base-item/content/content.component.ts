import { Component } from '@angular/core';
import { inject } from '@angular/core';
import { OnInit } from '@angular/core';
import { NgFor, NgClass } from '@angular/common';

import { ActivatedRoute } from '@angular/router';

import { MessageService } from 'primeng/api';
import { DialogService, DynamicDialogRef } from 'primeng/dynamicdialog';
import { ButtonModule } from 'primeng/button';
import { ProgressBar } from 'primeng/progressbar';
import { Skeleton } from 'primeng/skeleton';
import { ScrollPanelModule } from 'primeng/scrollpanel';

import { StorageService } from '@/storage.service'
import { CPlayer } from './player/player.component';

@Component({
  selector: 'app-content',
  standalone: true,
  imports: [
    NgFor, NgClass,
    ButtonModule, ProgressBar, Skeleton, ScrollPanelModule,
    CPlayer,
  ],
  templateUrl: './content.component.html',
  styleUrl: './content.component.scss',
  providers: [DialogService, MessageService],
})
export class CContent implements OnInit {

  private readonly route = inject(ActivatedRoute);
  selectedId!: string;

  data = [
    {
      id: '1000',
      code: 'f230fh0g3',
      name: 'Bamboo Watch',
      description: 'Product Description',
      image: 'bamboo-watch.jpg',
      price: 65,
      category: 'Accessories',
      quantity: 24,
      inventoryStatus: 'INSTOCK',
      rating: 5
    },
    {
      id: '1001',
      code: 'nvklal433',
      name: 'Black Watch',
      description: 'Product Description',
      image: 'black-watch.jpg',
      price: 72,
      category: 'Accessories',
      quantity: 61,
      inventoryStatus: 'INSTOCK',
      rating: 4
    },
    {
      id: '1002',
      code: 'zz21cz3c1',
      name: 'Blue Band',
      description: 'Product Description',
      image: 'blue-band.jpg',
      price: 79,
      category: 'Fitness',
      quantity: 2,
      inventoryStatus: 'LOWSTOCK',
      rating: 3
    },
    {
      id: '1003',
      code: '244wgerg2',
      name: 'Blue T-Shirt',
      description: 'Product Description',
      image: 'blue-t-shirt.jpg',
      price: 29,
      category: 'Clothing',
      quantity: 25,
      inventoryStatus: 'INSTOCK',
      rating: 5
    },
    {
      id: '1004',
      code: 'h456wer53',
      name: 'Bracelet',
      description: 'Product Description',
      image: 'bracelet.jpg',
      price: 15,
      category: 'Accessories',
      quantity: 73,
      inventoryStatus: 'INSTOCK',
      rating: 4
    }
  ];

  constructor(
    private storage: StorageService,
    private dialogService: DialogService,
    private ref: DynamicDialogRef,
  ) { }

  ngOnInit() { }

  closeDialog(data: any) {
    this.ref.close(data);
  }
}
