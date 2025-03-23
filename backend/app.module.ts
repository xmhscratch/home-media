import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';

import myconf from './app.config';

import { AppController } from './app.controller';
import { AppService } from './app.service';

import { StorageController } from './storage/storage.controller';
import { NyaaController } from './custom/nyaa.controller';

import { StorageService } from './storage/storage.service';
import { TreeService } from './storage/tree.service';
import { NyaaService } from './custom/nyaa.service';

@Module({
  imports: [ConfigModule.forRoot({
    load: [myconf],
    isGlobal: true,
  })],
  controllers: [AppController, StorageController, NyaaController],
  providers: [AppService, StorageService, TreeService, NyaaService],
})

export class AppModule { }
