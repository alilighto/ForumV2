import AbstractView from "./AbstractView.js";
import fetcher from "../pkg/fetcher.js";

const getUserByID = async (userID) => {
  const path = `/api/profile/${userID}`;
  const data = await fetcher.get(path);
  if (data && data.msg != undefined) {
    console.log(data);
    return;
  }
  if (data) {
    drawUser(data);
  } else {
    console.log(data);
  }
};

const getUserPosts = async (userID) => {
  const path = `/api/profile/posts/${userID}`;
  const posts = await fetcher.get(path);
  if (posts && posts.msg != undefined) {
    console.log(posts);
    return;
  }
  if (posts) {
    const postsDoc = document.getElementById("posts");
    postsDoc.textContent = "";
    for (let i = posts.length - 1; i >= 0; i--) {
      const post = posts[i];
      const el = newPostElement(post);
      postsDoc.append(el);
    }
  }
};

const getUserLikedPosts = async (userID) => {
  const path = `/api/profile/liked-posts/${userID}`;
  const posts = await fetcher.get(path);
  if (posts && posts.msg != undefined) {
    console.log(posts);
    return;
  }
  if (posts) {
    const postsDoc = document.getElementById("posts");
    postsDoc.textContent = "";
    for (let i = posts.length - 1; i >= 0; i--) {
      const post = posts[i];
      const el = newPostElement(post);
      postsDoc.append(el);
    }
  }
};

const getUserDislikedPosts = async (userID) => {
  const path = `/api/profile/disliked-posts/${userID}`;
  const posts = await fetcher.get(path);
  if (posts && posts.msg != undefined) {
    console.log(posts);
    return;
  }
  if (posts) {
    const postsDoc = document.getElementById("posts");
    postsDoc.textContent = "";
    for (let i = posts.length - 1; i >= 0; i--) {
      const post = posts[i];
      const el = newPostElement(post);
      postsDoc.append(el);
    }
  }
};

const newPostElement = (post) => {
  const el = document.createElement("div");
  el.classList.add("card");

  const titleEl = document.createElement("a");
  titleEl.classList.add("card-header");
  titleEl.setAttribute("href", `/post/${post.post_id}`);
  titleEl.setAttribute("data-link", "");
  titleEl.innerText = "Title: " + post.title;

  const authorEl = document.createElement("a");
  authorEl.classList.add("card-header");
  authorEl.setAttribute("href", `/user/${post.user_id}`);
  authorEl.setAttribute("data-link", "");
  authorEl.innerText = "Author: " + post.username;

  const body = document.createElement("div");
  body.classList.add("card-body");

  const tagsEl = document.createElement("h5");
  tagsEl.classList.add("card-title");
  for (let i = 0; i < post.tags.length; i++) {
    post.tags[i] = " #" + post.tags[i];
  }
  tagsEl.innerText = "Categories:" + post.tags.slice(0, -1);

  const dataEl = document.createElement("p");
  dataEl.classList.add("card-text");
  dataEl.innerText = post.data.substring(0, 300) + "...";

  body.append(tagsEl);
  body.append(dataEl);

  el.append(titleEl);
  el.append(authorEl);
  el.append(body);
  return el;
};

const drawUser = (user) => {
  document.getElementById("username").innerText = user.username;
  document.getElementById("email").innerText = user.email;
};

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Profile");
  }

  async getHtml() {
    return `
        <main class="main-content">
            <div class="profile-container">
                <div class="profile-header">
                    <img 
                        src="/src/assets/img/profile.jpg" 
                        alt="profile image" 
                        class="profile-image"
                    >
                    <div class="profile-info">
                        <h2 id="username" class="profile-username"></h2>
                        <p id="email" class="profile-email"></p>
                    </div>
                </div>

                <div class="profile-actions">
                    <select id="posts-filter" class="form-select">
                        <option value="created">Created Posts</option>
                        <option value="liked">Liked Posts</option>
                        <option value="disliked">Disliked Posts</option>
                    </select>
                </div>

                <div id="posts" class="profile-posts">
                    <!-- Posts will be dynamically populated here -->
                </div>
            </div>
        </main>
        `;
  }

  async init() {
    const userID = this.params.userID;
    
    // Fetch user and initial posts
    await getUserByID(userID);
    await getUserPosts(userID);

    // Posts filter event listener
    const postsFilter = document.getElementById("posts-filter");
    postsFilter.addEventListener("change", async (e) => {
      const selectedValue = postsFilter.value;

      switch (selectedValue) {
        case "created":
          await getUserPosts(userID);
          break;
        case "liked":
          await getUserLikedPosts(userID);
          break;
        case "disliked":
          await getUserDislikedPosts(userID);
          break;
        default:
          console.error("Invalid filter option");
      }
    });
  }
}
