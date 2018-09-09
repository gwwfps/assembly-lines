import xs from 'xstream';

const errorPrefix = 'error|';

const openWS = playerId =>
  new Promise(resolve => {
    const connection = new WebSocket(`ws://localhost:5555/ws/${playerId}`);
    connection.onopen = () => {
      resolve(connection);
    };
  });

const transformArgs = args =>
  JSON.stringify(
    Object.keys(args).reduce(
      (acc, key) => ({
        ...acc,
        [`${key[0].toUpperCase()}${key.substr(1)}`]: args[key]
      }),
      {}
    )
  );

export async function makeGameDriver(playerId) {
  const connection = await openWS(playerId);

  return function(sink$) {
    sink$.addListener({
      next: msg => {
        const { action, args } = msg;
        const payload = args ? `${action}|${transformArgs(args)}` : action;
        connection.send(payload);
      },
      error: () => {},
      complete: () => {}
    });

    let stateListener, errorListener;
    const stateSource = xs.create({
      start: listener => {
        connection.onerror = err => {
          listener.error(err);
        };
        stateListener = listener;
      },
      stop: () => {
        connection.close();
      }
    });
    const errorSource = xs.create({
      start: listener => {
        errorListener = listener;
      },
      stop: () => {}
    });

    connection.onmessage = ({ data }) => {
      if (data.startsWith(errorPrefix)) {
        errorListener && errorListener.next(data.substring(errorPrefix.length));
      } else {
        stateListener && stateListener.next(JSON.parse(data));
      }
    };

    return {
      error: errorSource,
      state: stateSource
    };
  };
}
