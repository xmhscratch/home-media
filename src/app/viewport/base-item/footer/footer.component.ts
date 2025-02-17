import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';

import { DynamicDialogRef } from 'primeng/dynamicdialog';
import { ButtonModule } from 'primeng/button';

@Component({
  selector: 'app-footer',
  standalone: true,
  imports: [
    ButtonModule, FormsModule,
  ],
  templateUrl: './footer.component.html',
  styleUrl: './footer.component.scss',
})
export class CFooter implements OnInit {

  constructor(
    private ref: DynamicDialogRef,
  ) { }

  ngOnInit(): void { }

  // closeDialog(data: any) {
  //   this.ref.close(data);
  // }
}
