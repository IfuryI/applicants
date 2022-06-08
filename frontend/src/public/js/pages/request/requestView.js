import {globalEventBus} from 'utils/eventbus';
import BaseView from '../baseView';
import {globalRouter} from 'utils/router';
import {PATHS} from 'utils/paths';
import {getFormValues} from 'utils/formDataWork';
import {OK_CODE, BAD_REQUEST, UNAUTHORIZED, INTERNAL_SERVER_ERROR} from 'utils/codes';
import {setListenersForHidingValidationError} from 'utils/setValidationResult';
import {INCORRECT_DATA, INCORRECT_LOGIN, SERVER_ERROR} from 'utils/errorMessages';
import {busEvents} from 'utils/busEvents';
import './request.tmpl';
import {userMeta} from 'utils/userMeta';
import {Navbar} from 'components/navbar';

/**
 * Представление страницы моделирования
 */
export default class RequestView extends BaseView {
    /**
     * Конструктор
     * @param {Element} parent - элемент для рендера
     */
    constructor(parent) {
        // eslint-disable-next-line no-undef
        super(parent, Handlebars.templates['request.hbs']);

        this.formSubmittedCallback = this.formSubmitted.bind(this);
    }

    /**
     * Проверка, если пользователь уже авторизован
     */
    render() {
        this.setLoginPage();
    }

    /**
     * Запуск рендера и установка колбеков
     */
    setLoginPage() {
        super.render();
        this.navbarComponent = new Navbar(document.getElementById('navbar'), {'authorized': true});
        this.navbarComponent.render();
        this.setEventListeners();
    }

    /**
     * "Деструктор" страницы
     */
    hide() {
        super.hide(this);
    }

    /**
     * Установка колбеков
     */
    setEventListeners() {
        document.getElementById('1-submit').onclick = function() {
            globalEventBus.emit(busEvents.REQUEST_START, 1);
        };
        document.getElementById('2-submit').onclick = function() {
            globalEventBus.emit(busEvents.REQUEST_START, 2);
        };
        document.getElementById('3-submit').onclick = function() {
            globalEventBus.emit(busEvents.REQUEST_START, 3);
        };
        document.getElementById('4-submit').onclick = function() {
            globalEventBus.emit(busEvents.REQUEST_START, 4);
        };
        document.getElementById('5-submit').onclick = function() {
            globalEventBus.emit(busEvents.REQUEST_START, 5);
        };
        document.getElementById('6-submit').onclick = function() {
            globalEventBus.emit(busEvents.REQUEST_START, 6);
        };
        document.getElementById('7-submit').onclick = function() {
            globalEventBus.emit(busEvents.REQUEST_START, 7);
        };
        document.getElementById('8-submit').onclick = function() {
            globalEventBus.emit(busEvents.REQUEST_START, 8);
        };
    }

    /**
     * Удаление колбеков
     */
    removeEventListeners() {
        document.getElementById('5-submit')?.removeEventListener('submit', this.formSubmittedCallback);
    }

    /**
     * Обработка отправки формы
     * @param {Object} event - событие отправки формы
     */
    formSubmitted(event) {
    }

    /**
     * Проверка статуса запроса на вход
     * @param {number} status - статус запроса
     */
    processLoginAttempt(status) {
        if (status === OK_CODE) {
            globalRouter.activate(`${PATHS.user}/${userMeta.getUsername()}`);
            return;
        }
        const errors = {
            [BAD_REQUEST]: INCORRECT_DATA,
            [UNAUTHORIZED]: INCORRECT_LOGIN,
            [INTERNAL_SERVER_ERROR]: SERVER_ERROR,
        };
        document.getElementById('validation-hint-login').innerText = errors[status];
    }
}
