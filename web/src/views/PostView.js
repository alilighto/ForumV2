import fetcher from "../pkg/fetcher.js";
import AbstractView from "./AbstractView.js";

const getPostPath = "/api/post/";
const sendCommentPath = "/api/comment/create";

const getPost = async (postID) => {
  const post = await fetcher.get(getPostPath + postID);
  if (post && post.msg != undefined) {
    return;
  }
  if (post) {
    document.getElementById("post-title").innerText = "Title: " + post.title;
    const userEl = document.getElementById("post-user-id");
    userEl.innerText = "Author: " + post.username;
    userEl.setAttribute("href", `/user/${post.user_id}`);
    const formattedCategories = post?.categories?.map((tag) => " #" + tag);
    document.getElementById("post-tags").innerText =
      "Categories:" + formattedCategories?.slice(0, -1);

    document.getElementById("post-data").innerText = post.data;
    document.getElementById("post-like-inner").innerText = post.likes;
    document.getElementById("post-dislike-inner").innerText = post.dislikes;

    const likeBtn = document.getElementById("post-like");
    likeBtn.addEventListener("click", () => {
      votePost(postID, 1);
    });

    const dislikeBtn = document.getElementById("post-dislike");
    dislikeBtn.addEventListener("click", () => {
      votePost(postID, 0);
    });

    const commentsDoc = document.getElementById("comments");
    commentsDoc.innerHTML = ""; // Clear previous comments

    if (post.comments.length > 0) {
      const commentText = document.createElement("h3");
      commentText.innerText = "Comments: ";
      commentsDoc.append(commentText);
    }

    for (let i = post.comments.length - 1; i >= 0; i--) {
      const comment = post.comments[i];
      const el = drawComments(comment);
      commentsDoc.append(el);
    }
  }
};

const votePost = async (postID, likeType) => {
  const path = "/api/post/vote";
  const body = {
    post_id: parseInt(postID),
    vote: likeType,
  };
  const data = await fetcher.post(path, body);
  if (data && data.msg) {
    return;
  }
  await getPost(postID);
};

const voteComment = async (commentID, likeType) => {
  const path = "/api/comment/vote";
  const body = {
    comment_id: parseInt(commentID),
    vote: likeType,
  };
  const data = await fetcher.post(path, body);
  if (data && data.msg) {
    return;
  }

  const postIDFromURL = window.location.pathname.split("/").pop();

  if (postIDFromURL) {
    await getPost(postIDFromURL);
  } else {
    console.error("Could not find post ID");
  }
};

const drawComments = (comment) => {
  const el = document.createElement("div");
  el.classList.add("card");

  const authorEl = document.createElement("a");
  authorEl.classList.add("card-header");
  authorEl.setAttribute("href", `/user/${comment.user_id}`);
  authorEl.setAttribute("data-link", "");
  authorEl.innerText = "Author: " + comment.username;

  const body = document.createElement("div");
  body.classList.add("card-body");

  const dataEl = document.createElement("p");
  dataEl.classList.add("card-text");
  dataEl.innerText = comment.data;
  body.append(dataEl);

  let likeButton = document.createElement("button");
  likeButton.className = "btn comment-like";
  likeButton.id = "comment-like";
  let likeIcon = document.createElement("i");
  likeIcon.className = "fa fa-thumbs-up fa-lg";
  likeIcon.id = "comment-like-inner";
  likeIcon.setAttribute("aria-hidden", "true");
  likeIcon.innerText = comment.likes;
  likeButton.appendChild(likeIcon);

  let dislikeButton = document.createElement("button");
  dislikeButton.className = "btn comment-dislike";
  dislikeButton.id = "comment-dislike";
  let dislikeIcon = document.createElement("i");
  dislikeIcon.className = "fa fa-thumbs-down fa-lg";
  dislikeIcon.id = "comment-dislike-inner";
  dislikeIcon.setAttribute("aria-hidden", "true");
  dislikeIcon.innerText = comment.dislikes;
  dislikeButton.appendChild(dislikeIcon);

  const votes = document.createElement("div");
  votes.appendChild(likeButton);
  votes.appendChild(dislikeButton);

  likeButton.addEventListener("click", () => {
    voteComment(comment.comment_id, 1);
  });
  dislikeButton.addEventListener("click", () => {
    voteComment(comment.comment_id, 0);
  });

  el.append(authorEl);
  el.append(body);
  const hr = document.createElement("hr");
  el.append(hr);
  el.append(votes);
  return el;
};

const sendComment = async (comment, postID) => {
  let body = {
    data: comment,
    post_id: parseInt(postID),
  };
  const data = await fetcher.post(sendCommentPath, body);
  if (data && data.msg !== undefined) {
    let showErr = document.getElementById("showError");
    showErr.innerHTML = data.msg;
    return;
  }
  await getPost(postID);
};

export default class extends AbstractView {
  constructor(params, user) {
    super(params);
    this.user = user;
    this.setTitle("Post");
  }

  async getHtml() {
    const isAuthorized = Boolean(this.user.id);
    return `
        <main class="main-content">
            <div class="post-container">
                <div class="post-details">
                    <h3 id="post-title"></h3>
                    <div id="post-user">
                        <a id="post-user-id"></a>
                    </div>
                    <h5 id="post-tags"></h5>
                    <div class="post-content">
                        <p id="post-data"></p>
                    </div>
                    <div class="post-actions">
                        <button class="btn post-like" id="post-like">
                            <i class="fa fa-thumbs-up fa-lg" id="post-like-inner" aria-hidden="true"></i>
                        </button>
                        <button class="btn post-dislike" id="post-dislike">
                            <i class="fa fa-thumbs-down fa-lg" id="post-dislike-inner" aria-hidden="true"></i>
                        </button>
                    </div>
                </div>
            </div>

            ${
              isAuthorized
                ? `
            <div class="comment-section">
                <form id="comment-form" class="comment-form">
                    <h3>Leave a comment here:</h3>
                    <textarea 
                        id="comment-input" 
                        class="form-control" 
                        rows="3" 
                        placeholder="Leave a comment"
                    ></textarea>
                    <button type="submit" class="btn btn-primary">Send</button>
                    <div class="error" id="showError"></div>
                </form>
            </div>
            `
                : ""
            }

            <div id="comments" class="comments-container">
                <!-- Comments will be dynamically populated here -->
            </div>
        </main>
        `;
  }

  async init() {
    const postID = this.params.postID;
    await getPost(postID);

    // Comment submission for authorized users
    const commentForm = document.getElementById("comment-form");
    commentForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const commentInput = document.getElementById("comment-input");
      const comment = commentInput.value.trim();
      if (comment) {
        await sendComment(comment, postID);
        commentInput.value = ""; // Clear the input after submission
      }
    });
  }
}
