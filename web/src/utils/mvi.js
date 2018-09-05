import isolate from '@cycle/isolate';

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
    const state$ = model(actions, componentSinks);
    const sinks = view(state$, componentSinks);
    return sinks;
  };
}
