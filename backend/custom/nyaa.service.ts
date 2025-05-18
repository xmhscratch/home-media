import { Injectable, OnModuleInit, OnModuleDestroy } from '@nestjs/common';

import { map, pick, join } from 'lodash-es';
import Parser from 'rss-parser';

import { Tree } from '../storage/tree/tree';
import { TreeMetaNode, IDatabase, IStatement, BindParams, SqlValue } from '../storage/tree/tree.d';
import { TreeService } from '../storage/tree.service';

type NyaaItem = {
    id: string
    title: string
    seeders: string
    leechers: string
    infoHash: string
    size: string
    trusted: string
    isoDate: string
}

type ExNyaaItem = NyaaItem & {
    seeders: number
    leechers: number
    trusted: boolean
    createdDate: string
    dataSource: string
    dataSourceType: number
}

type ExTreeMetaNode = TreeMetaNode & {
    seeders: number
    leechers: number
    // agoAsMinutes: number
}

@Injectable()
export class NyaaService implements OnModuleInit, OnModuleDestroy {

    parser: Parser<any> = new Parser({
        customFields: {
            item: [
                ['nyaa:seeders', 'seeders'],
                ['nyaa:leechers', 'leechers'],
                ['nyaa:infoHash', 'infoHash'],
                ['nyaa:size', 'size'],
                ['nyaa:trusted', 'trusted'],
            ]
        }
    });

    trackers = [
        'http://nyaa.tracker.wf:7777/announce',
        'udp://open.stealth.si:80/announce',
        'udp://tracker.opentrackr.org:1337/announce',
        'udp://exodus.desync.com:6969/announce',
        'udp://tracker.torrent.eu.org:451/announce',
    ];

    constructor(private treeService: TreeService) { }

    async onModuleInit(): Promise<void> { }

    async onModuleDestroy(): Promise<void> { }

    async importItems(rootId: string, _label: string, query: string): Promise<object> {
        const items = await this.fetchItems(query)

        const metaItems = map(items, (i: ExNyaaItem, _k: number): ExTreeMetaNode => {
            const { id, seeders, leechers, size, trusted } = i
            const { title, createdDate, dataSource, dataSourceType } = i

            const description = JSON.stringify({ seeders, leechers, size, trusted })
            return {
                ...<TreeMetaNode>{
                    id,
                    root: rootId,
                    description,
                    title,
                    createdDate,
                    dataSource,
                    dataSourceType,
                },
                seeders,
                leechers,
                // agoAsMinutes: (<any>new Date() - <any>new Date(createdDate)) / (1000 * 60),
            } as ExTreeMetaNode
        })

        await this.cleanOlderItems(rootId)

        await Promise.all(
            map(items, (item: any) => this.treeService.execute(
                rootId,
                'create',
                [rootId, item.id],
            ))
        )

        await Promise.all(
            map(metaItems, (metaItem: any) => this.treeService.execute(
                rootId,
                'insert',
                [metaItem],
            ))
        )

        return items
    }

    async fetchItems(query: string): Promise<Array<object>> {
        const feed = await this.parser
            .parseURL(`https://nyaa.si/?page=rss&${query || ''}`)
            .catch(console.log)

        const { items } = feed

        const t1 = map(items, (i: NyaaItem) => pick(i, [
            'title', 'seeders', 'leechers', 'infoHash', 'size', 'trusted', 'isoDate'
        ]))
        const t2 = map(t1, (v: NyaaItem, _k: number): ExNyaaItem => {
            const {
                infoHash,
                title, seeders, leechers,
                size, trusted, isoDate,
            } = v

            return <ExNyaaItem>{
                id: this.genId(infoHash),
                title,
                infoHash,
                seeders: parseInt(seeders),
                leechers: parseInt(leechers),
                size,
                trusted: !/(no|0|false)/gi.test(trusted),
                createdDate: (new Date(isoDate)).toISOString(),
                dataSource: `magnet:?xt=urn:btih:${infoHash}&dn=${title}&tr=${join(this.trackers, '&tr=')}`,
                dataSourceType: 2,
            }
        })

        return t2
    }

    cleanOlderItems(rootId: string) {
        return this.treeService.customEx(rootId, async (tree: Tree, db: IDatabase) => {
            let stmt: IStatement = (<IDatabase>db).prepare(`
                SELECT
                    node.id
                FROM
                    nodes AS node
                LEFT JOIN
                    nodes_meta AS meta
                ON (
                    node.id = meta.id
                )
                WHERE (
                    node.root = $rootId
                    AND ((strftime('%s','now') - strftime('%s',meta.created_date)) / 60) >= (60 * 24 * 5)
                )
                ORDER BY node.left;
            `)
            stmt.bind(<BindParams>{
                rootId: <SqlValue>rootId,
            })

            const cleanIds: string[] = []
            const targetIds: string[] = map(stmt.all(), (v: { id: string }): string => v['id'])

            for (let targetId of targetIds) {
                await tree.delete(targetId)
                cleanIds.push(targetId)
            }

            return cleanIds
        })
    }

    genId(infoHash: string): string {
        let txt = ""
        const primes = [5, 7, 11, 13, 17, 19, 23]
        for (let i = 0; i < infoHash.length; i++) {
            for (let j = 0; j < primes.length; j++) {
                if (!((i + 1) % primes[j])) { txt += infoHash[i]; continue }
            }
        }
        return txt
    }
}

// getSizeInt(ustr: string): number | BigInt {
//     const s = ['b', 'k', 'm', 'g', 't', 'p', 'e', 'z', 'y']
//     const us = /[\ ]{0,}([\d\.]+)[\ ]{0,}(bytes|b|([kmgtpezy]{1}[i]{0,1})b(yte|)[s]{0,1}){1}[\ ]{0,}/gi.exec(ustr)
//     const k = s.indexOf((us.length>3 ? (us[3][0]).toLowerCase() : 'b') || 'b')
//     return !k ? parseInt(us[1]) : BigInt(parseInt(us[1]) * (8<<7<<((k - 1) * 10)))
// }
