const keyPrefix = 'dossier#';

function get(key, fallback) {
  let value;

  try {
    value = window.localStorage.getItem(keyPrefix + key);
  } catch (exc) {
    if (exc instanceof DOMException) {
      return fallback;
    }

    throw exc;
  }

  if (value === null) {
    return fallback;
  }

  try {
    return JSON.parse(value);
  } catch (exc) {
    return fallback;
  }
}

function set(key, value) {
  let strValue = JSON.stringify(value);

  try {
    window.localStorage.setItem(keyPrefix + key, strValue);
  } catch (exc) {
    if (!(exc instanceof DOMException)) {
      throw exc;
    }
  }
}

export { get, set };

/* vim: set sw=2 sts=2 et : */
