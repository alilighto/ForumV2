import Home from "./views/HomeView.js";
import SignIn from "./views/SignInView.js";
import SignUp from "./views/SignUpView.js";
import CreatePost from "./views/CreatePostView.js";
import Post from "./views/PostView.js";
import Profile from "./views/ProfileView.js";
import NavBar from "./views/NavBarView.js";
import Utils from "./pkg/Utils.js";
import fetcher from "./pkg/fetcher.js";

const pathToRegex = path => new RegExp("^" + path.replace(/\//g, "\\/").replace(/:\w+/g, "(.+)") + "$");

const roles = {
    guest: 0,
    user: 1,
}

const getParams = match => {
    const values = match.result.slice(1);
    const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map(result => result[1]);

    return Object.fromEntries(keys.map((key, i) => {
        return [key, values[i]];
    }));
};


const navigateTo = url => {
    history.pushState(null, null, url);
    router();
};

const router = async () => {
    const routes = [
        { path: "/", view: Home, minRole: roles.guest},
        { path: "/sign-in", view: SignIn, minRole: roles.guest},
        { path: "/sign-up", view: SignUp, minRole: roles.guest},
        { path: "/create-post", view: CreatePost, minRole: roles.user},
        { path: "/post/:postID", view: Post, minRole: roles.guest},
        { path: "/user/:userID", view: Profile, minRole: roles.user},
        
    ];

    const potentialMatches = routes.map(route =>{
        return {
            route: route,
            result : location.pathname.match(pathToRegex(route.path))
        };
    });

    const checker = await fetcher.checkToken()
    if (checker &&!checker.checker){
        localStorage.setItem('role', roles.guest);
        localStorage.removeItem('id')
        localStorage.removeItem('token')
    }
    const user = Utils.getUser()

    if (!user.role) {
        user.role = roles.guest;
        localStorage.setItem('role', user.role);
    };
    const NavBarView = new NavBar(null, user);
    document.querySelector("#navbar").innerHTML = await NavBarView.getHtml();
    NavBarView.init();
    
    let match = potentialMatches.find(potentialMatches => potentialMatches.result !== null)
    if (!match) {
        Utils.showError(404)
        return
    }
    
    if (user.role < match.route.minRole) {
        Utils.showError(401, 'Please sign in to get access for this page')
        return  
    };

    const view = new match.route.view(getParams(match), user);
    document.querySelector("#app").innerHTML = await view.getHtml();
    view.init();
};

window.addEventListener("popstate", router)

window.addEventListener("storage", () => {
    const user = Utils.getUser()
    if (user.id == null) {
        location.reload()
    }
})


document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", e => {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            navigateTo(e.target.href);
        }
    });

    router();
});

export default {navigateTo, roles};