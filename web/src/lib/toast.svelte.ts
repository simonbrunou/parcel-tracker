export type ToastType = "success" | "error" | "info";

export interface Toast {
  id: number;
  message: string;
  type: ToastType;
}

let nextId = 0;
let toasts = $state<Toast[]>([]);
const timeouts = new Map<number, ReturnType<typeof setTimeout>>();

export function getToasts(): Toast[] {
  return toasts;
}

export function addToast(message: string, type: ToastType = "success") {
  const id = nextId++;
  toasts = [...toasts, { id, message, type }];
  const timer = setTimeout(() => dismissToast(id), 4000);
  timeouts.set(id, timer);
}

export function dismissToast(id: number) {
  const timer = timeouts.get(id);
  if (timer) {
    clearTimeout(timer);
    timeouts.delete(id);
  }
  toasts = toasts.filter((t) => t.id !== id);
}
