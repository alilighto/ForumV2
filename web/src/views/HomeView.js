import AbstractView from "./AbstractView.js";
import fetcher from "../pkg/fetcher.js";

const path = `/api/posts/`

const getPostsByCategory = async (category) =>{
    const data = await fetcher.get(path + category)
    if (data && data.msg !== undefined){
        console.log(data)
        return
    }else{
        const postsDoc = document.getElementById("posts")
        postsDoc.textContent = "";
        for (let i = data.length - 1; i >= 0; i--) {
            const post = data[i];
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

export default class extends AbstractView{
    constructor(params){
        super(params);
        this.setTitle("Forum");
    }
    async getHtml(){
        return `
        <style>
        .card{
            -webkit-box-shadow: 0 2px 2px 0 rgba(0, 0, 0, 0.2), 0 6px 10px 0 rgba(0, 0, 0, 0.3);
            box-shadow: 0 2px 2px 0 rgba(0, 0, 0, 0.2), 0 6px 10px 0 rgba(0, 0, 0, 0.3);
        }
        </style>
        <header class="py-3 mb-4 border-bottom">
            <div class="container d-flex flex-wrap justify-content-center">
            <form id="form-search" class="w-100 me-3" onsubmit="return false;">
                <input id="search" type="search" class="form-control" placeholder="Search by category" aria-label="Search">
            </form>
            </div>
        </header>
        <div id="posts"></div>
        `;
    }
    async init() {
        getPostsByCategory("ALL")
        const signInForm = document.getElementById("form-search")
        signInForm.addEventListener("submit", function () {
            let category = document.getElementById("search").value
            if (category.trim() === ""){
                category = "ALL"
            }
            getPostsByCategory(category.trim())
        })
    }
}