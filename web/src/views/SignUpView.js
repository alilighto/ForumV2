import AbstractView from "./AbstractView.js";
import redirect from "../index.js";
import fetcher from "../pkg/fetcher.js";

const path = "/api/signup"

const signup = async (email, username, password, rePassword) => {
    let body = {
        "email" : email,
        "username": username,
        "password" : password,
        "cfmpsw" : rePassword
    }

    const data = await fetcher.post(path, body)
    if (data && data.msg !== undefined){
        let showErr = document.getElementById("showError")
        showErr.innerHTML = data.msg
        return
    }
    redirect.navigateTo('/sign-in')
}

export default class extends AbstractView{
    constructor(params){
        super(params);
        this.setTitle("Sign-up");
    }
    async getHtml(){
        return `
        <style>
        .form-signup {
        max-width: 400px;
        padding: 15px;
        }
        
        .form-signup .form-floating:focus-within {
        z-index: 2;
        }
        
        .form-signup input[type="email"] {
        margin-bottom: -1px;
        border-bottom-right-radius: 0;
        border-bottom-left-radius: 0;
        }
        
        .form-signup input[type="password"] {
        border-top-left-radius: 0;
        border-top-right-radius: 0;
        }
        .form-signup .rePassword{
            margin-bottom: 10px;
        }
    </style>
    <main class="form-signup w-100 m-auto">
        <form id="form-signup" class="form-signup text-center" onsubmit="return false;">
            <h1 class="h1 mb-3 fw-normal">Please sign up</h1>
            <div class="form-floating">
                <input type="email" class="form-control" id="email" placeholder="name@example.com">
                <label for="email">Email address</label>
            </div>
            <div class="form-floating">
                <input type="text" class="form-control" id="username" placeholder="user">
                <label for="username">Username</label>
            </div>
            <div class="form-floating">
                <input type="password" class="form-control" id="password" placeholder="Password">
                <label for="password">Password</label>
            </div>
            <div class="form-floating rePassword">
                <input type="password" class="form-control" id="rePassword" placeholder="Password">
                <label for="rePassword">Repeat Password</label>
            </div>
            <button class="w-100 btn btn-lg btn-primary" type="submit">Sign up</button>
            <div id="showError"></div>
        </form>
    </main>
        `;
    }
    async init() {
        const signUpForm = document.getElementById("form-signup")
        signUpForm.addEventListener("submit", function () {
            const email = document.getElementById("email").value
            const username = document.getElementById("username").value
            const password = document.getElementById("password").value
            const rePassword = document.getElementById("rePassword").value
            signup(email, username,password, rePassword)
        })
    }
}