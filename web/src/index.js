import { run } from '@cycle/run';
import { makeDOMDriver } from '@cycle/dom';

import main from './app';
import { makeGameDriver } from './game-driver';
import { getId } from './utils/id';

run(main, {
  DOM: makeDOMDriver('#root'),
  game: makeGameDriver(getId())
});
