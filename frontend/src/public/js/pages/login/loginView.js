import {globalEventBus} from 'utils/eventbus';
import BaseView from '../baseView';
import {globalRouter} from 'utils/router';
import {PATHS} from 'utils/paths';
import {getFormValues} from 'utils/formDataWork';
import {OK_CODE, BAD_REQUEST, UNAUTHORIZED, INTERNAL_SERVER_ERROR} from 'utils/codes';
import {setListenersForHidingValidationError} from 'utils/setValidationResult';
import {INCORRECT_DATA, INCORRECT_LOGIN, SERVER_ERROR} from 'utils/errorMessages';
import {busEvents} from 'utils/busEvents';
import './login.tmpl';
import {userMeta} from 'utils/userMeta';
import {Navbar} from 'components/navbar';

/**
 * Представление страницы логина
 */
export default class LoginView extends BaseView {
    /**
     * Конструктор
     * @param {Element} parent - элемент для рендера
     */
    constructor(parent) {
        // eslint-disable-next-line no-undef
        super(parent, Handlebars.templates['login.hbs']);

        globalEventBus.on(busEvents.LOGIN_STATUS, this.processLoginAttempt.bind(this));

        this.formSubmittedCallback = this.formSubmitted.bind(this);
    }

    /**
     * Проверка, если пользователь уже авторизован
     */
    render() {
        if (userMeta.getAuthorized()) {
            globalRouter.activate(`${PATHS.user}/${userMeta.getUsername()}`);
            return;
        }
        this.setLoginPage();
    }

    /**
     * Запуск рендера и установка колбеков
     */
    setLoginPage() {
        super.render();

        this.navbarComponent = new Navbar(document.getElementById('navbar'), {'authorized': false});
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
        document.getElementById('login').addEventListener('submit', this.formSubmittedCallback);
        this.navbarComponent.setEventListeners();
        setListenersForHidingValidationError();
    }

    /**
     * Удаление колбеков
     */
    removeEventListeners() {
        document.getElementById('login')?.removeEventListener('submit', this.formSubmittedCallback);
        this.navbarComponent.removeEventListeners();
    }

    /**
     * Обработка отправки формы
     * @param {Object} event - событие отправки формы
     */
    formSubmitted(event) {
        event.preventDefault();
        globalEventBus.emit(busEvents.LOGIN_CLICKED, getFormValues(event.target));
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
