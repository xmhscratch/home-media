import { Test, TestingModule } from '@nestjs/testing';
import { NyaaController } from './nyaa.controller';

describe('NyaaController', () => {
  let controller: NyaaController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [NyaaController],
    }).compile();

    controller = module.get<NyaaController>(NyaaController);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });
});
