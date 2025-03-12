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
    _pendingWrites: Map<string, any> = new Map()

    constructor(
        private configService: ConfigService,
    ) { }

    async onModuleInit(): Promise<void> {
        const Database = require('better-sqlite3');
        this.Database = Database

        this.dbpool = new LRUCache({
            max: 10,
            maxSize: 50,
            ttl: 1000 * 30, /* 30s */
            updateAgeOnGet: true,
            updateAgeOnHas: false,

            sizeCalculation: (_value, _key) => 1,

            dispose: (db: IDatabase, rootId: string) => {
                if (this._pendingWrites.has(rootId)) { return }

                this._pendingWrites.set(
                    rootId,
                    setTimeout(async () => {
                        const buffer = db.serialize()
                        const pWriteFile = writeFile(this.getDbPath(rootId), buffer)

                        db.close()
                        await pWriteFile

                        clearTimeout(this._pendingWrites.get(rootId))
                        this._pendingWrites.delete(rootId)
                        console.log(`db ${rootId} saved`)
                    }, 500)
                )
            },
        })
        console.log("storage.services onModuleInit")
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

    async getDb(rootId: string) {
        let db: IDatabase

        if (this.dbpool.has(rootId)) {
            db = this.dbpool.get(rootId)
        } else {
            db = new this.Database(
                this.getDbPath(rootId),
                // { verbose: console.log },
            )
            this.dbpool.set(rootId, db)
        }
        db.pragma('journal_mode = WAL');

        return db
    }

    async closeDb(rootId: string, db: IDatabase) {
        return this.dbpool.delete(rootId)
    }

    getDbPath(rootId: string) {
        const rootPath = this.configService.get<string>('rootPath')
        return `${rootPath}/db/${rootId}.db`
    }
}
