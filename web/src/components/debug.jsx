import sampleCombine from 'xstream/extra/sampleCombine';

import mvi from '../utils/mvi';
import inputValue from '../utils/input-value';

export default mvi({
  intent: ({ DOM, game: { state } }) => ({
    click$: DOM.select('button').events('click'),
    actionInput$: inputValue(DOM.select('input[name="action"]')),
    argsInput$: inputValue(DOM.select('input[name="args"]')),
    state
  }),

  model: ({ click$, actionInput$, argsInput$, state }) => ({
    action$: click$
      .compose(sampleCombine(actionInput$, argsInput$))
      .map(([_click, action, args]) => action && args && { action, args })
      .filter(x => x),
    message$: state
  }),

  view: ({ action$, message$ }) => ({
    DOM: message$.map(data => (
      <div className="debug">
        <pre className="code">
          <code>{JSON.stringify(data)}</code>
        </pre>
        <input className="form-input" name="action" />
        <input className="form-input" name="args" />
        <button className="btn">Send</button>
      </div>
    )),
    game: action$.map(({ action, args }) => ({
      action: action,
      args: JSON.parse(args)
    }))
  })
});
