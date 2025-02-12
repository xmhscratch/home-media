import { Directive } from '@angular/core';

import { CItem } from './item.component';

@Directive({
  selector: '[appItem]',
  standalone: true,
})
export class ItemDirective {

  constructor() {
    // console.log(234234)
  }

}
