import AbstractView from "./AbstractView.js";
import TagEditor from "../pkg/tag-editor.js";
import fetcher from "../pkg/fetcher.js";
import redirect from "../index.js";


const path = "/api/post/create"

const createPost = async (title, text, tags) => {
    let body = {
        "title" : title,
        "data": text,
        "tags" : tags,
    }

    const data = await fetcher.post(path, body)
    if (data && data.msg !== undefined){
        let showErr = document.getElementById("showError")
        showErr.innerHTML = data.msg
        return
    }
    redirect.navigateTo(`/post/${data.post_id}`)
}

export default class extends AbstractView{
    constructor(params){
        super(params);
        this.setTitle("Create-post");
    }
    async getHtml(){
        return `
        <style>
        .form-createPost {
          max-width: 90wh;
          padding: 15px;
        }
        
        .form-createPost .form-floating:focus-within {
          z-index: 2;
        }
        .form-createPost textarea {
            width: 100%;
            height: 500px;
            resize: none;
        }
         .form-createPost .textblock__input {
            box-sizing: border-box;
            display: flex;
            flex-wrap: wrap;
            width: 100%;
            height: min-content;
            resize: none;
            overflow: hidden;
            color: #000000;
            font-size: 15px;
            border: 1px solid #a7a7a7;
            border-radius: 3px;
            padding: 5px;
            background-color: #ffffff;
        }
        
         .form-createPost .textblock__input::placeholder {
            color: rgba(240, 246, 252, .1);
        }
        
         .form-createPost .textblock__input>input {
            box-sizing: inherit;
            border: none;
            outline: none;
            background-color: #e7e7e7;
            color: inherit;
            min-width: 50px;
            flex: 1;
        }
        
         .form-createPost .textblock__input:focus,
         .form-createPost .textblock__input:hover,
         .form-createPost .textblock__input:active {
            border-color: #58a6ff;
            outline: none;
            box-shadow: rgba(0, 0, 0, .05) 0 0 0px 4px, #388BFD26 0 0 0px 4px;
        }

         .form-createPost .btn-tag {
            box-sizing: border-box;
            display: flex;
            width: max-content;
            height: max-content;
            background-color: #388BFD26;
            border-radius: 3px;
            padding: 5px;
            margin: 2px;
            font-size: 15px;
            text-decoration: none;
            color: #58a6ff;
        }
        
         .form-createPost .btn-tag>.remove {
            box-sizing: inherit;
            cursor: pointer;
            margin-left: 5px;
        }
        
         .form-createPost .btn-tag>.remove:hover {
            color: #388BFD26;
        }
        
         .form-createPost .btn-tag:hover {
            box-shadow: rgba(0, 0, 0, 0.1) 0px 0px 100px 100px inset;
        }
    </style>
    <main class="form-createPost w-100 m-auto">
        <div class="container">
        <form id="form-createPost" class="form-createPost <text-center>" onsubmit="return false;">
        <div class="mb-3">
            <label for="TitleInput" class="form-label">Titile</label>
            <input maxlength="58" name="title" type="text" class="form-control" id="TitleInput" required>
            <div class="form-text">Maximum of 58 characters</div>
        </div>
        <div class="mb-3">
            <label for="TextInput" class="form-label">Text</label>
            <textarea name="text" maxlength="10000" id="TextInput" rows="3" required></textarea>
            <div class="form-text">Maximum of 10000 characters</div>
        </div>
        <div class="textblock__input" id="b_TagEditor">
            <input type="text" placeholder="...">
        </div>
        <textarea type="text" name="tags" class="d-none" id="tb_TagEditor"></textarea>
        <div class="form-text">Maximum of 5 tags and length of 1 tag maximum of 16 characters</div>
        <button class="btn btn-primary">Post</button>
        <div id="showError"></div>
        </form>
        </div>
    </main>
        `;
    }
    async init() {
        const tagEditor = new TagEditor({
            BlockSelectorName: '#b_TagEditor',
            TextBlockSelectorName: '#tb_TagEditor',
            HasDoubles: false,
            ToLower: true,
            MaxTags: 5,
            Tags: ["LIGHT"]
        });
        const signInForm = document.getElementById("form-createPost")
        signInForm.addEventListener("submit", function () {
            const title = document.getElementById("TitleInput").value
            const data = document.getElementById("TextInput").value
            const tagElements = document.querySelectorAll(".btn-tag")
            const tags = [];

        tagElements.forEach(tagElement => {
            const tagValue = tagElement.querySelector('span').textContent;
            tags.push(tagValue);
        });
            createPost(title, data, tags)
        })
    }
};