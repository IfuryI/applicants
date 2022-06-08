'use strict';

import Validator from 'utils/validation.js'

QUnit.module('Validation test', function () {

    QUnit.test('email validation', function(assert) {
        const validator = new Validator();

        assert.deepEqual(validator.validateEmail('@'), []);
        assert.deepEqual(validator.validateEmail('example'), [
            'incorrect email'
        ]);
        assert.deepEqual(validator.validateEmail('example@example.com'), []);
    });

    QUnit.test('login validation', function (assert) {
        const validator = new Validator();

        assert.deepEqual(validator.validateLogin('grillow!'), [
            'should contain only latin letters, numbers and underscores'
        ]);
        assert.deepEqual(validator.validateLogin('grillow rules!'), [
            'should contain only latin letters, numbers and underscores'
        ]);
        assert.deepEqual(validator.validateLogin('admin'), [
            'should be from 6 to 20 characters long'
        ]);
        assert.deepEqual(validator.validateLogin('XxX_grillow1337_XxX'), []);
    });

    QUnit.test('password validation', function (assert) {
        const validator = new Validator();

        assert.deepEqual(validator.validatePassword('qwerty'), [
            'should be from 8 to 32 characters long',
            'should contain at least one digit'
        ]);
        assert.deepEqual(validator.validatePassword('qwerty123'), []);
        assert.deepEqual(validator.validatePassword('Qwerty123'), []);
        assert.deepEqual(validator.validatePassword('Qwerty123_'), []);
        assert.deepEqual(validator.validatePassword('Aa@4567'), [
            'should be from 8 to 32 characters long'
        ]);
        assert.deepEqual(validator.validatePassword('12345678'), [
            'should contain at least one letter'
        ]);
        assert.deepEqual(validator.validatePassword('Aa@45678'), []);
    });

});

