(function () {
  const API_BASE = "http://localhost:8080";

  const store = {
    state: {
      accessToken: localStorage.getItem("access_token") || "",
      refreshToken: localStorage.getItem("refresh_token") || "",
      role: localStorage.getItem("role") || "",
      notes: []
    },
    set(partial) {
      this.state = { ...this.state, ...partial };
      if ("accessToken" in partial) localStorage.setItem("access_token", this.state.accessToken || "");
      if ("refreshToken" in partial) localStorage.setItem("refresh_token", this.state.refreshToken || "");
      if ("role" in partial) localStorage.setItem("role", this.state.role || "");
    },
    clearAuth() {
      this.set({ accessToken: "", refreshToken: "", role: "" });
    }
  };

  function decodeRoleFromJWT(token) {
    try {
      const payload = token.split(".")[1];
      if (!payload) return "";
      const obj = JSON.parse(atob(payload.replace(/-/g, "+").replace(/_/g, "/")));
      return obj.role || "";
    } catch (_) {
      return "";
    }
  }

  function notify(message, isError) {
    const toast = document.getElementById("toast");
    toast.textContent = message;
    toast.className = "toast show" + (isError ? " error" : "");
    setTimeout(function () {
      toast.className = "toast";
    }, 2300);
  }

  async function apiRequest(path, options) {
    const init = { ...(options || {}) };
    init.headers = { "Content-Type": "application/json", ...(init.headers || {}) };
    if (store.state.accessToken) {
      init.headers.Authorization = "Bearer " + store.state.accessToken;
    }

    let res = await fetch(API_BASE + path, init);
    if (res.status === 401 && store.state.refreshToken && path !== "/refresh") {
      const refreshed = await refreshSession();
      if (refreshed) {
        init.headers.Authorization = "Bearer " + store.state.accessToken;
        res = await fetch(API_BASE + path, init);
      }
    }

    let body = {};
    const raw = await res.text();
    if (raw) {
      try {
        body = JSON.parse(raw);
      } catch (_) {
        body = { message: raw };
      }
    }

    if (!res.ok) {
      const msg = body.error || body.message || ("Request failed: " + res.status);
      throw new Error(msg);
    }
    return body;
  }

  const api = {
    register(payload) {
      return apiRequest("/register", { method: "POST", body: JSON.stringify(payload) });
    },
    login(payload) {
      return apiRequest("/login", { method: "POST", body: JSON.stringify(payload) });
    },
    refresh(payload) {
      return apiRequest("/refresh", { method: "POST", body: JSON.stringify(payload) });
    },
    logout(payload) {
      return apiRequest("/logout", { method: "POST", body: JSON.stringify(payload) });
    },
    listNotes() {
      return apiRequest("/api/notes");
    },
    createNote(payload) {
      return apiRequest("/api/notes", { method: "POST", body: JSON.stringify(payload) });
    },
    updateNote(id, payload) {
      return apiRequest("/api/notes/" + id, { method: "PUT", body: JSON.stringify(payload) });
    },
    deleteNote(id) {
      return apiRequest("/api/notes/" + id, { method: "DELETE" });
    }
  };

  async function refreshSession() {
    try {
      const data = await api.refresh({ refresh_token: store.state.refreshToken });
      const role = decodeRoleFromJWT(data.access_token);
      store.set({
        accessToken: data.access_token,
        refreshToken: data.refresh_token,
        role: role
      });
      return true;
    } catch (_) {
      store.clearAuth();
      if (location.hash !== "#/login") location.hash = "#/login";
      return false;
    }
  }

  function isAuthed() {
    return !!store.state.accessToken;
  }

  const app = document.getElementById("app");
  const logoutBtn = document.getElementById("logoutBtn");
  const editorModal = document.getElementById("editorModal");
  const editNoteForm = document.getElementById("editNoteForm");
  const closeModalBtn = document.getElementById("closeModalBtn");
  const cancelModalBtn = document.getElementById("cancelModalBtn");

  logoutBtn.addEventListener("click", async function () {
    try {
      if (store.state.refreshToken) {
        await api.logout({ refresh_token: store.state.refreshToken });
      }
    } catch (_) {
      // keep logout resilient even if backend token is already invalid
    }
    store.clearAuth();
    location.hash = "#/login";
    notify("Başarıyla çıkış yapıldı");
  });

  function closeEditorModal() {
    editorModal.classList.add("hidden");
    editorModal.setAttribute("aria-hidden", "true");
    editNoteForm.reset();
  }

  function openEditorModal(note) {
    editNoteForm.elements.id.value = String(note.id);
    editNoteForm.elements.title.value = note.title || "";
    editNoteForm.elements.content.value = note.content || "";
    editorModal.classList.remove("hidden");
    editorModal.setAttribute("aria-hidden", "false");
  }

  closeModalBtn.addEventListener("click", closeEditorModal);
  cancelModalBtn.addEventListener("click", closeEditorModal);
  editorModal.addEventListener("click", function (e) {
    if (e.target && e.target.getAttribute("data-close-modal") === "true") {
      closeEditorModal();
    }
  });

  editNoteForm.addEventListener("submit", async function (e) {
    e.preventDefault();
    const id = editNoteForm.elements.id.value;
    const title = editNoteForm.elements.title.value;
    const content = editNoteForm.elements.content.value;

    try {
      await api.updateNote(id, { title: title, content: content });
      closeEditorModal();
      await loadNotes();
      notify("Not güncellendi");
    } catch (err) {
      notify(err.message, true);
    }
  });

  function syncTopNav() {
    const protectedEls = document.querySelectorAll("[data-protected='true']");
    const publicOnlyEls = document.querySelectorAll("[data-public-only='true']");

    protectedEls.forEach(function (el) {
      el.style.display = isAuthed() ? "" : "none";
    });
    publicOnlyEls.forEach(function (el) {
      el.style.display = isAuthed() ? "none" : "";
    });
  }

  function cloneView(id) {
    return document.getElementById(id).content.cloneNode(true);
  }

  function renderHome() {
    app.innerHTML = "";
    app.appendChild(cloneView("home-view"));
  }

  function wireLogin() {
    app.innerHTML = "";
    app.appendChild(cloneView("login-view"));
    const form = document.getElementById("loginForm");
    form.addEventListener("submit", async function (e) {
      e.preventDefault();
      const fd = new FormData(form);
      try {
        const data = await api.login({
          email: fd.get("email"),
          password: fd.get("password")
        });
        const role = decodeRoleFromJWT(data.access_token);
        store.set({
          accessToken: data.access_token,
          refreshToken: data.refresh_token,
          role: role
        });
        notify("Başarıyla giriş yapıldı");
        location.hash = "#/dashboard";
      } catch (err) {
        notify(err.message, true);
      }
    });
  }

  function wireRegister() {
    app.innerHTML = "";
    app.appendChild(cloneView("register-view"));
    const form = document.getElementById("registerForm");
    form.addEventListener("submit", async function (e) {
      e.preventDefault();
      const fd = new FormData(form);
      try {
        await api.register({
          name: fd.get("name"),
          email: fd.get("email"),
          password: fd.get("password")
        });
        notify("Hesabınız başarıyla oluşturuldu. Lütfen giriş yapın.");
        location.hash = "#/login";
      } catch (err) {
        notify(err.message, true);
      }
    });
  }

  function renderNotes(notes) {
    const list = document.getElementById("notesList");
    list.innerHTML = "";
    if (!notes.length) {
      list.innerHTML = "<p class='muted'>No notes yet.</p>";
      return;
    }

    notes.forEach(function (note) {
      const node = document.getElementById("note-item").content.cloneNode(true);
      node.querySelector(".note-title").textContent = note.title;
      node.querySelector(".note-content").textContent = note.content;
      node.querySelector(".chip").textContent = "owner:" + note.user_id;

      const editBtn = node.querySelector(".edit-btn");
      const deleteBtn = node.querySelector(".delete-btn");

      editBtn.addEventListener("click", function () {
        openEditorModal(note);
      });

      deleteBtn.addEventListener("click", async function () {
        if (!confirm("Delete this note?")) return;
        try {
          await api.deleteNote(note.id);
          await loadNotes();
          notify("Notunuz silindi");
        } catch (err) {
          notify(err.message, true);
        }
      });

      list.appendChild(node);
    });
  }

  async function loadNotes() {
    const notes = await api.listNotes();
    store.set({ notes: notes });
    renderNotes(notes);

    const statsLine = document.getElementById("statsLine");
    if (statsLine) {
      statsLine.textContent = "Total notes visible to you: " + notes.length;
    }
  }

  function wireDashboard() {
    app.innerHTML = "";
    app.appendChild(cloneView("dashboard-view"));

    const panelTitle = document.getElementById("panelTitle");
    const sessionInfo = document.getElementById("sessionInfo");
    const role = store.state.role || "user";

    panelTitle.textContent = role === "admin" ? "Admin Panel" : "User Panel";
    sessionInfo.textContent = "Signed in as role: " + role;

    const createForm = document.getElementById("createNoteForm");
    const refreshBtn = document.getElementById("refreshNotesBtn");

    createForm.addEventListener("submit", async function (e) {
      e.preventDefault();
      const fd = new FormData(createForm);
      try {
        await api.createNote({
          title: fd.get("title"),
          content: fd.get("content")
        });
        createForm.reset();
        await loadNotes();
        notify("Not oluşturuldu");
      } catch (err) {
        notify(err.message, true);
      }
    });

    refreshBtn.addEventListener("click", async function () {
      try {
        await loadNotes();
        notify("Notlar yenilendi");
      } catch (err) {
        notify(err.message, true);
      }
    });

    loadNotes().catch(function (err) {
      notify(err.message, true);
    });
  }

  function routeGuard(path) {
    const publicPaths = ["/", "/login", "/register"];
    const authed = isAuthed();
    if (!authed && !publicPaths.includes(path)) return "/login";
    if (authed && (path === "/login" || path === "/register")) return "/dashboard";
    return path;
  }

  function parseRoute() {
    const hash = location.hash || "#/";
    const path = hash.replace(/^#/, "");
    return path || "/";
  }

  function renderRoute() {
    syncTopNav();
    const current = parseRoute();
    const next = routeGuard(current);
    if (next !== current) {
      location.hash = "#" + next;
      return;
    }

    if (next === "/") return renderHome();
    if (next === "/login") return wireLogin();
    if (next === "/register") return wireRegister();
    if (next === "/dashboard") return wireDashboard();

    app.innerHTML = "<section class='panel'><h2>404</h2><p class='muted'>Route not found.</p></section>";
  }

  window.addEventListener("hashchange", renderRoute);
  renderRoute();
})();
