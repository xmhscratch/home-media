import { Controller } from '@nestjs/common';
import { Get, Param, Res, Query } from '@nestjs/common';
import { Response } from 'express';

import { TreeService } from './tree.service';
import { StorageService } from './storage.service';
import { default as sampleTreeData } from './sample'

import { map } from 'lodash-es';

@Controller('storage')
export class StorageController {

    constructor(
        private treeService: TreeService,
        private storageService: StorageService,
    ) { }

    @Get(':rootId([a-f\\d]{24})/:nodeId([a-f\\d]{24})/node')
    async getInfo(
        @Param('rootId') rootId: string,
        @Param('nodeId') nodeId: string,
        @Query('extend') includeMeta: string,
        @Query('branch') onlyBranch: string,
        @Query('t') timestamp: string,
        @Res() res: Response,
    ) {
        const results = await this.treeService.execute(
            rootId,
            'getNode',
            [nodeId, {
                includeMeta: /(1|true|yes)/g.test(includeMeta) ? 1 : 0,
                onlyBranch: /(1|true|yes)/g.test(onlyBranch) ? 1 : 0,
            }],
        )
        return res.json(results)
    }

    @Get(':rootId([a-f\\d]{24})/:nodeId([a-f\\d]{24})/children')
    async getChildren(
        @Param('rootId') rootId: string,
        @Param('nodeId') nodeId: string,
        @Query('extend') includeMeta: string,
        @Query('branch') onlyBranch: string,
        @Query('t') timestamp: string,
        @Res() res: Response,
    ) {
        const results = await this.treeService.execute(
            rootId,
            'getChildren',
            [nodeId, {
                includeMeta: /(1|true|yes)/g.test(includeMeta) ? 1 : 0,
                onlyBranch: /(1|true|yes)/g.test(onlyBranch) ? 1 : 0,
            }],
        )
        return res.json(results)
    }

    @Get(':rootId([a-f\\d]{24})/:nodeId([a-f\\d]{24})/descendants')
    async getDescendants(
        @Param('rootId') rootId: string,
        @Param('nodeId') nodeId: string,
        @Query('extend') includeMeta: string,
        @Query('branch') onlyBranch: string,
        @Query('t') timestamp: string,
        @Res() res: Response,
    ) {
        const results = await this.treeService.execute(
            rootId,
            'getDescendants',
            [nodeId, {
                includeMeta: /(1|true|yes)/g.test(includeMeta) ? 1 : 0,
                onlyBranch: /(1|true|yes)/g.test(onlyBranch) ? 1 : 0,
            }],
        )
        return res.json(results)
    }

    @Get(':rootId([a-f\\d]{24})/:nodeId([a-f\\d]{24})/paths')
    async getPaths(
        @Param('rootId') rootId: string,
        @Param('nodeId') nodeId: string,
        @Query('extend') includeMeta: string,
        @Query('branch') onlyBranch: string,
        @Query('t') timestamp: string,
        @Res() res: Response,
    ) {
        const results = await this.treeService.execute(
            rootId,
            'getPaths',
            [nodeId, {
                includeMeta: /(1|true|yes)/g.test(includeMeta) ? 1 : 0,
                onlyBranch: /(1|true|yes)/g.test(onlyBranch) ? 1 : 0,
            }],
        )
        return res.json(results)
    }

    @Get(':rootId([a-f\\d]{24})/:nodeId([a-f\\d]{24})/adjacency')
    async getAdjacencyList(
        @Param('rootId') rootId: string,
        @Param('nodeId') nodeId: string,
        @Query('extend') includeMeta: string,
        @Query('branch') onlyBranch: string,
        @Query('t') timestamp: string,
        @Res() res: Response,
    ) {
        const results = await this.treeService.execute(
            rootId,
            'toAdjacencyList',
            [nodeId, {
                includeMeta: /(1|true|yes)/g.test(includeMeta) ? 1 : 0,
                onlyBranch: /(1|true|yes)/g.test(onlyBranch) ? 1 : 0,
            }],
        )
        return res.json(results)
    }

    @Get('roots')
    async getRoots(
        @Res() res: Response,
    ) {
        const rootItems = this.storageService.getRoots()
        return res.json(rootItems)
    }

    @Get(':rootId([a-f\\d]{24})/test')
    async getTest(
        @Param('rootId') rootId: string,
        @Res() res: Response,
    ) {
        const results = await this.treeService.execute(
            rootId,
            'import',
            [sampleTreeData as object[]],
        )
        return res.json(results)
    }
}
