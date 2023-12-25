import * as geometry from './geometry.js';
import * as defaultunit from './defaultunit.js';

function removeAllChildNodes(node) {
  while (node.lastChild) {
    node.removeChild(node.lastChild);
  }
}

function enableTooltip(owner, container) {
  bootstrap.Tooltip.getOrCreateInstance(owner, {
    selector: '[data-bs-toggle="tooltip"]',
    container: container,
  });
}

function cleanupTooltip(owner) {
  const tooltip = bootstrap.Tooltip.getInstance(owner);

  if (tooltip) {
    tooltip.dispose();
  }
}

function toTitleCase(str) {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

function createNowrapSpan(title) {
  const span = document.createElement('span');

  span.classList.add('text-nowrap');

  if (title) {
    span.title = title;
    span.dataset.bsToggle = "tooltip";
  }

  return span;
}

const getUnitNumberFormat = (() => {
  const cache = new Map();

  return (unit, signDisplay) => {
    const key = JSON.stringify([unit.name, signDisplay]);

    let numFormat = cache.get(key);

    if (!numFormat) {
      let fractionDigits = 0;

      switch (unit) {
      case geometry.Inch:
      case geometry.Centimeter:
        fractionDigits = 1;
        break;
      }

      numFormat = new Intl.NumberFormat(undefined, {
        style: 'decimal',
        minimumFractionDigits: fractionDigits,
        maximumFractionDigits: fractionDigits,
        signDisplay: signDisplay,
      });

      cache.set(key, numFormat);
    }

    return numFormat;
  };
})();

class LengthFormatter {
  _target = null;
  _points = new geometry.Length(0);
  _showSign = false;
  _unit = geometry.Points;
  _text = null;

  get target() {
    return this._target;
  }

  set target(target) {
    this._target = target;
  }

  get points() {
    return this._points.pt;
  }

  set points(pt) {
    if (pt !== this._points.pt) {
      this._points = new geometry.Length(pt, geometry.Point);
      this._text = null;
    }
  }

  get showSign() {
    return this._showSign;
  }

  set showSign(showSign) {
    this._showSign = showSign;
    this._text = null;
  }

  get unit() {
    return this._unit;
  }

  set unit(unit) {
    if (unit !== this._unit) {
      this._unit = unit;
      this._text = null;
    }
  }

  render() {
    if (this._target !== null) {
      if (this._text === null) {
        const numFormat = getUnitNumberFormat(this._unit, this._showSign ? 'exceptZero' : 'auto')

        this._text = numFormat.format(this._points.toUnit(this._unit));
      }

      this._target.innerText = this._text;
    }
  }
}

class GeometryPoint extends HTMLElement {
  static observedAttributes = ["left-pt", "top-pt"];

  _stopDefaultUnitObserver = null;
  _elUnit = null;

  constructor() {
    super();
    this._formatters = {
      left: new LengthFormatter(),
      top: new LengthFormatter(),
    };
  }

  connectedCallback() {
    removeAllChildNodes(this);

    this._stopDefaultUnitObserver =
      defaultunit.observe(this._onDefaultUnitChange.bind(this));

    const el = {
      left: createNowrapSpan('Left'),
      top: createNowrapSpan('Top'),
    };

    this._elUnit = createNowrapSpan();

    this._setUnit(defaultunit.get());
    this._formatters.left.target = el.left;
    this._formatters.top.target = el.top;

    const elWrapper = document.createElement('span');

    elWrapper.append('(', el.left, ', ', el.top, ') ', this._elUnit);

    this.append(elWrapper);

    enableTooltip(this, elWrapper);

    this._render();
  }

  disconnectedCallback() {
    cleanupTooltip(this);

    if (this._stopDefaultUnitObserver !== null) {
      this._stopDefaultUnitObserver();
    }

    this._elUnit = null;

    Object.values(this._formatters).forEach((formatter) => {
      formatter.target = null;
    });
  }

  attributeChangedCallback(name, oldValue, newValue) {
    const target = {
      'left-pt': this._formatters.left,
      'top-pt': this._formatters.top,
    }[name];

    if (target) {
      target.points = parseFloat(newValue);
      this._render();
    }
  }

  _setUnit(unit) {
    Object.values(this._formatters).forEach((formatter) => {
      formatter.unit = unit;
    });
  }

  _onDefaultUnitChange(unit) {
    this._setUnit(unit);
    this._render();
  }

  _render() {
    Object.values(this._formatters).forEach((formatter) => {
      formatter.render();
    });

    if (this._elUnit !== null) {
      this._elUnit.innerText = this._formatters.left.unit.name;
    }
  }
}

class _GeometryRectFormatters {
  constructor() {
    this._left = new LengthFormatter();
    this._top = new LengthFormatter();

    this._right = new LengthFormatter();
    this._bottom = new LengthFormatter();

    this._width = new LengthFormatter();
    this._height = new LengthFormatter();

    this._width.showSign = true;
    this._height.showSign = true;
  }

  reset() {
    Object.values(this).forEach(formatter => {
      formatter.points = 0;
      formatter.target = null;
    });
  }

  get unit() {
    return this._left.unit;
  }

  set unit(unit) {
    Object.values(this).forEach(formatter => {
      formatter.unit = unit;
    });
  }

  setEdgeTarget(edge, target) {
    this[`_${edge}`].target = target;
  }

  setEdgePoints(edge, points) {
    this[`_${edge}`].points = points;

    switch (edge) {
    case 'left':
    case 'right':
      this._width.points = this._right.points - this._left.points;
    case 'top':
    case 'bottom':
      this._height.points = this._bottom.points - this._top.points;
    }
  }

  render() {
    Object.values(this).forEach(formatter => formatter.render());
  }
}

class GeometryRect extends HTMLElement {
  static observedAttributes = ["left-pt", "top-pt", "right-pt", "bottom-pt"];

  _stopDefaultUnitObserver = null;
  _elUnit = null;

  constructor() {
    super();
    this._formatters = new _GeometryRectFormatters();
  }

  connectedCallback() {
    removeAllChildNodes(this);

    this._stopDefaultUnitObserver =
      defaultunit.observe(this._onDefaultUnitChange.bind(this));

    const elSpan = new Object;

    ['left', 'top', 'right', 'bottom', 'width', 'height'].forEach(edge => {
      const span = createNowrapSpan(toTitleCase(edge));

      this._formatters.setEdgeTarget(edge, span);

      elSpan[edge] = span;
    });

    this._formatters.unit = defaultunit.get();

    const
      elWrapper = document.createElement('span'),
      elUnit = createNowrapSpan();

    this._elUnit = elUnit;

    elWrapper.append(
      '(', elSpan.left, ', ', elSpan.top, ')\u2013(',
      elSpan.right, ', ', elSpan.bottom, ') [',
      elSpan.width, '\u00D7', elSpan.height, '] ',
      elUnit,
    );

    this.append(elWrapper);

    enableTooltip(this, elWrapper);

    this._render();
  }

  disconnectedCallback() {
    cleanupTooltip(this);

    if (this._stopDefaultUnitObserver !== null) {
      this._stopDefaultUnitObserver();
    }

    this._elUnit = null;
    this._formatters.reset();
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (name.endsWith('-pt')) {
      const edge = name.replace(/-pt$/, '');

      this._formatters.setEdgePoints(edge, parseFloat(newValue));
      this._render();
    }
  }

  _onDefaultUnitChange(unit) {
    this._formatters.unit = unit;
    this._render();
  }

  _render() {
    this._formatters.render();

    if (this._elUnit !== null) {
      this._elUnit.innerText = this._formatters.unit.name;
    }
  }
}

customElements.define("dossier-geometry-point", GeometryPoint);
customElements.define("dossier-geometry-rect", GeometryRect);

/* vim: set sw=2 sts=2 et : */
