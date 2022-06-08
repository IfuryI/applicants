import {globalEventBus} from 'utils/eventbus';
import {API} from 'utils/api';
import {globalRouter} from 'utils/router';
import {PATHS} from 'utils/paths';
import {OK_CODE} from 'utils/codes';
import {AUTH_ERROR} from 'utils/errorMessages';
import {busEvents} from 'utils/busEvents';

/**
 *  Модель страницы профиля
 */
export default class ProfileModel {
    /**
     * Конструктор
     */
    constructor() {
        globalEventBus.on(busEvents.GET_PROFILE_DATA, this.getProfileData.bind(this));
    }

    /**
     * Получение данных профиля
     * @param {string} username - имя пользователя
     */
    getProfileData(username) {
        Promise.all([API.getUser(username)])
            .then((responses) => {
                if (responses.some((resp) => resp.status !== OK_CODE)) {
                    throw new Error(AUTH_ERROR);
                }
                const userData = responses.map((resp) => resp.data);
                return userData[0];
            })
            .then((userData) => {
                globalEventBus.emit(busEvents.SET_PROFILE_DATA, userData);
            })
            .catch((err) => {
                if (err.message === AUTH_ERROR) {
                    globalRouter.activate(PATHS.login);
                }
            });
    }
}
