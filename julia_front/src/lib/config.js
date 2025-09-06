// 환경 변수 설정
// Docker 환경에서는 프록시를 통해 API 호출
export const VAPID_PUBLIC = import.meta.env.VITE_VAPID_PUBLIC || '';
export const BACKEND_HOST = import.meta.env.VITE_BACKEND_HOST || 'http://localhost:8080';
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';
export const APP_TITLE = import.meta.env.VITE_APP_TITLE || 'Julia';
export const APP_VERSION = import.meta.env.VITE_APP_VERSION || '1.0.0';
export const DEBUG_MODE = import.meta.env.VITE_DEBUG_MODE === 'true';

// API 엔드포인트들
export const API_ENDPOINTS = {
	API_BASE_URL: API_BASE_URL,
	LOGIN: `${API_BASE_URL}/auth/login`,
	LOGOUT: `${API_BASE_URL}/auth/logout`,
	ME: `${API_BASE_URL}/auth/me`,
	MAKEUPS: `${API_BASE_URL}/makeups`,
	MAKEUPS_BY_MONTH: (yearMonth) => `${API_BASE_URL}/makeups/month/${yearMonth}`,
	MAKEUPS_BY_DATE: (date) => `${API_BASE_URL}/makeups/date/${date}`,
	MAKEUPS_BY_USER: (userId) => `${API_BASE_URL}/makeups/user/${userId}`,
	MAKEUPS_BY_USER_AND_DATE: (userId, date) => `${API_BASE_URL}/makeups/user/${userId}/date/${date}`,
	UPDATE_MAKEUP: (userId, date, time) =>
		`${API_BASE_URL}/makeups/user/${userId}/date/${date}/time/${time}`,
	DELETE_MAKEUP: (userId, date, time) =>
		`${API_BASE_URL}/makeups/user/${userId}/date/${date}/time/${time}`,
	// 시험 기간 관련 엔드포인트
	EXAM_PERIODS: `${API_BASE_URL}/exam-periods`,
	EXAM_PERIODS_BY_CLASS: (classId) => `${API_BASE_URL}/exam-periods/class/${classId}`,
	EXAM_PERIOD_BY_ID: (examPeriodId) => `${API_BASE_URL}/exam-periods/${examPeriodId}`,
	// 관리자 메모 관련 엔드포인트
	ADMIN_MEMO: `${API_BASE_URL}/admin/memo`,
};
