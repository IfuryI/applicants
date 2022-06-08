import {globalEventBus} from 'utils/eventbus';
import {API} from 'utils/api';
import {busEvents} from 'utils/busEvents';
import {userMeta} from 'utils/userMeta';
import {CREATED} from 'utils/codes';

/**
 *  Модель страницы регистрации
 */
export default class Model {
    /**
     * Конструктор
     */
    constructor() {
        globalEventBus.on(busEvents.SIGNUP_CLICKED, this.createUser.bind(this));
    }

    /**
     * Создание (регистрация) пользователя
     * @param {Object} data - данные о пользователе
     */
    createUser(data) {
        API.signup(data)
            .then((res) => {
                userMeta.setAuthorized(res.status === CREATED);
                userMeta.setUsername(data.username);
                globalEventBus.emit(busEvents.SIGNUP_STATUS, res.status);
            });
    }
}
