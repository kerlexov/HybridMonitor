import {Cookies} from "quasar";
import axios from "axios";

export default axios.create({baseURL:  'http://127.0.0.1:9393', headers:{Authorization: 'Token ' + Cookies.get('token')}});
