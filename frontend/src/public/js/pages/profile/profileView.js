import {globalEventBus} from 'utils/eventbus';
import BaseView from '../baseView';
import {busEvents} from 'utils/busEvents';
import {Navbar} from 'components/navbar';
import {userMeta} from 'utils/userMeta';
import './profile.tmpl';

/**
 * Представление страницы профиля
 */
export default class ProfileView extends BaseView {
    /**
     * Конструктор
     * @param {Element} parent - элемент для рендера
     */
    constructor(parent) {
        // eslint-disable-next-line no-undef
        super(parent, Handlebars.templates['profile.hbs']);

        globalEventBus.on(busEvents.SET_PROFILE_DATA, this.setProfileData.bind(this));
    }

    /**
     * Запуск рендера
     * @param {string} username - имя пользователя
     */
    render(username) {
        globalEventBus.emit(busEvents.GET_PROFILE_DATA, username);
    }

    /**
     * Установка данных профиля
     * @param {Object} data - данные профиля
     */
    setProfileData(data) {
        this.username = data.username;
        console.log(data)
        super.render(data);

        this.navbarComponent = new Navbar(document.getElementById('navbar'),
            {'authorized': userMeta.getAuthorized()});
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
        document.getElementById('follow-button')?.addEventListener('click', this.followClickedCallback);
        this.navbarComponent.setEventListeners();
    }

    /**
     * Удаление колбеков
     */
    removeEventListeners() {
        document.getElementById('follow-button')?.removeEventListener('click', this.followClickedCallback);
        this.navbarComponent.removeEventListeners();
    }
}
