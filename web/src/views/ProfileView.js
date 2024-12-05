import Utils from "../pkg/Utils.js";
import fetcher from "../pkg/fetcher.js";
import AbstractView from "./AbstractView.js";

const getUserByID = async (userID) => {
    const path = `/api/profile/${userID}`
    const data = await fetcher.get(path)
    if (data && data.msg != undefined){
        console.log(data)
        return
    }
    if (data) {
        drawUser(data)
    }else {
        console.log(data)
    }
}

const getUserLikedPosts = async (userID) =>{
    const path =`/api/profile/liked-posts/${userID}`
    const posts = await fetcher.get(path)
    if (posts && posts.msg != undefined){
        console.log(post)
        return
    }
    if (posts){
        const postsDoc = document.getElementById("posts")
        postsDoc.textContent = "";
        for (let i = posts.length - 1; i >= 0; i--) {
            const post = posts[i];
            const el = newPostElement(post);
            postsDoc.append(el);
        }
    }
}
const getUserDislikedPosts = async (userID) =>{
    const path =`/api/profile/disliked-posts/${userID}`
    const posts = await fetcher.get(path)
    if (posts && posts.msg != undefined){
        console.log(post)
        return
    }
    if (posts){
        const postsDoc = document.getElementById("posts")
        postsDoc.textContent = "";
        for (let i = posts.length - 1; i >= 0; i--) {
            const post = posts[i];
            const el = newPostElement(post);
            postsDoc.append(el);
        }
    }
}

const getUserPosts = async (userID) =>{
    const path =`/api/profile/posts/${userID}`
    const posts = await fetcher.get(path)
    if (posts && posts.msg != undefined){
        console.log(post)
        return
    }
    if (posts){
        const postsDoc = document.getElementById("posts")
        postsDoc.textContent = "";
        for (let i = posts.length - 1; i >= 0; i--) {
            const post = posts[i];
            const el = newPostElement(post);
            postsDoc.append(el);
        }
    }
}
const newPostElement = (post) =>{
    const el = document.createElement("div")
    el.classList.add("card")

    const titleEl = document.createElement("a")
    titleEl.classList.add("card-header")
    titleEl.setAttribute("href", `/post/${post.post_id}`)
    titleEl.setAttribute("data-link", "")
    titleEl.innerText = "Title: " + post.title

    const authorEl = document.createElement("a")
    authorEl.classList.add("card-header")
    authorEl.setAttribute("href", `/user/${post.user_id}`)
    authorEl.setAttribute("data-link", "")
    authorEl.innerText = "Author: " + post.username

    const body = document.createElement("div")
    body.classList.add("card-body")
    
    const tagsEl = document.createElement("h5")
    tagsEl.classList.add("card-title")
    for (let i = 0; i < post.tags.length; i++){
        post.tags[i] = " #" + post.tags[i] 
    }
    tagsEl.innerText = "Categories:" + post.tags.slice(0, -1)
    
    const dataEl = document.createElement("p")
    dataEl.classList.add("card-text")
    dataEl.innerText = post.data.substring(0, 300)+"..."

    body.append(tagsEl)
    body.append(dataEl)

    el.append(titleEl)
    el.append(authorEl)
    el.append(body)
    return el
}
const drawUser = (user) =>{
    document.getElementById("username").innerText = user.username
    document.getElementById("email").innerText = user.email
}

export default class extends AbstractView{
    constructor(params){
        super(params);
        this.setTitle("Profile");
    }
    async getHtml(){
        return `
        <style>
            .card{
                background-color: rgba(0,0,0,.03) !important;
            }
        </style>
        <div class="container justify-content-center align-items-center " style="max-width: 300px;">
            <div class="card justify-content-center align-items-center" style="max-width: 300px; ">
                <img src="/src/assets/img/profile.jpg" alt="profile image" class="img-fluid img-thumbnail mt-4 mb-2"
                style="width: 150px; z-index: 1">
                <div class="card-body">
                    <h5 id="username" class="card-header"></h5>
                    <h5 id="email" class="card-header"></h5>
                    <select id="options" class="form-select" aria-label="Default select example">
                        <option value="created">Created posts</option>
                        <option value="liked">Liked posts</option>
                        <option value="disliked">Disliked posts</option>
                    </select>
                </div>
                <div class="error" id="showError"></div>
            </div>
        </div>
        <hr>
        <div id="posts"></div>
        `;
    }
    async init() {
        const userID = this.params.userID
        getUserByID(userID)
        getUserPosts(userID)
        // Находим элемент select по его id
        let selectElement = document.getElementById("options");

        // Добавляем обработчик события "change" для элемента select
        selectElement.addEventListener("change", function() {
            // Получаем выбранное значение
            let selectedValue = selectElement.value;
            
            // Делаем что-то с полученным значением
            if (selectedValue == "created"){
                getUserPosts(userID)
            } else if (selectedValue == "liked"){
                getUserLikedPosts(userID)
            } else if (selectedValue == "disliked") {
                getUserDislikedPosts(userID)
            } else {
                Utils.showError(400, "invalid option")
            }
        });
    }
}