'use strict';

import { EventBus } from 'utils/eventbus.js'

QUnit.module('EventBus test', function () {

    QUnit.test('Add single callback and call it', function (assert) {
        const bus = new EventBus();
        let a = 0;
        bus.on('sample', () => { a = a + 1; });
        bus.emit('sample');

        assert.strictEqual(a, 1);
	});

	QUnit.test('Add multiple callbacks and call them', function (assert) {
        const bus = new EventBus();
        let a = 0;
        let b = 0;
        bus.on('sample', () => { a = a + 1; });
        bus.on('sample', () => { b = b + 1337; });
        bus.emit('sample');

        assert.strictEqual(a, 1);
        assert.strictEqual(b, 1337);
	});

    QUnit.test('Callbacks are independent', function (assert) {
        const bus = new EventBus();
        let a = 0;
        bus.on('sample', () => { a = 1; });
        bus.on('sample', () => { a = 2; throw new Error('Test Error: this is okay'); });
        bus.on('sample', () => { a = 3; });
        bus.emit('sample');

        assert.strictEqual(a, 3);
    });

    QUnit.test('Arguments', function (assert) {
        const bus = new EventBus();
        let a = 0;
        const b = [];
        bus.on('sample', (value) => { a = value; });
        bus.on('sample', (value) => { b.push(value); });
        bus.emit('sample', 1337);

        assert.strictEqual(a, 1337);
        assert.deepEqual(b, [ 1337 ]);

        bus.emit('sample', 228);

        assert.strictEqual(a, 228);
        assert.deepEqual(b, [ 1337, 228 ]);
    });

});

