import { createBrowserRouter, Navigate } from "react-router-dom";
import { Layout } from "../components/Layout";
import { CartPage } from "../pages/CartPage";
import { LoginPage } from "../pages/LoginPage";
import { OrdersPage } from "../pages/OrdersPage";
import { RegisterPage } from "../pages/RegisterPage";
import { RestaurantDetailsPage } from "../pages/RestaurantDetailsPage";
import { RestaurantsPage } from "../pages/RestaurantsPage";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Layout />,
    children: [
      { index: true, element: <Navigate to="/restaurants" replace /> },
      { path: "login", element: <LoginPage /> },
      { path: "register", element: <RegisterPage /> },
      { path: "restaurants", element: <RestaurantsPage /> },
      { path: "restaurants/:restaurantID", element: <RestaurantDetailsPage /> },
      { path: "cart", element: <CartPage /> },
      { path: "orders", element: <OrdersPage /> },
    ],
  },
]);
