import ProfileModel from './profileModel';
import ProfileView from './profileView';


/**
 * Контроллер страницы профиля
 */
export default class ProfileController {
    /**
     * Конструктор
     * @param {Element} parent - элемент для рендера
     */
    constructor(parent) {
        this.view = new ProfileView(parent);
        this.model = new ProfileModel();
    }
}
