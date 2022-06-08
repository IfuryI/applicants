import {globalEventBus} from 'utils/eventbus';
import {API} from 'utils/api';
import {busEvents} from 'utils/busEvents';
import {OK_CODE} from 'utils/codes';
import {userMeta} from 'utils/userMeta';


/**
 *  Модель страницы логина
 */
export default class LoginModel {
    /**
     * Конструктор
     */
    constructor() {
        globalEventBus.on(busEvents.LOGIN_CLICKED, this.checkLogin.bind(this));
    }

    /**
     * Проверка успешности авторизации пользователя
     * @param {Object} userData - данные пользователя
     */
    checkLogin(userData) {
        API.login(userData)
            .then((res) => {
                userMeta.setAuthorized(res.status === OK_CODE);
                userMeta.setUsername(userData.username);
                globalEventBus.emit(busEvents.LOGIN_STATUS, res.status);
            });
    }
}
