import xs from 'xstream';

import mvi from './utils/mvi';

export default mvi({
  intent: ({ DOM }) => ({
    click$: DOM.events('click')
  }),

  model: ({ click$ }) => ({
    click$
  }),

  view: ({ click$ }) => ({
    DOM: xs.of(<div>init</div>),
    WS: click$.mapTo('some message')
  })
});
