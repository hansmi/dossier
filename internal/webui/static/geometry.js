class LengthUnit {
  constructor(name, pointsPerUnit) {
    this._name = name;
    this._pointsPerUnit = pointsPerUnit;
  }

  static fromName(name) {
    const unit = _unitByName[name];

    if (!unit) {
      throw new Error(`Unknown unit ${name}`);
    }

    return unit;
  }

  get name() {
    return this._name;
  }

  get pointsPerUnit() {
    return this._pointsPerUnit;
  }
}

function newUnit(name, pointsPerUnit) {
  return Object.freeze(new LengthUnit(name, pointsPerUnit));
}

const
  Point = newUnit('pt', 1),
  Inch = newUnit('in', 72 * Point.pointsPerUnit),
  Centimeter = newUnit('cm', Inch.pointsPerUnit / 2.54),
  Millimeter = newUnit('mm', Inch.pointsPerUnit / 25.4);

const ALL_UNITS = Object.freeze([
  Point,
  Inch,
  Centimeter,
  Millimeter,
]);

const _unitByName = Object.fromEntries(ALL_UNITS.map((unit) => [unit.name, unit]));

class Length {
  constructor(value, unit) {
    if (unit === undefined || unit === null) {
      unit = Point;
    }

    this._value = value * unit.pointsPerUnit;
  }

  get pt() {
    return this._value;
  }

  toUnit(unit) {
    return this._value / unit.pointsPerUnit;
  }
}

export {
  ALL_UNITS,
  Point,
  Inch,
  Centimeter,
  Millimeter,
  LengthUnit,
  Length,
};

/* vim: set sw=2 sts=2 et : */
