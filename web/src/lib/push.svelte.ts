import { getVAPIDKey, subscribePush, unsubscribePush } from "./api";

let pushSupported = $state(false);
let pushEnabled = $state(false);
let pushLoading = $state(false);
let registration: ServiceWorkerRegistration | null = null;

export function isPushSupported(): boolean {
  return pushSupported;
}

export function isPushEnabled(): boolean {
  return pushEnabled;
}

export function isPushLoading(): boolean {
  return pushLoading;
}

export async function initPush(): Promise<void> {
  if (!("serviceWorker" in navigator) || !("PushManager" in window)) {
    return;
  }
  pushSupported = true;

  try {
    registration = await navigator.serviceWorker.register("/sw.js");
    const existing = await registration.pushManager.getSubscription();
    pushEnabled = existing !== null;
  } catch (err) {
    console.error("Service worker registration failed:", err);
    pushSupported = false;
  }
}

export async function togglePush(): Promise<void> {
  if (!registration) return;
  pushLoading = true;

  try {
    if (pushEnabled) {
      const sub = await registration.pushManager.getSubscription();
      if (sub) {
        await unsubscribePush(sub.endpoint);
        await sub.unsubscribe();
      }
      pushEnabled = false;
    } else {
      const permission = await Notification.requestPermission();
      if (permission !== "granted") {
        return;
      }

      const { key } = await getVAPIDKey();
      const sub = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: urlBase64ToUint8Array(key),
      });

      const json = sub.toJSON();
      await subscribePush({
        endpoint: sub.endpoint,
        p256dh: json.keys!.p256dh!,
        auth: json.keys!.auth!,
      });
      pushEnabled = true;
    }
  } catch (err) {
    console.error("Push toggle failed:", err);
  } finally {
    pushLoading = false;
  }
}

function urlBase64ToUint8Array(base64String: string): Uint8Array {
  const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding)
    .replace(/-/g, "+")
    .replace(/_/g, "/");
  const rawData = atob(base64);
  const outputArray = new Uint8Array(rawData.length);
  for (let i = 0; i < rawData.length; i++) {
    outputArray[i] = rawData.charCodeAt(i);
  }
  return outputArray;
}
