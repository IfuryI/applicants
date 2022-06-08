import RequestModel from './requestModel';
import RequestView from './requestView';


/**
 * Контроллер страницы авторизации
 */
export default class RequestController {
    /**
     * Конструктор
     * @param {Element} parent - элемент для рендера
     */
    constructor(parent) {
        this.view = new RequestView(parent);
        this.model = new RequestModel();
    }
}
