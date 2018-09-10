import xs from 'xstream';
import sampleCombine from 'xstream/extra/sampleCombine';

import mvi from '../utils/mvi';
import PlayerStatus from '../constants/player-status';
import GamePhase from '../constants/game-phase';

const renderLobby = (name, { Players }) => (
  <div className="columns">
    <div className="column col-12">
      <div className="panel">
        <div className="panel-header">
          <div className="panel-title">Waiting to start</div>
        </div>
        <div className="panel-body">
          Players ({Object.keys(Players).length}
          ):{' '}
          {Object.keys(Players).map(player => (
            <span className="chip">
              {player === name && <i className="icon icon-people" />} {player}
            </span>
          ))}
        </div>
        <div className="panel-footer">
          <button className="btn btn-primary start">Start game</button>{' '}
          <button className="btn btn-error leave">Leave lobby</button>
        </div>
      </div>
    </div>
  </div>
);

export default mvi({
  intent: ({ DOM, game: { state } }) => ({
    clickStart$: DOM.select('button.start').events('click'),
    clickLeave$: DOM.select('button.leave').events('click'),
    inGame$: state.map(({ Status }) => Status === PlayerStatus.IN_GAME),
    state$: state.map(({ State }) => State),
    name$: state.map(({ Name }) => Name)
  }),

  model: ({ clickStart$, clickLeave$, inGame$, state$, name$ }) => ({
    start$: clickStart$,
    leave$: clickLeave$,
    show$: inGame$,
    name$,
    inLobby$: state$.map(({ Phase }) => Phase === GamePhase.LOBBY),
    state$
  }),

  view: ({ start$, leave$, inLobby$, show$, state$, name$ }) => ({
    DOM: xs
      .combine(inLobby$, show$, state$, name$)
      .map(
        ([inLobby, show, state, name]) =>
          show ? (inLobby ? renderLobby(name, state) : null) : null
      ),
    game: xs.merge(
      start$.map(() => ({
        action: 'StartGame'
      })),
      leave$.map(() => ({
        action: 'LeaveLobby'
      }))
    )
  })
});
