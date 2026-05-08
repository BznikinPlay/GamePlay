import { createRouter, createWebHistory } from "vue-router";
import HomePage from "../views/HomePage.vue";
import MyRentals from "../views/MyRentals.vue";
import CheckoutPage from "../views/CheckoutPage.vue";
import ProfilePage from "../views/ProfilePage.vue";

const routes = [
  { path: "/", component: HomePage },
  { path: "/rentals", component: MyRentals },
  { path: "/checkout/:id?", component: CheckoutPage, name: "checkout" },
  { path: "/profile", component: ProfilePage, name: "profile" },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
