import xs from 'xstream';

import mvi from '../utils/mvi';
import Debug from './debug';
import Lobby from './lobby';
import Game from './game';

export default mvi({
  view: () => ({
    DOM: ({ components$ }) =>
      components$.map(({ debug, lobby, game }) => (
        <div>
          <div className="container grid-lg">
            <div className="columns">
              <div className="column col-12">
                <h3>Assembly Lines the Game</h3>
              </div>
            </div>
            {lobby}
            {game}
          </div>
          {debug}
        </div>
      )),
    game: ({ debug$, lobby$, game$ }) =>
      xs.merge(debug$, lobby$, game$).startWith({
        action: 'FetchState'
      })
  }),

  components: {
    debug: {
      main: Debug
    },
    lobby: {
      main: Lobby
    },
    game: {
      main: Game
    }
  }
});
