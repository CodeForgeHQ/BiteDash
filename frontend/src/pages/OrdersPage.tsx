export function OrdersPage() {
  return (
    <div className="page">
      <h1>Orders</h1>
      <div className="placeholder">
        <p>
          Orders page will be available after the backend implements{" "}
          <code>GET /orders</code>.
        </p>
        <p className="muted">
          You can place orders from the Cart page using{" "}
          <code>POST /orders/make</code>.
        </p>
      </div>
    </div>
  );
}
