'use strict'

const path = require('path')

module.exports = [
  {
    mode: 'none',
    entry: {
      import: 'javascript-state-machine',
    },
    output: {
      filename: 'state-machine.js',
      path: path.resolve(__dirname, 'static'),
      library: {
        type: 'umd',
        name: 'StateMachine',
      },
    },
  },
];

/* vim: set sw=2 sts=2 et : */
