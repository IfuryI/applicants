import SettingsModel from './settingsModel';
import SettingsView from './settingsView';


/**
 * Контроллер страницы настроек профиля
 */
export default class SettingsController {
    /**
     * Конструктор
     * @param {Element} parent - элемент для рендера
     */
    constructor(parent) {
        this.view = new SettingsView(parent);
        this.model = new SettingsModel();
    }
}
