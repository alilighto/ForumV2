import AbstractView from "./AbstractView.js";
import redirect from "../index.js";
import fetcher from "../pkg/fetcher.js";
import Utils from "../pkg/Utils.js";

const path = `/api/signin`

const signIn = async (email, password) => {
    let body = {
        "email" : email,
        "password" : password
    }

    const data = await fetcher.post(path, body)
    if (data && data.msg !== undefined){
        let showErr = document.getElementById("showError")
        showErr.innerHTML = data.msg
        return
    }
    localStorage.setItem("token", data.token)
    const payload = Utils.parseJwt(data.token)
    localStorage.setItem("id", payload.id)
    localStorage.setItem("role", redirect.roles.user)
    redirect.navigateTo('/')
}

export default class extends AbstractView{
    constructor(params){
        super(params);
        this.setTitle("Sign-in");
    }
    async getHtml(){
        return `
        <style>
        .form-signin {
          max-width: 400px;
          padding: 15px;
        }
        
        .form-signin .form-floating:focus-within {
          z-index: 2;
        }
        
        .form-signin input[type="email"] {
          margin-bottom: -1px;
          border-bottom-right-radius: 0;
          border-bottom-left-radius: 0;
        }
        
        .form-signin input[type="password"] {
          margin-bottom: 10px;
          border-top-left-radius: 0;
          border-top-right-radius: 0;
        }
    </style>
    <main class="form-signin w-100 m-auto">
        <form id="form-signin" class="form-signin text-center" onsubmit="return false;">
            <h1 class="h1 mb-3 fw-normal">Please sign in</h1>
            <div class="form-floating">
                <input type="email" class="form-control" id="email" placeholder="name@example.com">
                <label for="email">Email address</label>
            </div>
            <div class="form-floating">
                <input type="password" class="form-control" id="password" placeholder="Password">
                <label for="password">Password</label>
                </div>
            <button class="w-100 btn btn-lg btn-primary" type="submit">Sign in</button>
            <br/>
            <div id="showError"></div>
        </form>
    </main>
        `;
    }
    async init() {
        const signInForm = document.getElementById("form-signin")
        signInForm.addEventListener("submit", function () {
            const email = document.getElementById("email").value
            const password = document.getElementById("password").value
            signIn(email, password)
        })
    }
}