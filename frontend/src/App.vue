<template>
  <div id="app">
    <nav class="navbar">
      <div class="nav-container">
        <h1 class="logo">
          <span class="logo-icon">🎮</span>
          <span class="logo-text">ConsoleRent</span>
        </h1>
        <div class="nav-links">
          <router-link to="/">Главная</router-link>
          <router-link v-if="isLoggedIn" to="/rentals">Мои аренды</router-link>
          <button
            v-if="!isLoggedIn"
            @click="showAuthModal = true"
            class="auth-btn"
          >
            Войти / Регистрация
          </button>
          <div v-else class="user-menu">
            <span class="username">{{ currentUser?.username }}</span>
            <button @click="logout" class="auth-btn logout-btn">Выйти</button>
          </div>
        </div>
      </div>
    </nav>

    <router-view @auth-required="showAuthModal = true" />

    <!-- Модальное окно авторизации -->
    <div v-if="showAuthModal" class="modal" @click.self="showAuthModal = false">
      <div class="modal-content auth-modal">
        <button class="modal-close" @click="showAuthModal = false">✕</button>

        <div class="auth-tabs">
          <button
            :class="['tab-btn', { active: authMode === 'login' }]"
            @click="authMode = 'login'"
          >
            Вход
          </button>
          <button
            :class="['tab-btn', { active: authMode === 'register' }]"
            @click="authMode = 'register'"
          >
            Регистрация
          </button>
        </div>

        <!-- Форма входа -->
        <form v-if="authMode === 'login'" @submit.prevent="login">
          <input
            v-model="loginForm.email"
            type="email"
            placeholder="Email"
            required
            class="input"
          />
          <input
            v-model="loginForm.password"
            type="password"
            placeholder="Пароль"
            required
            class="input"
          />
          <button type="submit" class="btn-primary" :disabled="loading">
            {{ loading ? "Вход..." : "Войти" }}
          </button>
        </form>

        <!-- Форма регистрации -->
        <form v-if="authMode === 'register'" @submit.prevent="register">
          <input
            v-model="registerForm.username"
            type="text"
            placeholder="Имя пользователя"
            required
            class="input"
          />
          <input
            v-model="registerForm.email"
            type="email"
            placeholder="Email"
            required
            class="input"
          />
          <input
            v-model="registerForm.password"
            type="password"
            placeholder="Пароль (мин. 6 символов)"
            required
            class="input"
          />
          <button type="submit" class="btn-primary" :disabled="loading">
            {{ loading ? "Регистрация..." : "Зарегистрироваться" }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";
import axios from "axios";

export default {
  name: "App",
  setup() {
    const isLoggedIn = ref(false);
    const currentUser = ref(null);
    const showAuthModal = ref(false);
    const authMode = ref("login");
    const loading = ref(false);

    const loginForm = ref({
      email: "",
      password: "",
    });

    const registerForm = ref({
      username: "",
      email: "",
      password: "",
    });

    // Настройка axios интерсептора для добавления токена
    const setupAxiosInterceptor = () => {
      axios.interceptors.request.use(
        (config) => {
          const token = localStorage.getItem("auth_token");
          if (token) {
            config.headers.Authorization = `Bearer ${token}`;
          }
          return config;
        },
        (error) => {
          return Promise.reject(error);
        },
      );

      // Интерсептор для обработки 401 ошибок
      axios.interceptors.response.use(
        (response) => response,
        (error) => {
          if (error.response && error.response.status === 401) {
            console.log("Unauthorized, redirecting to login");
            logout();
            showAuthModal.value = true;
          }
          return Promise.reject(error);
        },
      );
    };

    // Сохранение токена в localStorage
    const saveToken = (token) => {
      localStorage.setItem("auth_token", token);
    };

    // Загрузка пользователя из токена
    const loadUser = () => {
      const token = localStorage.getItem("auth_token");
      if (token) {
        try {
          // Декодируем токен для получения информации о пользователе
          const payload = JSON.parse(atob(token.split(".")[1]));
          currentUser.value = {
            id: payload.user_id,
            username: payload.username,
            email: payload.email,
          };
          isLoggedIn.value = true;
        } catch (e) {
          console.error("Invalid token", e);
          logout();
        }
      }
    };

    const login = async () => {
      loading.value = true;
      try {
        const response = await axios.post("/api/login", loginForm.value);
        saveToken(response.data.token);
        currentUser.value = response.data.user;
        isLoggedIn.value = true;
        showAuthModal.value = false;
        loginForm.value = { email: "", password: "" };
        // Перезагружаем страницу для обновления состояния
        window.location.reload();
      } catch (error) {
        console.error("Login error:", error);
        alert(error.response?.data?.error || "Ошибка входа");
      } finally {
        loading.value = false;
      }
    };

    const register = async () => {
      loading.value = true;
      try {
        const response = await axios.post("/api/register", registerForm.value);
        saveToken(response.data.token);
        currentUser.value = response.data.user;
        isLoggedIn.value = true;
        showAuthModal.value = false;
        registerForm.value = { username: "", email: "", password: "" };
        // Перезагружаем страницу для обновления состояния
        window.location.reload();
      } catch (error) {
        console.error("Register error:", error);
        alert(error.response?.data?.error || "Ошибка регистрации");
      } finally {
        loading.value = false;
      }
    };

    const logout = () => {
      localStorage.removeItem("auth_token");
      isLoggedIn.value = false;
      currentUser.value = null;
      window.location.reload();
    };

    onMounted(() => {
      setupAxiosInterceptor();
      loadUser();
    });

    return {
      isLoggedIn,
      currentUser,
      showAuthModal,
      authMode,
      loading,
      loginForm,
      registerForm,
      login,
      register,
      logout,
    };
  },
};
</script>

<style scoped>
.navbar {
  background: white;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.5rem;
  font-weight: bold;
}

.logo-icon {
  font-size: 1.8rem;
}

.logo-text {
  background: linear-gradient(135deg, #0066cc 0%, #0052a0 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.nav-links {
  display: flex;
  gap: 2rem;
  align-items: center;
}

.nav-links a {
  text-decoration: none;
  color: #4a5568;
  font-weight: 500;
  transition: all 0.3s;
}

.nav-links a:hover {
  color: #0066cc;
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.username {
  color: #0066cc;
  font-weight: 500;
}

.auth-btn {
  background: linear-gradient(135deg, #0066cc 0%, #0052a0 100%);
  color: white;
  border: none;
  padding: 0.6rem 1.5rem;
  border-radius: 25px;
  cursor: pointer;
  transition: all 0.3s;
  font-weight: 500;
}

.auth-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 102, 204, 0.3);
}

.logout-btn {
  background: #e53e3e;
}

.logout-btn:hover {
  background: #c53030;
  box-shadow: 0 5px 15px rgba(229, 62, 62, 0.3);
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(5px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 20px;
  position: relative;
  max-width: 450px;
  width: 90%;
  padding: 2rem;
  animation: slideIn 0.3s ease;
}

.modal-close {
  position: absolute;
  top: 1rem;
  right: 1rem;
  background: #e2e8f0;
  border: none;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  cursor: pointer;
  font-size: 1.2rem;
  transition: all 0.3s;
}

.modal-close:hover {
  background: #cbd5e0;
}

.auth-tabs {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
  border-bottom: 2px solid #e2e8f0;
}

.tab-btn {
  flex: 1;
  padding: 0.75rem;
  background: none;
  border: none;
  font-size: 1rem;
  font-weight: 500;
  color: #718096;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
}

.tab-btn.active {
  color: #0066cc;
}

.tab-btn.active::after {
  content: "";
  position: absolute;
  bottom: -2px;
  left: 0;
  right: 0;
  height: 2px;
  background: #0066cc;
}

form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.input {
  padding: 0.75rem;
  border: 2px solid #e2e8f0;
  border-radius: 10px;
  font-size: 1rem;
  transition: all 0.3s;
}

.input:focus {
  outline: none;
  border-color: #0066cc;
  box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
}

.btn-primary {
  padding: 0.75rem;
  background: linear-gradient(135deg, #0066cc 0%, #0052a0 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 102, 204, 0.3);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

@keyframes slideIn {
  from {
    transform: translateY(-50px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}
</style>
