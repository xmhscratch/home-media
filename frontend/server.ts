import { NestFactory } from '@nestjs/core';

import { APP_BASE_HREF } from '@angular/common';
import { CommonEngine } from '@angular/ssr';

import express from 'express';
import { parseInt as ldParseInt } from 'lodash-es';

import { fileURLToPath } from 'node:url';
import { dirname, join, resolve } from 'node:path';

import { AppModule } from './app.module';
import bootstrap from './main.server';

// The Express app is exported so that it can be used by serverless Functions.
export async function app(): Promise<express.Express> {
    const app = await NestFactory.create(AppModule);
    const httpAdapter = app.getHttpAdapter();
    const server = httpAdapter.getInstance();

    const serverDistFolder = dirname(fileURLToPath(import.meta.url));
    const browserDistFolder = resolve(serverDistFolder, '../browser');
    const indexHtml = join(serverDistFolder, 'index.server.html');

    const commonEngine = new CommonEngine();

    server.set('view engine', 'html');
    server.set('views', browserDistFolder);

    // Example Express Rest API endpoints
    // server.get('/api/**', (req, res) => { });
    // Serve static files from /browser
    server.get('**', express.static(browserDistFolder, {
        maxAge: '1y',
        index: 'index.html',
    }));

    // All regular routes use the Angular engine
    server.get('**', (req: any, res: any, next: (err?: Error) => void) => {
        const { protocol, originalUrl, baseUrl, headers } = req;

        commonEngine
            .render({
                bootstrap,
                documentFilePath: indexHtml,
                url: `${protocol}://${headers.host}${originalUrl}`,
                publicPath: browserDistFolder,
                providers: [
                    { provide: APP_BASE_HREF, useValue: baseUrl },
                ],
            })
            .then((html) => res.send(html))
            .catch((err) => next(err));
    });

    return server;
}

async function run(): Promise<void> {
    const port = ldParseInt(process.env['PORT'] || '4000')

    // Start up the Node server
    const server = await app();
    server.listen(port, '0.0.0.0', () => {
        console.log(`Node Express server listening on http://0.0.0.0:${port}`);
    });
}

run();
