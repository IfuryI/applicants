import ModelingModel from './modelingModel';
import ModelingView from './modelingView';


/**
 * Контроллер страницы авторизации
 */
export default class ModelingController {
    /**
     * Конструктор
     * @param {Element} parent - элемент для рендера
     */
    constructor(parent) {
        this.view = new ModelingView(parent);
        this.model = new ModelingModel();
    }
}
