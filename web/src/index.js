import 'babel-polyfill';

import { run } from '@cycle/run';
import { makeDOMDriver } from '@cycle/dom';

import main from './components/app';
import { makeGameDriver } from './game-driver';
import { getId } from './utils/id';

(async () => {
  run(main, {
    DOM: makeDOMDriver('#root'),
    game: await makeGameDriver(getId())
  });
})();
