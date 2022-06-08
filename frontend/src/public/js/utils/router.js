'use strict';

import {PATHS} from './paths';
import {findAscendingTag} from './findAscendingTag';

/**
 * Роутер
 */
export class Router {
    /**
     * Конструктор - создание роутов
     */
    constructor() {
        this.routes = new Map();
        this._currentPath = '';
        this._currentParameters = '';
    }

    /**
     * Добавление нового роута
     * @param {string} path - путь
     * @param {Object} view - объект представления с методом render()
     */
    register(path, view) {
        this.routes.set(path, view);
    }

    /**
     * Запуск роутера
     */
    start() {
        window.addEventListener('load', () => {
            this.activate(location.pathname, location.search);
        });

        window.addEventListener('popstate', () => {
            this.activate(location.pathname, location.search);
        });

        window.addEventListener('click', (event) => {
            const {target} = event;
            const link = findAscendingTag(target, 'A');
            const href = (link !== null) ? link.href : null;
            this.activate(href, '', event);
        });
    }

    /**
     * Поиск пути и выделение параметров из url
     * @param {string|null} path - путь
     * @param {string} query - query-параметры
     * @param {Object} event - событие, порожденное кликом
     */
    activate(path, query='', event=null) {
        if (path === null) {
            return;
        }

        for (const i in PATHS) {
            if (path.includes(PATHS[i])) {
                event?.preventDefault();
                const params = path.substring(path.indexOf(PATHS[i]) + PATHS[i].length + 1);
                this.pushState(PATHS[i], query, {}, params);
                return;
            }
        }
    }

    /**
     * Формирование пути в адресной строке и history.pushState()
     * @param {string} path - путь
     * @param {string} query - query-параметры
     * @param {Object} state - объект состояния
     * @param {string} parameters - параметры пути
     */
    pushState(path = '/', query, state = {}, parameters = '') {
        let newPath = path;
        if (query) {
            parameters += (location.pathname.slice(-1) === '/') ? query : '/' + query;
        }
        if (parameters) {
            newPath = `${path}/${parameters}`;
        }
        if (newPath !== location.pathname) {
            history.pushState(state, document.title, newPath);
            // если переходим на новую страницу, прокрутим страницу наверх
            this.scrollToTop();
        }

        this.handlePath(path, parameters);
    }

    /**
     * Запуск рендера представления по пути
     * @param {string} path - путь
     * @param {string} parameters - параметры пути
     */
    handlePath(path, parameters = '') {
        this.currentView?.hide();
        this.currentView = this.routes.get(path);
        this.currentView.render(parameters);
        this.scrollToTop();
        this._currentPath = path;
        this._currentParameters = parameters;
    }

    /**
     * Получить текущий путь
     * @return {string} - текуший путь
     */
    get currentPath() {
        return this._currentPath;
    }

    /**
     * Получить текущие параметры страницы
     * @return {string} - текущие параметры
     */
    get currentParameters() {
        return this._currentParameters;
    }

    /**
     * После перехода на другой путь, прокрутим положение в документе до начала страницы
     */
    scrollToTop() {
        window.scrollTo(0, 0);
    }

    /**
     * Переход на предыдущую страницу
     */
    back() {
        window.history.back();
    }

    /**
     * Переход на следующую страницу
     */
    forward() {
        window.history.forward();
    }
}

export const globalRouter = new Router();
