/**
 * Базовый класс компонента
 */
export class Component {
    /**
     * Конструктор компонента
     * @param {Object} parent - родитель компонента
     * @param {Object} state - перечиление полей-состояний компонента
     */
    constructor(parent, state) {
        this.parent = parent;
        this.state = state;
    }

    /**
     * Отрисовка компонента
     */
    render() {
        this.parent?.insertAdjacentHTML('beforeend', this.tmpl());
    }

    /**
     * Скрыть компонент
     */
    hide() {
        this.removeEventListeners();
        this.parent.innerHTML = '';
    }

    /**
     * HTML-код компонента
     */
    tmpl() {
    }

    /**
     * Установить листенеры компоненту
     */
    setEventListeners() {
    }

    /**
     * Убрать листенеры компонента
     */
    removeEventListeners() {
    }
}
