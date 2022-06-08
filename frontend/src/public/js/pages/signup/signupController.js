import SignupModel from './signupModel';
import SignupView from './signupView';


/**
 * Контроллер страницы регистрации
 */
export default class SignupController {
    /**
     * Конструктор
     * @param {Element} parent - элемент для рендера
     */
    constructor(parent) {
        this.view = new SignupView(parent);
        this.model = new SignupModel();
    }
}
