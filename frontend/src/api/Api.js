import Repo from "src/api/Repo"

export default {
  async get(api){
    return Repo.get(api);
  },
  async post(api, object){
    return Repo.post(api, object);
  },
  async put(api, object){
    return Repo.put(api, object);
  }
}
