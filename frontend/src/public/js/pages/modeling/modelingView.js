import {globalEventBus} from 'utils/eventbus';
import BaseView from '../baseView';
import {globalRouter} from 'utils/router';
import {PATHS} from 'utils/paths';
import {getFormValues, getModelingFormValues} from 'utils/formDataWork';
import {OK_CODE, BAD_REQUEST, UNAUTHORIZED, INTERNAL_SERVER_ERROR} from 'utils/codes';
import {setListenersForHidingValidationError} from 'utils/setValidationResult';
import {INCORRECT_DATA, INCORRECT_LOGIN, SERVER_ERROR} from 'utils/errorMessages';
import {busEvents} from 'utils/busEvents';
import './modeling.tmpl';
import {userMeta} from 'utils/userMeta';
import {Navbar} from 'components/navbar';

/**
 * Представление страницы моделирования
 */
export default class ModelingView extends BaseView {
    /**
     * Конструктор
     * @param {Element} parent - элемент для рендера
     */
    constructor(parent) {
        // eslint-disable-next-line no-undef
        super(parent, Handlebars.templates['modeling.hbs']);

        globalEventBus.on(busEvents.MODEL_GOOD, this.processModeling.bind(this));

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
        document.getElementById('modeling').addEventListener('submit', this.formSubmittedCallback);
    }

    /**
     * Удаление колбеков
     */
    removeEventListeners() {
        document.getElementById('login')?.removeEventListener('submit', this.formSubmittedCallback);
    }

    /**
     * Обработка отправки формы
     * @param {Object} event - событие отправки формы
     */
    formSubmitted(event) {
        event.preventDefault();
        globalEventBus.emit(busEvents.MODEL_START, getModelingFormValues(event.target));
    }

    /**
     * 
     * @param {number} status - статус запроса
     */
     processModeling(status) {
        if (status === OK_CODE) {
            document.getElementById('validation-hint-login').innerText = "Моделирование запущено успешно!";
        }

        document.getElementById('validation-hint-login').innerText = "Ошибка запуска моделирования";
    }
}
