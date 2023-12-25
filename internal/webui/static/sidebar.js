import * as settings from './settings.js';

const WIDTH_SETTING_NAME = 'sidebar_width';

const elSidebar = document.getElementById('sidebar');

function initSidebar() {
  let timeout = false;

  const observer = new ResizeObserver((entries) => {
    window.clearTimeout(timeout);

    for (const entry of entries) {
      timeout = window.setTimeout(() => {
        settings.set(WIDTH_SETTING_NAME, entry.borderBoxSize[0].inlineSize);
      }, 250);
      break;
    }
  });

  observer.observe(elSidebar, {
    box: 'content-box',
  });

  elSidebar.style.resize = 'horizontal';

  let current = settings.get(WIDTH_SETTING_NAME, null);
  if (current !== null && current > 1) {
    elSidebar.style.width = `${current}px`;
  }
}

initSidebar();

/* vim: set sw=2 sts=2 et : */
