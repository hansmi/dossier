import * as settings from './settings.js';
import * as geometry from './geometry.js';

const kShowKindSetting = 'show_kind';
const kShowEmptySetting = 'show_empty';
const kShowSketchNodeSetting = 'show_sketch_node';

const divSidebar = document.getElementById('sidebar');

const cbShowLayout = document.getElementById('page_filter_show_layout');
const divShowKindGroup = document.getElementById('page_filter_show_kind_group');
const rbShowNone = document.getElementById('page_filter_show_none');
const rbShowBlocks = document.getElementById('page_filter_show_blocks');
const rbShowLines = document.getElementById('page_filter_show_lines');
const cbShowEmpty = document.getElementById('page_filter_show_empty');
const cbSketchShowValid = document.getElementById('sketch_show_valid');

const divViewer = document.getElementById('dossier_viewer');

const divNodeDialogTemplate = document.getElementById('dossier_page_node_dialog_template');
const elNodeDialogKind = document.getElementById('dossier_page_node_dialog_kind');
const elNodeDialogBounds = document.getElementById('dossier_page_node_dialog_bounds');
const elNodeDialogText = document.getElementById('dossier_page_node_dialog_text');
const elNodeDialogTextDetails = document.getElementById('dossier_page_node_dialog_text_details');

const kDocBlockClass = 'overlay_doc_block';
const kDocLineClass = 'overlay_doc_line';
const kBlocksVisibleClass = 'dossier_doc_blocks_visible';
const kLinesVisibleClass = 'dossier_doc_lines_visible';
const kEmptyVisibleClass = 'dossier_doc_empty_element_visible';
const kSketchNodeVisibleClass = 'dossier_sketch_node_visible';
const kSketchNodeInfoClass = 'dossier_sketch_node_info';
const kSketchNodeHighlightClass = 'dossier_sketch_node_highlight';
const kSketchNodeSearchAreaHighlightClass = 'dossier_sketch_node_search_area_highlight';

function applyShowLayout() {
  const toggle = (token, force) => {
    divViewer.classList.toggle(token, force);
  };

  toggle(kBlocksVisibleClass, rbShowBlocks.checked);
  toggle(kLinesVisibleClass, rbShowLines.checked);
  toggle(kEmptyVisibleClass, cbShowEmpty.checked);
  toggle(kSketchNodeVisibleClass, cbSketchShowValid.checked);
}

function initFilter() {
  cbShowEmpty.checked = settings.get(kShowEmptySetting, false);
  cbShowEmpty.addEventListener('change', (ev) => {
    settings.set(kShowEmptySetting, ev.target.checked);
    applyShowLayout();
  });

  const selectedValue = settings.get(kShowKindSetting, null);

  let selected = [
    rbShowLines,
    rbShowBlocks,
  ].find((i) => i.value === selectedValue);

  if (selected === undefined) {
    selected = rbShowNone;
  }

  selected.checked = true;

  divShowKindGroup.addEventListener('change', (ev) => {
    const rbChecked = divShowKindGroup.querySelector('input[name="page_filter_show_kind"]:checked');

    if (rbChecked) {
      settings.set(kShowKindSetting, rbChecked.value);
    }

    applyShowLayout();
  });

  cbSketchShowValid.checked = settings.get(kShowSketchNodeSetting, true);
  cbSketchShowValid.addEventListener('change', (ev) => {
    settings.set(kShowSketchNodeSetting, ev.target.checked);
    applyShowLayout();
  });

  applyShowLayout();
}

class PositionTool {
  _infoPopper = null;
  _infoWrapper = null;
  _info = null;

  constructor(element) {
    this._elem = element;
    this._physWidth = new geometry.Length(this._elem.dataset.widthPt, geometry.Point);
    this._physHeight = new geometry.Length(this._elem.dataset.heightPt, geometry.Point);

    this._ensurePointer();
    this._ensureInfo();
  }

  _ensurePointer() {
    this._elPointer = document.createElement('div');
    this._elPointer.classList.add('position-absolute', 'invisible');
    this._elem.appendChild(this._elPointer);
  }

  _ensureInfo() {
    const wrapper = document.createElement('div');

    wrapper.innerHTML = `
      <div class="tooltip dossier_viewer_tooltip dossier_viewer_pointer_tooltip" role="tooltip">
        <div class="tooltip-inner">
          <dossier-geometry-point></dossier-geometry-point>
        </div>
      </div>
    `;

    this._infoWrapper = wrapper.children[0];
    this._info = this._infoWrapper.querySelector('dossier-geometry-point');

    this._elem.appendChild(this._infoWrapper);
  }

  activate() {
    this._elem.addEventListener('pointerenter', this._onPointerEnter.bind(this), {
      passive: true,
    });

    this._elem.addEventListener('pointermove', this._onPointerMove.bind(this), {
      passive: true,
    });

    this._elem.addEventListener('pointerleave', this._onPointerLeave.bind(this), {
      passive: true,
    });
  }

  _createInfoPopper() {
    return window.Popper.createPopper(this._elPointer, this._infoWrapper, {
      placement: 'right-start',
      modifiers: [
        {
          name: 'offset',
          options: {
            offset: [10, 10],
          },
        },
      ],
    });
  }

  _eventData(ev) {
    const elemBounds = this._elem.getBoundingClientRect();
    const posAbs = new DOMPointReadOnly(
      ev.clientX - elemBounds.x,
      ev.clientY - elemBounds.y,
    );
    const posPct = new DOMPointReadOnly(
      posAbs.x * 100.0 / elemBounds.width,
      posAbs.y * 100.0 / elemBounds.height,
    );

    return { elemBounds, posAbs, posPct };
  }

  _onPointerEnter(ev) {
    if (this._infoPopper === null) {
      this._infoPopper = this._createInfoPopper();
      this._infoWrapper.classList.add('show');
    }
  }

  _onPointerMove(ev) {
    if (this._infoPopper !== null) {
      const data = this._eventData(ev);

      const info = this._info;

      info.setAttribute('left-pt', data.posPct.x * this._physWidth.pt / 100);
      info.setAttribute('top-pt', data.posPct.y * this._physHeight.pt / 100);

      const elPointerStyle = this._elPointer.style;

      elPointerStyle.left = `${data.posPct.x}%`;
      elPointerStyle.top = `${data.posPct.y}%`;

      this._infoPopper.update();
    }
  }

  _onPointerLeave(ev) {
    this._infoWrapper.classList.remove('show');

    if (this._infoPopper !== null) {
      this._infoPopper.destroy();
      this._infoPopper = null;
    }
  }
}

class MeasurementTool {
  _listenerAbortController = null;
  _origin = null;
  _frame = null;
  _infoPopper = null;
  _infoPopperPlacement = null;
  _infoWrapper = null;
  _info = null;

  constructor(element) {
    this._elem = element;
    this._physWidth = new geometry.Length(this._elem.dataset.widthPt, geometry.Point);
    this._physHeight = new geometry.Length(this._elem.dataset.heightPt, geometry.Point);

    this._ensureFrame();
    this._ensureInfo();

    this._sm = new window.StateMachine({
      init: 'hidden',
      transitions: [
        { name: 'down', from: ['hidden', 'visible'], to: 'movetest' },

        { name: 'move', from: ['movetest', 'measure'], to: 'measure' },

        { name: 'up', from: 'movetest', to: 'hidden' },
        { name: 'up', from: 'measure', to: 'visible' },

        { name: 'cancel', from: ['down', 'movetest', 'measure'], to: 'hidden' },
      ],
      methods: {
        // onTransition: (lifecycle) => console.log("transition", lifecycle),
        // onEnterState: (lifecycle) => console.log("enter state", lifecycle),

        onEnterHidden: this._onEnterHidden.bind(this),
        onBeforeDown: this._onBeforeDown.bind(this),
        onLeaveMovetest: this._onLeaveMovetest.bind(this),
        onEnterMeasure: this._onEnterMeasure.bind(this),
        onAfterMove: this._onAfterMove.bind(this),
        onAfterUp: this._onAfterUp.bind(this),
        onAfterCancel: this._onAfterCancel.bind(this),
      },
    });
  }

  _ensureFrame() {
    const frame = document.createElement('div');

    frame.classList.add('dossier_viewer_selection', 'text-body');

    this._frame = frame;

    this._elem.appendChild(frame);
  }

  _ensureInfo() {
    const wrapper = document.createElement('div');

    wrapper.innerHTML = `
      <div class="tooltip dossier_viewer_tooltip dossier_viewer_selection_tooltip" role="tooltip">
        <div class="tooltip-inner">
          <dossier-geometry-rect></dossier-geometry-rect>
        </div>
      </div>
    `;

    this._infoWrapper = wrapper.children[0];
    this._info = this._infoWrapper.querySelector('dossier-geometry-rect');

    this._elem.appendChild(this._infoWrapper);
  }

  activate() {
    this._elem.addEventListener('pointerdown', this._onPointerDown.bind(this), {
      capture: true,
    });
  }

  _eventData(ev) {
    const elemBounds = this._elem.getBoundingClientRect();
    const posAbs = new DOMPointReadOnly(
      ev.clientX - elemBounds.x,
      ev.clientY - elemBounds.y,
    );
    const posPct = new DOMPointReadOnly(
      posAbs.x * 100.0 / elemBounds.width,
      posAbs.y * 100.0 / elemBounds.height,
    );

    return {
      target: ev.target,
      elemBounds,
      posAbs,
      posPct,
    };
  }

  _createInfoPopper() {
    return window.Popper.createPopper(this._frame, this._infoWrapper, {
      modifiers: [
        {
          name: 'flip',
          options: {
            boundary: this._elem,
            flipVariations: false,
            allowedAutoPlacements: [],
          },
        },
      ],
    });
  }

  _updateInfoPopperPlacement(data) {
    const placement =
      (this._origin.posPct.y < data.posPct.y ? 'top' : 'bottom') + '-' +
      (this._origin.posPct.x < data.posPct.x ? 'start' : 'end');

    if (placement !== this._infoPopperPlacement) {
      this._infoPopper.setOptions(options => ({
        ...options,
        placement: placement,
      }));
      this._infoPopperPlacement = placement;
    }
  }

  _setVisibility(value) {
    this._frame.classList.toggle('d-none', !value);
    this._frame.classList.toggle('d-block', value);
    this._infoWrapper.classList.toggle('show', value);
  }

  _onEnterHidden(lifecycle) {
    if (this._infoPopper !== null) {
      this._infoPopper.destroy();
      this._infoPopper = null;
    }

    this._setVisibility(false);
  }

  _onBeforeDown(lifecycle, data) {
    if (this._listenerAbortController === null) {
      this._listenerAbortController = new AbortController();
    }

    this._elem.addEventListener('pointercancel', this._onPointerCancel.bind(this), {
      passive: true,
      signal: this._listenerAbortController.signal,
    });

    this._elem.addEventListener('pointermove', this._onPointerMove.bind(this), {
      passive: true,
      signal: this._listenerAbortController.signal,
    });

    this._elem.addEventListener('pointerup', this._onPointerUp.bind(this), {
      signal: this._listenerAbortController.signal,
    });

    this._origin = data;
  }

  _onLeaveMovetest(lifecycle, data) {
    const distance = Math.max(
      Math.abs(this._origin.posAbs.x - data.posAbs.x),
      Math.abs(this._origin.posAbs.y - data.posAbs.y),
    );

    return lifecycle.to === 'hidden' || distance > 5;
  }

  _onEnterMeasure(lifecycle) {
    this._setVisibility(true);

    if (this._infoPopper === null) {
      this._infoPopper = this._createInfoPopper();
      this._infoPopperPlacement = null;
    }
  }

  _onAfterMove(lifecycle, data) {
    const clamp = (value) => Math.max(0, Math.min(100, value));

    const left = clamp(Math.min(this._origin.posPct.x, data.posPct.x));
    const top = clamp(Math.min(this._origin.posPct.y, data.posPct.y));
    const width = clamp(Math.max(this._origin.posPct.x, data.posPct.x)) - left;
    const height = clamp(Math.max(this._origin.posPct.y, data.posPct.y)) - top;

    const frameStyle = this._frame.style;

    frameStyle.left = `${left}%`;
    frameStyle.top = `${top}%`;
    frameStyle.width = `${width}%`;
    frameStyle.height = `${height}%`;

    if (this._infoPopper !== null) {
      const info = this._info;

      info.setAttribute('left-pt', left * this._physWidth.pt / 100);
      info.setAttribute('top-pt', top * this._physHeight.pt / 100);
      info.setAttribute('right-pt', (left + width) * this._physWidth.pt / 100);
      info.setAttribute('bottom-pt', (top + height) * this._physHeight.pt / 100);

      this._updateInfoPopperPlacement(data);
      this._infoPopper.update();
    }
  }

  _reset() {
    if (this._listenerAbortController !== null) {
      this._listenerAbortController.abort();
      this._listenerAbortController = null;
    }

    this._origin = null;
  }

  _onAfterUp(lifecycle) {
    switch (lifecycle.from) {
    case 'movetest':
      const target = this._origin?.target;

      if (target) {
        window.setTimeout(target.click.bind(target), 0);
      }

      break;

    default:
      // Swallow a possible click event.
      window.addEventListener('click', (ev) => {
        ev.stopImmediatePropagation();
        ev.preventDefault();
      }, { capture: true, once: true });
      break;
    }

    this._reset();
  }

  _onAfterCancel() {
    this._reset();
  }

  _onPointerDown(ev) {
    if (ev.isPrimary && ev.buttons == 1 && this._sm.can('down')) {
      ev.preventDefault();
      ev.stopImmediatePropagation();

      this._elem.setPointerCapture(ev.pointerId);

      this._sm.down(this._eventData(ev));
    }
  }

  _onPointerCancel() {
    if (this._sm.can('cancel')) {
      this._sm.cancel();
    }
  }

  _onPointerMove(ev) {
    if (this._sm.can('move')) {
      this._sm.move(this._eventData(ev));
    }
  }

  _onPointerUp(ev) {
    if (this._sm.can('up')) {
      ev.preventDefault();
      ev.stopImmediatePropagation();

      this._sm.up(this._eventData(ev));
    }
  }
}

function initViewer() {
  divViewer.querySelectorAll('img, a').forEach((i) => {
    i.draggable = false;
  });

  (new PositionTool(divViewer)).activate();
  (new MeasurementTool(divViewer)).activate();

  divViewer.addEventListener('click', (ev) => {
    const overlay = ev.target.closest('.dossier_viewer_overlay');

    if (overlay) {
      ev.preventDefault();
      ev.stopPropagation();

      if (overlay.classList.contains('dossier_sketch_node') && overlay.dataset.infoId) {
        const elNodeInfo = document.getElementById(overlay.dataset.infoId);
        if (elNodeInfo) {
          elNodeInfo.scrollIntoView();
        }
      }
    }
  });

  divNodeDialogTemplate.addEventListener('show.bs.modal', (ev) => {
    const modal = event.target;
    const button = event.relatedTarget;
    const overlay = button.closest('.dossier_viewer_overlay');

    elNodeDialogKind.value = overlay.dataset.nodeKind;

    const bounds = JSON.parse(overlay.dataset.nodeBounds);

    Object.entries(bounds).forEach(([key, value]) => {
      elNodeDialogBounds.setAttribute(key, value);
    });

    const text = JSON.parse(overlay.dataset.nodeText);

    elNodeDialogText.value = text;

    let lines = text.split(/(?<=\r\n|\r|\n)/);

    const directEscape = new Map([
      [' ', '\u2423'],
      ['\\', '\\\\'],
      ['\r', '\\r'],
      ['\n', '\\n'],
      ['\t', '\\t'],
      ['\e', '\\e'],
    ]);

    lines = lines.map((line) => line.replace(/(?:\s|\\|[^\u0020-\u007F])/g, c => {
      const replacement = directEscape.get(c, null);

      if (replacement) {
        return replacement;
      }

      return '\\u' + c.charCodeAt(0).toString(16).toUpperCase().padStart(4, '0');
    }));

    elNodeDialogTextDetails.value = lines.join('\n');
  });
}

function initSidebar() {
  const onSketchNodeInfoMouseEnter = (ev) => {
    const id = ev.target.id;
    const modified = new Array();
    const overlay = divViewer.querySelector(`.dossier_sketch_node[data-info-id="${id}"]`);

    if (overlay) {
      overlay.classList.add(kSketchNodeHighlightClass);
      modified.push(overlay.classList);
    }

    divViewer.querySelectorAll(`.dossier_sketch_node_search_area[data-info-id="${id}"]`).forEach((el) => {
      el.classList.add(kSketchNodeSearchAreaHighlightClass);
      modified.push(el.classList);
    });

    if (modified.length > 0) {
      ev.target.addEventListener('mouseleave', (ev) => {
        modified.forEach((classList) => classList.remove(
          kSketchNodeHighlightClass,
          kSketchNodeSearchAreaHighlightClass,
        ));
      }, { passive: true, once: true });
    }
  };

  divSidebar.querySelectorAll(`.${kSketchNodeInfoClass}`).forEach((el) => {
    el.addEventListener('mouseenter', onSketchNodeInfoMouseEnter);
  });
}

initFilter();
initViewer();
initSidebar();

/* vim: set sw=2 sts=2 et : */
