export default (...values) =>
  values.reduce((acc, value, i) => ({ ...acc, [value]: i }), {});
