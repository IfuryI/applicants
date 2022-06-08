import {setValidationResult, setListenersForHidingValidationError} from 'utils/setValidationResult';
import {globalEventBus} from 'utils/eventbus';
import BaseView from '../baseView';
import Validator from 'utils/validation';
import {globalRouter} from 'utils/router';
import {PATHS} from 'utils/paths';
import {OK_CODE, PAYLOAD_TOO_LARGE} from 'utils/codes';
import {PASSWORDS_MISMATCH} from 'utils/errorMessages';
import {busEvents} from 'utils/busEvents';
import './settings.tmpl';
import {Navbar} from 'components/navbar';
import {userMeta} from 'utils/userMeta';


/**
 * Представление страницы настроек профиля
 */
export default class SettingsView extends BaseView {
    /**
     * Конструктор
     * @param {Element} parent - элемент для рендера
     */
    constructor(parent) {
        // eslint-disable-next-line no-undef
        super(parent, Handlebars.templates['settings.hbs']);

        this.settings = {};
        this.input = {};

        globalEventBus.on(busEvents.SET_SETTINGS_DATA, this.setSettings.bind(this));
        globalEventBus.on(busEvents.RESPONSE_CHANGE_SETTINGS, this.displayServerResponse.bind(this));
        globalEventBus.on(busEvents.AVATAR_UPLOADED, this.displayServerResponseAvatar.bind(this));

        this.saveClickedCallback = this.saveClicked.bind(this);
        this.uploadAvatarClickedCallback = this.uploadAvatarClicked.bind(this);
    }

    /**
     * Запуск рендера
     */
    render() {
        globalEventBus.emit(busEvents.GET_SETTINGS_DATA);
    }

    /**
     * Изменение настроек
     * @param {Object} data - новые данные
     */
    setSettings(data) {
        this.settings = data;
        super.render(this.settings);

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
        document.getElementById('settings-save-button').addEventListener('click', this.saveClickedCallback);
        document.getElementById('avatar-upload-button').addEventListener('click', this.uploadAvatarClickedCallback);

        this.navbarComponent.setEventListeners();

        this.setPhotoPreviewListener();

        setListenersForHidingValidationError();
    }

    /**
     * Показ превью фото
     */
    setPhotoPreviewListener() {
        const file = document.getElementById('avatar');
        const preview = document.getElementById('avatar-img');

        file.addEventListener('change', (event) => {
            if (file.files.length === 0) {
                return;
            }

            const avatar = file.files[0];
            if (!this.validateAvatar(avatar)) {
                return;
            }

            if (avatar.type.indexOf('image') === -1) {
                return;
            }

            const fileReader = new FileReader();
            fileReader.onload = (event) => {
                if (getComputedStyle(preview, null).display === 'none') {
                    preview.style.display = 'block';
                }

                preview.src = event.target.result;
            };

            fileReader.readAsDataURL(avatar);
        });
    }

    /**
     * Удаление колбеков
     */
    removeEventListeners() {
        document.getElementById('settings-save-button').removeEventListener('click', this.saveClickedCallback);
        document.getElementById('avatar-upload-button').removeEventListener('click', this.uploadAvatarClickedCallback);
        this.navbarComponent.removeEventListeners();
    }

    /**
     * Обработчик нажатия на кнопку загрузки аватара
     */
    uploadAvatarClicked() {
        const selectedFile = document.getElementById('avatar').files[0];
        globalEventBus.emit(busEvents.UPLOAD_AVATAR, selectedFile);
    }

    /**
     * Обработчик нажатия на кнопку сохранения
     */
    saveClicked() {
        this.inputSettings();
        if (this.validateSettings()) {
            this.sendInput(this.deltaSettings());
        } else {
            document.getElementById('settings-server-response').innerHTML = '';
        }
        this.input = {};
    }

    /**
     * Проверка изменений данных о пользователе
     * @return {Object} - обновленные настройки
     */
    deltaSettings() {
        const settings = {};
        if (this.settings.fullname !== this.input.fullname) {
            settings.fullname = this.input.fullname;
        }
        if (this.settings.email !== this.input.email) {
            settings.email = this.input.email;
        }
        if (this.input.password1.length !== 0) {
            settings.password = this.input.password1;
        }
        return settings;
    }

    /**
     * Получение новых данных
     */
    inputSettings() {
        const email = document.getElementById('user-email').value;
        const password1 = document.getElementById('user-password').value;
        const password2 = document.getElementById('user-password-repeat').value;
        this.input = {
            email,
            password1,
            password2,
        };
    }

    /**
     * Проверка данных на валидность
     * @return {boolean} - статус наличия ошибок
     */
    validateSettings() {
        const validator = new Validator();

        const passwordErrors = validator.validatePassword(this.input.password1);
        const emailErrors = validator.validateEmail(this.input.email);

        if (this.input.password1 !== this.input.password2) {
            passwordErrors.push(PASSWORDS_MISMATCH);
        }
        if (this.input.password1.length === 0 && this.input.password2.length === 0) {
            passwordErrors.length = 0;
        }

        [
            [
                document.getElementById('user-email'),
                document.getElementById('settings-errors-email'),
                emailErrors,
            ],
            [
                document.getElementById('user-password'),
                document.getElementById('settings-errors-password'),
                passwordErrors,
            ],
        ].forEach(([inputField, inputHint, errors]) => setValidationResult(inputField, inputHint, errors));

        return ([...passwordErrors, ...emailErrors].length === 0);
    }

    /**
     * Проверка на валидность аватара
     * @param {string} avatar - путь к фото
     * @return {boolean} - статус наличия ошибок
     */
    validateAvatar(avatar) {
        const validator = new Validator();

        const avatarErrors = validator.validateAvatar(avatar);

        document.getElementById('settings-avatar-errors').innerHTML = avatarErrors.join('<br>');

        [
            [
                document.getElementById('avatar'),
                document.getElementById('settings-avatar-errors'),
                avatarErrors,
            ],
        ].forEach(([inputField, inputHint, errors]) => setValidationResult(inputField, inputHint, errors));

        return (avatarErrors.length === 0);
    }

    /**
     * Отправка запроса на изменение настроек
     * @param {Object} settings - измененные данные
     */
    sendInput(settings) {
        if (JSON.stringify(settings) !== '{}') {
            globalEventBus.emit(busEvents.REQUEST_CHANGE_SETTINGS, settings);
        }
    }

    /**
     * Показать изменения
     * @param {number} status - статус запроса
     */
    displayServerResponse(status) {
        if (status === OK_CODE) {
            globalRouter.activate(PATHS.user + '/' + userMeta.getUsername());
        } else {
            document.getElementById('settings-server-response').innerHTML = 'Ошибка сохранения настроек!';
        }
    }

    /**
     * Показать изменения аватара
     * @param {number} status - статус запроса
     */
    displayServerResponseAvatar(status) {
        let errorMessage = '';
        switch (status) {
        case OK_CODE:
            this.render();
            return;
        case PAYLOAD_TOO_LARGE:
            errorMessage = 'Слишком большой размер картинки!';
            break;
        default:
            errorMessage = 'Ошибка загрузки аватарки!';
            break;
        }
        document.getElementById('settings-avatar-errors').innerHTML = errorMessage;
    }
}
