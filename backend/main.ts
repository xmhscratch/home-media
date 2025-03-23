import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';

async function bootstrap() {
    const port = parseInt((process.env['PORT'] || '4000')) + 100;

    const app = await NestFactory.create(AppModule);

    app.enableCors();
    app.enableShutdownHooks();

    await app.listen(port, () => {
        // console.log(`NestJS server listening on http://localhost:${port}`)
    });
}
bootstrap();
