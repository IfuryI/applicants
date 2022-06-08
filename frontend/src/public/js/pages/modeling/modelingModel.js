import {globalEventBus} from 'utils/eventbus';
import {API} from 'utils/api';
import {busEvents} from 'utils/busEvents';
import {OK_CODE} from 'utils/codes';
import {userMeta} from 'utils/userMeta';


/**
 *  Модель страницы логина
 */
export default class ModelingModel {
    /**
     * Конструктор
     */
    constructor() {
        globalEventBus.on(busEvents.MODEL_START, this.startModeling.bind(this));
    }

    /**
     * Запуск моделирования
     * @param {Object} userData - введеный конфиг
     */
     startModeling(userData) {
        API.startModeling(userData)
            .then((res) => {
                userMeta.setAuthorized(res.status === OK_CODE);
                userMeta.setUsername(userData.username);
                globalEventBus.emit(busEvents.MODEL_GOOD, res.status);
            });
    }
}
