import AbstractView from "./AbstractView.js";
import Utils from "../pkg/Utils.js";

export default class extends AbstractView{
    constructor(params, user){
        super(params);
        this.user = user
    }
    async getHtml(){
        const isAuthorized = Boolean(this.user.id)
        return `
        <div class="container-fluid">
            <ul class="navbar-nav me-auto mb-2 mb-sm-0">
                <li class="nav-item">
                    <a href="/" class="nav-link" id="home-button" data-link>Home</a>
                </li>
        `+ (isAuthorized ?
        `
                <li class="nav-item">
                    <a href="/create-post" class="nav-link" id="create-post-button" data-link>Create Post</a>
                </li>
                <li class="nav-item">
                    <a href="/user/${this.user.id}" class="nav-link" id="Profile-button" data-link>Profile</a>
                </li>
                <li class="nav-item">
                    <a href="/" class="nav-link" id="sign-out-button" data-link>Sign Out</a>
                </li>
            </ul> 
        </div>
        `:
        `       <li class="nav-item">
                    <a href="/sign-up" class="nav-link" id="sign-in-button" data-link>Sign Up</a>
                </li>
                <li class="nav-item">
                    <a href="/sign-in" class="nav-link" id="sign-up-button" data-link>Sign In</a>
                </li>
            </ul>   
        </div>
        `
        );
    }
    async init() {
        document.body.addEventListener('click', (event) => {
            if (event.target.id === 'sign-out-button') {
                console.log("log");
                Utils.logOut();
                window.location.reload();
            }
        });
    }
};