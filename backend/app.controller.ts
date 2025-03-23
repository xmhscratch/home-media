import { Controller } from '@nestjs/common';
import { Get, Res } from '@nestjs/common';
import { Response } from 'express';
// import { ConfigService } from '@nestjs/config';

import { AppService } from './app.service';

@Controller()
export class AppController {
  constructor(
    // private readonly configService: ConfigService,
    private readonly appService: AppService,
  ) { }

  @Get('/ping')
  async ping(
    @Res() res: Response,
  ) {
    return res.json({})
  }
}
