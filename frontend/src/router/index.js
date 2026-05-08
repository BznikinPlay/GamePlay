import { createRouter, createWebHistory } from "vue-router";
import HomePage from "../views/HomePage.vue";
import MyRentals from "../views/MyRentals.vue";
import CheckoutPage from "../views/CheckoutPage.vue";

const routes = [
  { path: "/", component: HomePage },
  { path: "/rentals", component: MyRentals },
  { path: "/checkout/:id?", component: CheckoutPage, name: "checkout" },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
