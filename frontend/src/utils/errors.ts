const ERROR_MESSAGES: Record<string, string> = {
  "authentication required": "Требуется авторизация",
  "session expired. please login again.": "Сессия истекла. Войдите снова.",
  "invalid request body": "Некорректное тело запроса",
  "email already registered": "Email уже зарегистрирован",
  "failed to register user": "Не удалось зарегистрировать пользователя",
  "invalid email or password": "Неверный email или пароль",
  "failed to login user": "Не удалось выполнить вход",
  "failed to list restaurants": "Не удалось загрузить список ресторанов",
  "invalid restaurant id": "Некорректный ID ресторана",
  "restaurant not found": "Ресторан не найден",
  "product not found": "Продукт не найден",
  "product is not available": "Продукт недоступен",
  "insufficient product stock": "Недостаточно товара на складе",
  "failed to add item to cart": "Не удалось добавить товар в корзину",
  "failed to get cart": "Не удалось загрузить корзину",
  "cart not found": "Корзина не найдена",
  "failed to clear cart": "Не удалось очистить корзину",
  "invalid product id": "Некорректный ID продукта",
  "cart item not found": "Товар в корзине не найден",
  "failed to remove cart item": "Не удалось удалить товар из корзины",
  "active cart not found": "Активная корзина не найдена",
  "cart is empty": "Корзина пуста",
  "failed to create order": "Не удалось создать заказ",
  "invalid user identity": "Некорректная учётная запись",
  "invalid token": "Недействительный токен",
};

export function translateError(message: string): string {
  const normalized = message.trim();
  const translated = ERROR_MESSAGES[normalized.toLowerCase()];

  if (translated) {
    return translated;
  }

  const requestFailed = /^request failed \((\d+)\)$/i.exec(normalized);
  if (requestFailed) {
    return `Запрос не выполнен (${requestFailed[1]})`;
  }

  return normalized;
}

export function getErrorMessage(error: unknown, fallback: string): string {
  if (error instanceof Error) {
    return translateError(error.message);
  }

  return fallback;
}
