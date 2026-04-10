export type ToastType = "success" | "error" | "info";

export interface Toast {
  id: number;
  message: string;
  type: ToastType;
}

let nextId = 0;
let toasts = $state<Toast[]>([]);

export function getToasts(): Toast[] {
  return toasts;
}

export function addToast(message: string, type: ToastType = "success") {
  const id = nextId++;
  toasts = [...toasts, { id, message, type }];
  setTimeout(() => dismissToast(id), 4000);
}

export function dismissToast(id: number) {
  toasts = toasts.filter((t) => t.id !== id);
}
