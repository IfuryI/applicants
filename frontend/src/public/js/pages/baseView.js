'use strict';


/**
 * Базовый класс представления
 */
export default class BaseView {
    /**
     * Конструктор
     * @param {Element} parent - элемент для рендера
     * @param {function} template - функция шаблона
     */
    constructor(parent, template) {
        this.parent = parent;
        this.template = template;
    }

    /**
     * Рендер шаблона с данными
     * @param {Object} data - данные для страницы
     */
    render(data) {
        this.parent.innerHTML = this.template(data);
    }

    /**
     * "Деструктор" страницы
     * @param {Object} page - объект страницы, у которой определен removeEventListeners()
     */
    hide(page) {
        page.removeEventListeners();
    }
}
