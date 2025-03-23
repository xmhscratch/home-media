import { Component, WritableSignal, signal, Signal, OnInit, OnDestroy, computed } from '@angular/core';
import { MenuItem } from 'primeng/api';
import { IINode } from '../../../types/storage';
// import { StorageService } from '@/storage.service';

@Component({
  selector: 'app-item',
  standalone: true,
  imports: [],
  templateUrl: './item.component.html',
  styleUrl: './item.component.scss'
})
export class CItem implements IINode, MenuItem {
  id: string = ''
  root: string = ''

  title?: string = ''
  depth?: number = 0
}
