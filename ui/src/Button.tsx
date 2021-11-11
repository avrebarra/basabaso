export function Button({ children }) {
  return (
    <button className="p-2 rounded bg-blue-200 hover:bg-blue-600 transition">
      {children}
    </button>
  );
}
