import {API} from './api';
import {OK_CODE} from './codes';

/**
 * Класс, предоставляющий служебную информацию о текущем пользователе
 */
class UserMeta {
    /**
     * Конструктор, изначально ставим статус "Не авторизован"
     */
    constructor() {
        this.authStatus = false;
        this.username = null;
    }

    /**
     * Выставить в localStorage значение авторизованности
     * @return {Promise<void>} - промис с действием "выставление статуса авторизации"
     */
    async initMeta() {
        API.getCurrentUser().then((res) => {
            this.setAuthorized(res.status === OK_CODE);
            this.setUsername(res.data.username);
        });
    }

    /**
     * Получаем статус авторизованности
     * @return {boolean} - авторизован ли пользователь
     */
    getAuthorized() {
        return this.authStatus;
    }

    /**
     * Устанавливаем статус авторизованности
     * @param {boolean} authStatus - авторизован ли пользователь
     */
    setAuthorized(authStatus) {
        this.authStatus = authStatus;
    }

    /**
     * Получаем идентификатор текущего пользователя
     * @return {string} - имя пользователя
     */
    getUsername() {
        return this.username;
    }

    /**
     * Устанавливаем идентификатор текущего пользователя
     * @param {string|null} username - имя пользователя
     */
    setUsername(username) {
        this.username = username;
    }
}

export const userMeta = new UserMeta();
