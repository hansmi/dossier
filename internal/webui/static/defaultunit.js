import * as geometry from './geometry.js';
import * as settings from './settings.js';

const SETTING_NAME = 'length_unit';
const EVENT_NAME = 'dossier-default-unit-changed';

function get() {
  const name = settings.get(SETTING_NAME, null);

  let result = null;

  if (name) {
    result = geometry.LengthUnit.fromName(name);
  }

  if (!(result && result.name)) {
    result = geometry.Point;
  }

  return result;
}

function set(unit) {
  const old = get();

  if (old === unit || old.name === unit.name) {
    return;
  }

  settings.set(SETTING_NAME, unit.name);

  const ev = new CustomEvent(EVENT_NAME, {
    detail: unit,
  });

  document.dispatchEvent(ev);
}

// Invoke the given function after the default length unit has changed. The
// returned function removes the installed event handler.
function observe(fn) {
  const handler = (ev) => {
    fn(ev.detail);
  };

  document.addEventListener(EVENT_NAME, handler, {
    passive: true,
  });

  return () => {
    document.removeEventListener(EVENT_NAME, handler);
  };
}

export {
  get,
  set,
  observe,
};

/* vim: set sw=2 sts=2 et : */
