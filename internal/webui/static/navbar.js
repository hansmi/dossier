import * as defaultunit from './defaultunit.js';
import * as geometry from './geometry.js';

const elDefaultUnit = document.getElementById('dossier-default-unit-select');

(() => {
  const def = defaultunit.get();
  let selectedIndex = 0;

  geometry.ALL_UNITS.forEach((unit) => {
    const opt = document.createElement("option");

    opt.value = unit.name;
    opt.text = unit.name;

    elDefaultUnit.add(opt);

    if (def.name === unit.name) {
      selectedIndex = opt.index;
    }
  });

  elDefaultUnit.selectedIndex = selectedIndex;
})();

elDefaultUnit.addEventListener("input", (ev) => {
  const selectedOptions = elDefaultUnit.selectedOptions;

  if (selectedOptions.length > 0) {
    const unit = geometry.LengthUnit.fromName(selectedOptions[0].value);

    defaultunit.set(unit);
  }
});

/* vim: set sw=2 sts=2 et : */
