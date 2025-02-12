import { writeFile } from 'node:fs/promises'
import { Injectable, Scope, OnModuleInit, OnModuleDestroy } from '@nestjs/common';
import { ConfigService } from '@nestjs/config'

import { LRUCache } from 'lru-cache'
import { has, set, unset, zipObject } from 'lodash-es'

import { IDatabase } from './tree/tree.d';

@Injectable({ scope: Scope.DEFAULT })
export class StorageService implements OnModuleInit, OnModuleDestroy {

    RootList: Map<string, string> = new Map()
    Database: any
    dbpool: LRUCache<string, IDatabase, object>
    _pendingWrites: object = {}

    constructor(
        private configService: ConfigService,
    ) { }

    async onModuleInit(): Promise<void> {
        const Database = require('better-sqlite3');
        this.Database = Database

        this.dbpool = new LRUCache({
            max: 10,
            maxSize: 50,
            ttl: 1000 * 60, /* 1 mins */
            updateAgeOnGet: true,
            updateAgeOnHas: true,

            sizeCalculation: (_value, _key) => 1,

            dispose: async (db: IDatabase, rootId: string) => {
                if (has(this._pendingWrites, rootId)) { return }

                const buffer = db.serialize()
                const pWriteFile = writeFile(this.getDbPath(rootId), buffer)

                db.close()
                await pWriteFile

                unset(this._pendingWrites, rootId)
            },
        })
    }

    async onModuleDestroy(): Promise<void> {
        console.log("storage.services onModuleDestroy")
    }

    registerRoot(rootId: string, label: string): Map<string, string> {
        return this.RootList.set(rootId, label)
    }

    getRoots(): Array<object> {
        return Array.from(this.RootList, (vals, index) => {
            return zipObject(['rootId', 'label'], vals)
        })
    }

    getBottom(): Array<object> {
        return Array.from(this.RootList, (vals, index) => {
            return zipObject(['rootId', 'label'], vals)
        })
    }

    async getDb(rootId: string) {
        if (this.dbpool.has(rootId)) {
            return this.dbpool.get(rootId)
        }

        const db: IDatabase = new this.Database(
            this.getDbPath(rootId),
            // { verbose: console.log },
        )
        db.pragma('journal_mode = WAL');

        this.dbpool.set(rootId, db)
        return db
    }

    async closeDb(rootId: string, db: IDatabase) {
        set(this._pendingWrites, rootId, true)
        return this.dbpool.delete(rootId)
    }

    getDbPath(rootId: string) {
        const rootPath = this.configService.get<string>('rootPath')
        return `${rootPath}/db/${rootId}.db`
    }
}
