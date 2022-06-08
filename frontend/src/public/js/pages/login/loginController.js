import LoginModel from './loginModel';
import LoginView from './loginView';


/**
 * Контроллер страницы авторизации
 */
export default class LoginController {
    /**
     * Конструктор
     * @param {Element} parent - элемент для рендера
     */
    constructor(parent) {
        this.view = new LoginView(parent);
        this.model = new LoginModel();
    }
}
