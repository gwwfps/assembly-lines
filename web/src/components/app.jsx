import xs from 'xstream';

import mvi from '../utils/mvi';
import Debug from './debug';
import Lobby from './lobby';

export default mvi({
  view: () => ({
    DOM: ({ components$ }) =>
      components$.map(({ debug, lobby }) => (
        <div>
          <div className="container grid-lg">
            <div className="columns">
              <div className="column col-12">
                <h3>Assembly Lines the Game</h3>
              </div>
            </div>
            {lobby}
          </div>
          {debug}
        </div>
      )),
    game: ({ debug$, lobby$ }) =>
      xs.merge(debug$, lobby$).startWith({
        action: 'FetchState'
      })
  }),

  components: {
    debug: {
      main: Debug
    },
    lobby: {
      main: Lobby
    }
  }
});
