import {PATHS} from 'utils/paths';
import {globalRouter} from 'utils/router';
import ProfileController from './public/js/pages/profile/profileController';
import LoginController from './public/js/pages/login/loginController';
import SignupController from './public/js/pages/signup/signupController';
import SettingsController from './public/js/pages/settings/settingsController';
import {registerHandlebarsHelpers} from 'utils/handlebarsHelpers';
import ModelingController from './public/js/pages/modeling/modelingController';
import RequestController from './public/js/pages/request/requestController';
import './public/scss/compound.scss';
import {userMeta} from 'utils/userMeta';

window.addEventListener('DOMContentLoaded', async () => {
    const application = document.getElementById('app');

    await userMeta.initMeta();
    
    const source = new EventSource("http://localhost:3000")
    source.onmessage = (msg) => {
      var x = document.getElementById('one');
      x.value += msg.data + "\n";

      console.log(msg.data)
    }

    registerHandlebarsHelpers();

    globalRouter.register(PATHS.main, new ModelingController(application).view);
    globalRouter.register(PATHS.request, new RequestController(application).view);
    globalRouter.register(PATHS.user, new ProfileController(application).view);
    globalRouter.register(PATHS.login, new LoginController(application).view);
    globalRouter.register(PATHS.signup, new SignupController(application).view);
    globalRouter.register(PATHS.settings, new SettingsController(application).view);

    globalRouter.start();
});

