import isolate from '@cycle/isolate';
import xs from 'xstream';

const identity = x => x;

function render({ main, props }, sources) {
  return isolate(main)({ ...sources, props });
}

export default function makeMVIComponent({
  model = identity,
  view,
  intent = identity,
  components = {}
}) {
  return function(sources) {
    const componentSinks = Object.keys(components).reduce((acc, key) => {
      const component = components[key];
      const sinks = Array.isArray(component)
        ? component.map(c => render(c, sources))
        : render(component, sources);
      return { ...acc, [key]: sinks };
    }, {});

    const actions = intent(sources, componentSinks);
    const state = model(actions, componentSinks);
    const sinkCreators = view(state, componentSinks);

    const sinks = Object.keys(sinkCreators).reduce((acc, type) => {
      let combinedSink;
      const create = sinkCreators[type];
      if (typeof create === 'function') {
        const keys = Object.keys(componentSinks).filter(
          key => componentSinks[key][type]
        );
        const componentSinksForType = keys.reduce(
          (accum, key) => ({
            ...accum,
            [`${key}$`]: componentSinks[key][type]
          }),
          {}
        );
        if (keys.length > 1) {
          componentSinksForType.components$ = xs
            .combine(
              ...Object.values(componentSinksForType).map(sink =>
                sink.startWith(undefined)
              )
            )
            .map(values =>
              values.reduce(
                (accum, value, i) => ({
                  ...accum,
                  [keys[i]]: value
                }),
                {}
              )
            );
        }
        combinedSink = create(componentSinksForType);
      } else {
        combinedSink = create;
      }

      return { ...acc, [type]: combinedSink };
    }, {});

    return sinks;
  };
}
