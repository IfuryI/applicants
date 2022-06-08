'use strict';


/**
 * Класс ивентбаса
 */
export class EventBus {
    #callbacks;

    /**
     * Конструктор
     */
    constructor() {
        this.#callbacks = {};
    }

    /**
     * Подписка на событие
     * @param {string} name - имя события
     * @param {function} callback - колбек функция
     */
    on(name, callback) {
        if (!(name in this.#callbacks)) {
            this.#callbacks[name] = [];
        }

        this.#callbacks[name].push(callback);
    }

    /**
     * Вызов события
     * @param {string} name - имя события
     * @param {Object} args - аргументы функции колбека
     */
    emit(name, ...args) {
        if (name in this.#callbacks) {
            this.#callbacks[name].forEach((callback) => {
                try {
                    callback(...args);
                } catch (error) {
                    console.error(error);
                }
            });
        }
    }
}

export const globalEventBus = new EventBus();
