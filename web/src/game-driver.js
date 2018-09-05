import xs from 'xstream';
import { adapt } from '@cycle/run/lib/adapt';

export function makeGameDriver(playerId) {
  const connection = new WebSocket(`ws://localhost:5555/ws/${playerId}`);

  return function(sink$) {
    sink$.addListener({
      next: msg => {
        const { action, argsJson } = msg;
        connection.send(`${action}|${JSON.stringify(argsJson)}`);
      },
      error: () => {},
      complete: () => {}
    });

    const source = xs.create({
      start: listener => {
        connection.onerror = err => {
          listener.error(err);
        };
        connection.onmessage = msg => {
          console.log(msg);
          listener.next(msg.data);
        };
      },
      stop: () => {
        this.connection.close();
      }
    });

    return adapt(source);
  };
}
