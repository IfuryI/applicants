(function() {
  var template = Handlebars.template, templates = Handlebars.templates = Handlebars.templates || {};
templates['profile.hbs'] = template({"1":function(container,depth0,helpers,partials,data) {
    return "                        <button id=\"button-profile-settings\" class=\"head__settings-button\">\n                            <a class=\"settings-button__settings-link\" href=\"/settings\">Настройки</a>\n                        </button>\n";
},"compiler":[8,">= 4.3.0"],"main":function(container,depth0,helpers,partials,data) {
    var stack1, helper, alias1=depth0 != null ? depth0 : (container.nullContext || {}), alias2=container.hooks.helperMissing, alias3="function", alias4=container.escapeExpression, lookupProperty = container.lookupProperty || function(parent, propertyName) {
        if (Object.prototype.hasOwnProperty.call(parent, propertyName)) {
          return parent[propertyName];
        }
        return undefined
    };

  return "<div id=\"navbar\"></div>\n\n<div class=\"main\">\n    <div class=\"main__profile-card\">\n        <div class=\"profile-card__container\">\n            <div class=\"container__avatar-container\">\n                <img class=\"avatar-container__avatar\" src=\""
    + alias4(((helper = (helper = lookupProperty(helpers,"avatar") || (depth0 != null ? lookupProperty(depth0,"avatar") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"avatar","hash":{},"data":data,"loc":{"start":{"line":7,"column":59},"end":{"line":7,"column":69}}}) : helper)))
    + "\" alt=\"\">\n            </div>\n\n            <div class=\"profile-card__content\">\n                <div class=\"content__head\">\n                    <div class=\"head__first-row\">\n                        <p id=\"user-full-name\" class=\"head__user-name\">"
    + alias4(((helper = (helper = lookupProperty(helpers,"username") || (depth0 != null ? lookupProperty(depth0,"username") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"username","hash":{},"data":data,"loc":{"start":{"line":13,"column":71},"end":{"line":13,"column":83}}}) : helper)))
    + "</p>\n                    </div>\n                    <span id=\"user-email\" class=\"head__user-email\">"
    + alias4(((helper = (helper = lookupProperty(helpers,"email") || (depth0 != null ? lookupProperty(depth0,"email") : depth0)) != null ? helper : alias2),(typeof helper === alias3 ? helper.call(alias1,{"name":"email","hash":{},"data":data,"loc":{"start":{"line":15,"column":67},"end":{"line":15,"column":76}}}) : helper)))
    + "</span>\n"
    + ((stack1 = lookupProperty(helpers,"if").call(alias1,(lookupProperty(helpers,"eq")||(depth0 && lookupProperty(depth0,"eq"))||alias2).call(alias1,(depth0 != null ? lookupProperty(depth0,"username") : depth0),(lookupProperty(helpers,"getUsername")||(depth0 && lookupProperty(depth0,"getUsername"))||alias2).call(alias1,{"name":"getUsername","hash":{},"data":data,"loc":{"start":{"line":16,"column":39},"end":{"line":16,"column":52}}}),{"name":"eq","hash":{},"data":data,"loc":{"start":{"line":16,"column":26},"end":{"line":16,"column":53}}}),{"name":"if","hash":{},"fn":container.program(1, data, 0),"inverse":container.noop,"data":data,"loc":{"start":{"line":16,"column":20},"end":{"line":20,"column":27}}})) != null ? stack1 : "")
    + "                </div>\n\n            </div>\n        </div>\n    </div>\n</div>\n";
},"useData":true});
})();