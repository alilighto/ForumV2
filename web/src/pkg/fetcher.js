import Utils from "./Utils.js"
import redirect from "../index.js";

const fetcher = {
    get: async (path, body) =>{
        return makeRequest(path, body, "GET")
    },
    post: async(path, body) =>{
        return makeRequest(path, body, "POST")
    },
    checkToken: async() =>{
        const url = `http://${API_HOST_NAME}/api/is-valid`
        const options = {
            mode: 'cors',
            method: "GET",
        }
        const token = localStorage.getItem("token")
        if (token != undefined){
            options.headers = new Headers({
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            })
        }
        const response = await fetch(url, options).catch((e) =>{
            console.log(e)
            Utils.showError(503)
            return
        })
        var responseBody
        try{
            responseBody = await response.json()
        }  catch{
            console.log("some unexpected error: with json")
            return
        }
        if (!response.ok){
            Utils.showError(response.status, responseBody.msg)
            return responseBody
        }
        return responseBody
    }
}
// bad request: UNIQUE constraint failed: users.username
const makeRequest = async(path, body, method) => {
    const url = `http://${API_HOST_NAME}${path}`
    const options = {
        mode: 'cors',
        method: method,
        body: JSON.stringify(body)
    }
    const token = localStorage.getItem("token")
    if (token != undefined){
        options.headers = new Headers({
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
        })
    }
    const response = await fetch(url, options).catch((e) =>{
        console.log(e)
        Utils.showError(503)
        return
    })

    var responseBody
    try{
        responseBody = await response.json()
    }  catch  {
        return
    }

    if (response.status == 401 || response.status == 403) {
        Utils.logOut()
        redirect.navigateTo("/sign-in")
        return responseBody
    }
    if (response.status == 404){
        Utils.showError(response.status)
        return
    }
    if (response.status == 400){
        return responseBody
    }
    if (!response.ok){
        Utils.showError(response.status, responseBody.msg)
        return responseBody
    }
    return responseBody
}


export default fetcher

