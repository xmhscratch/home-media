// import { readFile } from 'node:fs/promises';
// import { createReadStream } from 'node:fs';

import { Controller } from '@nestjs/common';
import { Get, Post, Param, Header, Res, Query, RawBody } from '@nestjs/common';
import { Response } from 'express';

import { StorageService } from '../storage/storage.service';
import { NyaaService } from './nyaa.service';

@Controller('nyaa')
export class NyaaController {

    constructor(
        private nyaaService: NyaaService,
        private storageService: StorageService,
    ) { }

    @Get('/update')
    @Header('Content-Type', 'application/json')
    @Header('Accept', 'application/json')
    async getUpToDate(
        @Query('label') label: string,
        @Query('rootId') rootId: string,
        @Query('query') query: string,
        @Res() res: Response,
    ) {
        this.storageService.registerRoot(rootId, label)
        const a = await this.nyaaService.importItems(rootId, label, query)
        return res.json({})
        // return res.json({ rootId, label, query })
    }

    @Get('/test')
    @Header('Content-Type', 'application/json')
    @Header('Accept', 'application/json')
    async getTest(
        @Res() res: Response,
    ) {
        const label = "SubsPlease"
        const rootId = "67a61a61f2f1b1cbb8164e08"
        const query = "f=0&c=1_2&q=%5BSubsPlease%5D+1080p"
        this.storageService.registerRoot(rootId, label)
        const a = await this.nyaaService.importItems(rootId, label, query)
        return res.json(a)
    }

    // @Get('sql-wasm.wasm')
    // @Header('Content-Type', 'application/wasm')
    // @Header('Content-Disposition', 'attachment; filename="sql-wasm.wasm"')
    // async getInfo(@Res() res) {
    //     return createReadStream(require.resolve('sql.js/dist/sql-wasm.wasm')).pipe(res);
    // }
}
