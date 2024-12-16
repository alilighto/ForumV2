import AbstractView from "./AbstractView.js";
import fetcher from "../pkg/fetcher.js";
import redirect from "../index.js";

const path = "/api/post/create";

const createPost = async (title, text, categories) => {
  let body = {
    title: title,
    data: text,
    categories: categories,
  };

  const data = await fetcher.post(path, body);
  if (data && data.msg !== undefined) {
    let showErr = document.getElementById("showError");
    showErr.innerHTML = data.msg;
    return;
  }
  redirect.navigateTo(`/post/${data.post_id}`);
};

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Create-post");
  }

  async getHtml() {
    const categories = await this.getCategories();

    return `
    <main class="form-createPost w-100 m-auto">
        <div class="container">
        <form id="form-createPost" class="form-createPost <text-center>" onsubmit="return false;">
        <div class="mb-3">
            <label for="TitleInput" class="form-label">Title</label>
            <input maxlength="58" name="title" type="text" class="form-control" id="TitleInput" required>
            <div class="form-text">Maximum of 58 characters</div>
        </div>
        <div class="mb-3">
            <label for="TextInput" class="form-label">Text</label>
            <textarea name="text" maxlength="10000" id="TextInput" rows="3" required></textarea>
            <div class="form-text">Maximum of 10000 characters</div>
        </div>
        <div class="mb-3">
            <label for="CategoryInput" class="form-label">Category</label>
            <select id="CategoryInput" class="form-control" required>
                <option value="" disabled selected>Select a category</option>
                ${categories
                  .map(
                    (category) =>
                      `<option value="${category.name}">${category.name}</option>`
                  )
                  .join("")}
            </select>
            <div class="form-text">Select one category</div>
        </div>
        <button class="btn btn-primary">Post</button>
        <div id="showError" class="error-message"></div>
        </form>
        </div>
    </main>
        `;
  }

  async init() {
    const createPostForm = document.getElementById("form-createPost");
    createPostForm.addEventListener("submit", function (e) {
      e.preventDefault();

      const title = document.getElementById("TitleInput").value.trim();
      const data = document.getElementById("TextInput").value.trim();
      const categorySelect = document.getElementById("CategoryInput");
      const category = categorySelect.value;

      // Validate inputs
      if (!title || !data || !category) {
        document.getElementById("showError").textContent =
          "Please fill in all fields";
        return;
      }

      createPost(title, data, [category]);
    });
  }
}
