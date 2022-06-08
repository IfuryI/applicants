import {globalEventBus} from 'utils/eventbus';
import {API} from 'utils/api';
import {busEvents} from 'utils/busEvents';
import {OK_CODE} from 'utils/codes';
import {userMeta} from 'utils/userMeta';


/**
 *  Модель страницы логина
 */
export default class RequestModel {
    /**
     * Конструктор
     */
    constructor() {
        globalEventBus.on(busEvents.REQUEST_START, this.customRequest.bind(this));
    }

    /**
     * Проверка успешности авторизации пользователя
     * @param {Object} userData - данные пользователя
     */
     customRequest(btnType) {
        API.customRequest(btnType)
            .then((res) => {

            });
    }
}
