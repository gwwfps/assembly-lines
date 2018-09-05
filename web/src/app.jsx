import xs from 'xstream';
import sampleCombine from 'xstream/extra/sampleCombine';

import mvi from './utils/mvi';
import Init from './init';

export default mvi({
  intent: ({ DOM, game }) => ({
    click$: DOM.select('button').events('click'),
    input$: DOM.select('input')
      .events('input')
      .map(ev => ev.target.value),
    game
  }),

  // model: ({ click$ }) => {
  //   click$;
  // },

  view: ({ click$, input$, game }) => ({
    DOM: game.startWith('abc').map(msg => (
      <div>
        <h1>Assembly Lines the Game</h1>
        <p>msg: {msg}</p>
        <input />
        <button>Send</button>
      </div>
    )),
    game: click$
      .compose(sampleCombine(input$))
      .map(([click, input]) => input)
      .filter(x => x)
      .map(action => ({
        action,
        argsJson: {
          Question: 0
        }
      }))
  }),

  components: {
    init: {
      main: Init
    }
  }
});
