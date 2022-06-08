(function() {
  var template = Handlebars.template, templates = Handlebars.templates = Handlebars.templates || {};
templates['navbar.hbs'] = template({"1":function(container,depth0,helpers,partials,data) {
    var helper, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "            <div class=\"header__item header__button\">\n                <a href=\"/user/"
    + container.escapeExpression(((helper = (helper = lookupProperty(helpers,"getUsername") || (depth0 != null ? lookupProperty(depth0,"getUsername") : depth0)) != null ? helper : container.hooks.helperMissing),(typeof helper === "function" ? helper.call(depth0 != null ? depth0 : (container.nullContext || {}),{"name":"getUsername","hash":{},"data":data,"loc":{"start":{"line":17,"column":31},"end":{"line":17,"column":46}}}) : helper)))
    + "\">Профиль</a>\n            </div>\n            <div class=\"header__item header__button\">\n                <a id=\"logout-button\">Выйти</a>\n            </div>\n";
},"3":function(container,depth0,helpers,partials,data) {
    return "            <div class=\"header__item header__button\">\n                <a href=\"/login\">Войти</a>\n            </div>\n            <div class=\"header__item header__button\">\n                <a href=\"/signup\">Регистрация</a>\n            </div>\n";
},"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    var stack1, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "<div id=\"header\" class=\"header\">\n    <div class=\"header__section\">\n        <div class=\"header__item header__logo\">\n            <a href=\"\">Сервис приёма</a>\n        </div>\n        <div class=\"header__item header__button\">\n            <a href=\"/\">Моделирование</a>\n        </div>\n        <div class=\"header__item header__button\" id=\"reqbtn\">\n            <a href=\"/request\">Создание заявки в очередь</a>\n        </div>\n    </div>\n\n    <div class=\"header__section\">\n"
    + ((stack1 = lookupProperty(helpers,"if").call(depth0 != null ? depth0 : (container.nullContext || {}),(depth0 != null ? lookupProperty(depth0,"authorized") : depth0),{"name":"if","hash":{},"fn":container.program(1, data, 0),"inverse":container.program(3, data, 0),"data":data,"loc":{"start":{"line":15,"column":8},"end":{"line":29,"column":15}}})) != null ? stack1 : "")
    + "    </div>\n\n    <a href=\"javascript:void(0);\" id=\"bars-icon\" class=\"header__bars-icon\">\n        <i class=\"fa fa-bars fa-button\"></i>\n    </a>\n</div>\n";
},"useData":true});
})();