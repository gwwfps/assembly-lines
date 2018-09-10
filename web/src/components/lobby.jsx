import xs from 'xstream';
import sampleCombine from 'xstream/extra/sampleCombine';

import mvi from '../utils/mvi';
import inputValue from '../utils/input-value';
import PlayerStatus from '../constants/player-status';
import LobbyCard from '../views/lobby-card';

const isLobby = ({ Status }) => Status === PlayerStatus.LOBBY;

export default mvi({
  intent: ({ DOM, game: { state } }) => ({
    clickJoin$: DOM.select('button.join').events('click'),
    clickNew$: DOM.select('button.new').events('click'),
    nameInput$: inputValue(DOM.select('input[name="name"]')),
    sheetNameInput$: inputValue(DOM.select('input[name="sheetName"]')),
    state
  }),

  model: ({ clickJoin$, clickNew$, nameInput$, sheetNameInput$, state }) => ({
    join$: clickJoin$
      .compose(sampleCombine(nameInput$, sheetNameInput$))
      .map(([click, name, sheetName]) => ({
        id: click.target.value,
        name,
        sheetName
      }))
      .filter(x => x),
    new$: clickNew$
      .compose(sampleCombine(nameInput$, sheetNameInput$))
      .map(([_click, name, sheetName]) => ({ name, sheetName }))
      .filter(x => x),
    enabled$: xs
      .combine(nameInput$, sheetNameInput$)
      .map(([name, sheetName]) => !!(name && sheetName))
      .startWith(false),
    show$: state.map(isLobby),
    lobbies$: state.filter(isLobby).map(({ State }) => State)
  }),

  view: ({ join$, new$, enabled$, lobbies$, show$ }) => ({
    DOM: xs.combine(lobbies$, show$, enabled$).map(
      ([lobbies, show, enabled]) =>
        show ? (
          <div className="columns">
            {Object.keys(lobbies).length ? (
              Object.keys(lobbies).map(id =>
                LobbyCard({ id, players: lobbies[id], enabled })
              )
            ) : (
              <div className="column col-6">
                <div className="empty">
                  <p className="empty-title h5">
                    No lobby waiting for game start.
                  </p>
                  <p className="empty-subtitle">Start a new one to play.</p>
                </div>
              </div>
            )}
            <div className="column col-6">
              <div className="panel">
                <div className="panel-header">
                  <div className="panel-title">Player setup</div>
                </div>
                <div className="panel-body">
                  <div className="form-group">
                    <label className="form-label">Player name</label>
                    <input
                      type="text"
                      className="form-input"
                      name="name"
                      placeholder="Required"
                    />
                  </div>
                  <div className="form-group">
                    <label className="form-label">Factory name</label>
                    <input
                      type="text"
                      className="form-input"
                      name="sheetName"
                      placeholder="Required"
                    />
                  </div>
                </div>
                <div className="panel-footer">
                  <button className="btn new" disabled={!enabled}>
                    Start new lobby
                  </button>{' '}
                  or join existing one
                </div>
              </div>
            </div>
          </div>
        ) : null
    ),
    game: xs.merge(
      join$.map(args => ({
        action: 'JoinLobby',
        args
      })),
      new$.map(args => ({
        action: 'StartLobby',
        args
      }))
    )
  })
});
