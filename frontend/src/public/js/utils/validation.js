'use strict';


/**
 * Класс для валидации данных
 */
export default class Validator {
    /**
     * Проверка email
     * @param {string} email - адрес почты
     * @return {[string]} - список ошибок валидации
     */
    validateEmail(email) {
        const error = [];

        if (!(/@/.test(email))) {
            error.push('Некорректный формат Email');
        }

        return error;
    }

    /**
     * Проверка login
     * @param {string} login - логин
     * @return {[string]} - список ошибок валидации
     */
    validateLogin(login) {
        const error = [];

        if (!(login.length >= 6 && login.length <= 20)) {
            error.push('Логин должен быть длиной от 6 до 20 символов!');
        }
        if (!(/^[a-zA-Z0-9_]*$/.test(login))) {
            error.push('Логин должен состоять только из латинских букв, чисел и нижних подчеркиваний!');
        }

        return error;
    }

    /**
     * Проверка пароля
     * @param {string} password - пароль
     * @return {[string]} - список ошибок валидации
     */
    validatePassword(password) {
        const error = [];

        if (!(password.length >= 8 && password.length <= 32)) {
            error.push('Пароль должен быть длиной от 8 до 32!');
        }
        if (!(/\d/.test(password))) {
            error.push('Пароль должен содержать хотя бы одну цифру!');
        }
        if (!(/[a-zA-Z]/.test(password))) {
            error.push('Пароль должен содержать хотя бы одну букву!');
        }

        return error;
    }

    /**
     * Проверка аватара
     * @param {string} avatar - аватар
     * @return {[string]} - список ошибок валидации
     */
    validateAvatar(avatar) {
        const error = [];

        if (!avatar) {
            error.push('Нельзя выбрать пустой аватар!');
            return error;
        }

        if (!(['image/png', 'image/jpeg'].includes(avatar.type))) {
            error.push('Некорректный формат картинки! Используйте png или jpeg');
        }

        if (!(avatar.size <= 5242880)) {
            error.push('Размер файла не должен превышать 5MB!');
        }

        return error;
    }

    /**
     * Проверка рецензии
     * @param {Object} review - поля рецензии
     * @return {[string]} - список ошибок валидации
     */
    validateReview(review) {
        const error = [];

        if (review.review_type === undefined || review.title.length === 0 || review.content.length === 0) {
            error.push('Заполните все поля!');
        }

        return error;
    }
}
