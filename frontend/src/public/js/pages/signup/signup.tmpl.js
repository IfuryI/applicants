(function() {
  var template = Handlebars.template, templates = Handlebars.templates = Handlebars.templates || {};
templates['signup.hbs'] = template({"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    return "<div id=\"navbar\"></div>\n\n\n<div class=\"main\">\n    <div class=\"form-box\">\n        <div class=\"form-box__button-box\">\n            <div id=\"button-box__current-button\" class=\"button-box__button-wrap switcher-on-signup\"></div>\n            <a href=\"/login\">\n                <button type=\"button\" class=\"button-box__toggle-button\" id=\"login-button\">Войти</button>\n            </a>\n            <button type=\"button\" class=\"button-box__toggle-button button-box__button_text_white\" id=\"signup-button\">\n                Регистрация\n            </button>\n        </div>\n\n        <form id=\"signup\" class=\"form-box__form form-box__input-group\" novalidate>\n            <input name=\"username\" type=\"text\" class=\"form__input-field\" id=\"username-input\"\n                   placeholder=\"Имя пользователя\">\n            <p class=\"form__validation-hint\" id=\"validation-hint-login\"><em></em></p>\n            <input name=\"email\" type=\"email\" class=\"form__input-field\" id=\"email-input\" placeholder=\"Email\">\n            <p class=\"form__validation-hint\" id=\"validation-hint-email\"><em></em></p>\n            <input name=\"password\" type=\"password\" class=\"form__input-field\" id=\"password-input\" placeholder=\"Пароль\">\n            <p class=\"form__validation-hint\" id=\"validation-hint-password\"><em></em></p>\n            <p class=\"form__validation-hint form__send-form-hint\" id=\"validation-hint-signup\"><em></em></p>\n            <button type=\"submit\" id=\"signup-submit\" class=\"form__submit-button\">Зарегистрироваться</button>\n        </form>\n    </div>\n</div>\n";
},"useData":true});
})();