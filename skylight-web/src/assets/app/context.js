// import { reject } from "core-js/fn/promise"
import API from "./api"

export class Context {
    constructor(data={}) {
        this.cluster = data.cluster
        this.region = data.region
        this.project = data.project
        this.user = data.user
        this.roles = data.roles || []
    }
    setCluster(cluster) {
        this.cluster = cluster
        this.save()
    }
    setRegion(region) {
        this.region = region
        this.save()
    }
    isAdmin() {
        return this.roles.indexOf('admin') >= 0
    }
    save() {
        let data = JSON.stringify(this)
        localStorage.setItem('context', data)
    }
}

function loadFromLocalStorage() {
    let data = localStorage.getItem('context')
    if (!data) {
        return null
    }
    return new Context(JSON.parse(data))
}
async function loadFromServer() {
    let auth = (await API.system.isLogin()).auth
    let roles = []
    for (let i in auth.roles) { roles.push(auth.roles[i].name) }
    return new Context({
        cluster: auth.cluster, region: auth.region,
        project: auth.project, user: auth.user, roles: roles,
    })
}

export async function GetContext() {
    let context = await loadFromServer()
    context.save()
    return context
}
export function GetLocalContext() {
    return loadFromLocalStorage()
}
