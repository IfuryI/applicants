import BaseView from '../baseView';
import {globalEventBus} from 'utils/eventbus';
import {globalRouter} from 'utils/router';
import {PATHS} from 'utils/paths';
import {getFormValues} from 'utils/formDataWork';
import Validator from 'utils/validation';
import {setValidationResult, setListenersForHidingValidationError} from 'utils/setValidationResult';
import {BAD_REQUEST, CREATED, INTERNAL_SERVER_ERROR} from 'utils/codes';
import {ALREADY_EXISTS, INCORRECT_DATA} from 'utils/errorMessages';
import {busEvents} from 'utils/busEvents';
import './signup.tmpl';
import {userMeta} from 'utils/userMeta';
import {Navbar} from 'components/navbar';


/**
 * Представление страницы регистрации
 */
export default class SignupView extends BaseView {
    /**
     * Конструктор
     * @param {Element} parent - элемент для рендера
     */
    constructor(parent) {
        // eslint-disable-next-line no-undef
        super(parent, Handlebars.templates['signup.hbs']);

        globalEventBus.on(busEvents.SIGNUP_STATUS, this.processSignupAttempt.bind(this));

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
        this.setSignupPage();
    }

    /**
     * Запуск рендера и подписка на события
     */
    setSignupPage() {
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
        document.getElementById('signup').addEventListener('submit', this.formSubmittedCallback);
        this.navbarComponent.setEventListeners();
        setListenersForHidingValidationError();
    }

    /**
     * Удаление колбеков
     */
    removeEventListeners() {
        document.getElementById('signup')?.removeEventListener('submit', this.formSubmittedCallback);
        this.navbarComponent.removeEventListeners();
    }

    /**
     * Обработка отправки формы регистрации
     * @param {Object} event - событие отправки формы
     */
    formSubmitted(event) {
        event.preventDefault();

        const data = getFormValues(event.target);

        const validator = new Validator();
        const loginError = validator.validateLogin(data.username);
        const emailError = validator.validateEmail(data.email);
        const passwordError = validator.validatePassword(data.password);

        [
            [
                document.getElementById('username-input'),
                document.getElementById('validation-hint-login'),
                loginError,
            ],
            [
                document.getElementById('email-input'),
                document.getElementById('validation-hint-email'),
                emailError,
            ],
            [
                document.getElementById('password-input'),
                document.getElementById('validation-hint-password'),
                passwordError,
            ],
        ].forEach(([inputField, inputHint, errors]) => setValidationResult(inputField, inputHint, errors));

        if ([loginError, emailError, passwordError].every((error) => error.length === 0)) {
            globalEventBus.emit(busEvents.SIGNUP_CLICKED, data);
        }
    }

    /**
     * Обработка статуса после запроса регистрации
     * @param {number} status - статус запроса
     */
    processSignupAttempt(status) {
        if (status === CREATED) {
            globalRouter.activate(`${PATHS.user}/${userMeta.getUsername()}`);
            return;
        }
        const errors = {
            [BAD_REQUEST]: INCORRECT_DATA,
            [INTERNAL_SERVER_ERROR]: ALREADY_EXISTS,
        };
        document.getElementById('validation-hint-signup').innerText = errors[status];
    }
}
