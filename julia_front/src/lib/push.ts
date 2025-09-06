// src/lib/push.ts (프론트)
import { VAPID_PUBLIC } from './config.js';

const VAPID_PUBLIC_KEY = VAPID_PUBLIC;

async function fetchVapidKey() {
  const res = await fetch('/api/push/vapid-public', { cache: 'no-store' });
  if (!res.ok) throw new Error('vapid-public fetch failed: ' + res.status);
  return (await res.text()).trim().replace(/^"+|"+$/g, '');
}

function urlBase64ToUint8Array(base64: string) {
  const padding = '='.repeat((4 - (base64.length % 4)) % 4);
  const base64Safe = (base64 + padding).replace(/-/g, '+').replace(/_/g, '/');
  const raw = atob(base64Safe);
  return Uint8Array.from([...raw].map(c => c.charCodeAt(0)));
}

export async function enablePush(userId?: string) {
  if (!('serviceWorker' in navigator) || !('PushManager' in window)) throw new Error('Push not supported');
  if (!userId) throw new Error('User ID is required');

  const reg = await navigator.serviceWorker.register('/sw.js', { scope: '/' });
  const perm = await Notification.requestPermission();
  if (perm !== 'granted') throw new Error('Notification permission denied');

  const vapid = await fetchVapidKey();
  const key = urlBase64ToUint8Array(vapid);
  if (key.length !== 65 || key[0] !== 0x04) throw new Error('Invalid VAPID key');

  const sub = await reg.pushManager.subscribe({
    userVisibleOnly: true,
    applicationServerKey: key,
  });

  const j = sub.toJSON(); // { endpoint, expirationTime, keys: { p256dh, auth } }

  const payload = {
    user_id: userId,
    endpoint: j.endpoint,
    keys: {
      p256dh: j.keys?.p256dh,
      auth:   j.keys?.auth,
    },
    // 필요하면 명시
    content_encoding: 'aes128gcm',
  };

  await fetch('/api/push/subscriptions', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
    credentials: 'include',
  });

  return sub;
}




export async function disablePush(userId?: string) {
  if (!('serviceWorker' in navigator) || !('PushManager' in window)) {
    throw new Error('Push not supported');
  }

  if (!userId) {
    throw new Error('User ID is required');
  }

  const reg = await navigator.serviceWorker.ready;
  const subscription = await reg.pushManager.getSubscription();
  
  if (!subscription) {
    throw new Error('No active subscription found');
  }

  console.log('Deleting subscription with endpoint:', subscription.endpoint);

  // 브라우저에서 구독 해제
  await subscription.unsubscribe();

  // 서버에서 특정 구독 삭제
  const response = await fetch(`/api/push/subscriptions/${userId}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ endpoint: subscription.endpoint }),
    credentials: 'include'
  });

  if (!response.ok) {
    const errorText = await response.text();
    console.error('Delete subscription failed:', response.status, errorText);
    throw new Error(`Failed to delete subscription: ${response.status} ${errorText}`);
  }

  console.log('Subscription deleted successfully');
  return true;
}

export async function getPushSubscription() {
  if (!('serviceWorker' in navigator) || !('PushManager' in window)) {
    return null;
  }

  const reg = await navigator.serviceWorker.ready;
  const subscription = await reg.pushManager.getSubscription();
  
  return subscription;
}

export async function getUserSubscriptions(userId: string) {
  try {
    const response = await fetch(`/api/push/subscriptions/${userId}`, {
      method: 'GET',
      credentials: 'include'
    });

    if (!response.ok) {
      // 500 에러는 실제 서버 에러, 404는 구독이 없는 경우
      if (response.status === 500) {
        console.error('Server error getting subscriptions:', response.status);
        throw new Error('Server error getting subscriptions');
      }
      // 기타 에러들도 처리
      console.error('Failed to get user subscriptions:', response.status);
      return [];
    }

    const data = await response.json();
    return data.data || [];
  } catch (error) {
    console.error('Error getting user subscriptions:', error);
    // 네트워크 에러 등은 빈 배열 반환 (구독이 없는 것으로 처리)
    return [];
  }
}
