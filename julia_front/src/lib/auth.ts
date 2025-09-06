// 쿠키 관련 유틸리티 함수들
export function getCookie(name: string): string | null {
	const value = `; ${document.cookie}`;
	const parts = value.split(`; ${name}=`);
	if (parts.length === 2) return parts.pop()?.split(';').shift() || null;
	return null;
}

export function setCookie(name: string, value: string, days: number = 7): void {
	const expires = new Date();
	expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);
	document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/;SameSite=Strict`;
}

export function deleteCookie(name: string): void {
	document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 UTC;path=/;`;
}

// 세션 쿠키 이름 (백엔드와 일치해야 함)
export const SESSION_COOKIE_NAME = 'my-session';

// 인증 상태 확인
export function isAuthenticated(): boolean {
	const cookie = getCookie(SESSION_COOKIE_NAME);
	console.log('세션 쿠키 확인:', cookie);
	console.log('인증 상태:', cookie !== null);
	return cookie !== null;
}

// 로그아웃 처리
export function logout(): void {
	deleteCookie(SESSION_COOKIE_NAME);
}

// 401 응답 처리 함수
export function handleUnauthorized(): void {
	console.log('401 Unauthorized 응답 - 로그인 페이지로 리디렉트');
	logout();
	// 브라우저 환경에서만 window 객체 사용
	if (typeof window !== 'undefined') {
		window.location.href = '/';
	}
}

// API 요청 시 쿠키 포함
export async function fetchWithAuth(url: string, options: RequestInit = {}): Promise<Response> {
	const defaultOptions: RequestInit = {
		...options,
		credentials: 'include', // 쿠키 포함
		headers: {
			'Content-Type': 'application/json',
			...options.headers,
		},
	};

	const response = await fetch(url, defaultOptions);

	// 401 응답 처리
	if (response.status === 401) {
		handleUnauthorized();
	}

	return response;
}
