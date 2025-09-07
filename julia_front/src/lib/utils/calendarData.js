// 캘린더 데이터 로딩 유틸리티 모듈

import { loadExamPeriods, filterExamPeriodsForDate } from './examPeriod.js';

// 반 정보를 가져오는 함수
export async function loadClasses(fetchWithAuth, API_ENDPOINTS) {
	try {
		const response = await fetchWithAuth(`${API_ENDPOINTS.API_BASE_URL}/classes`);
		if (response.ok) {
			return await response.json();
		} else {
			console.error('Failed to load classes');
			return [];
		}
	} catch (error) {
		console.error('Error loading classes:', error);
		return [];
	}
}

// 시험 기간 데이터를 로드하는 함수
export async function loadExamPeriodsData(fetchWithAuth, API_ENDPOINTS) {
	try {
		return await loadExamPeriods(fetchWithAuth, API_ENDPOINTS);
	} catch (error) {
		console.error('Error loading exam periods:', error);
		return [];
	}
}

// 캘린더 데이터를 로드하는 함수
export async function loadCalendarData(fetchWithAuth, API_ENDPOINTS, currentYear, currentMonth, userRole, currentUser) {
	try {
		const yearMonth = `${currentYear}-${String(currentMonth + 1).padStart(2, '0')}`;

		// admin인 경우 모든 보강 일정을, 아닌 경우 자신의 보강 일정만 조회
		let res;
		if (userRole === 'admin') {
			res = await fetchWithAuth(API_ENDPOINTS.MAKEUPS_BY_MONTH(yearMonth));
		} else {
			// 학생인 경우 자신의 보강 일정만 조회
			const userId = currentUser?.id || '';
			res = await fetchWithAuth(API_ENDPOINTS.MAKEUPS_BY_USER(userId));
			// 월별 필터링은 클라이언트에서 처리
		}

		let makeupData = res.ok ? await res.json() : [];

		// 학생인 경우 현재 월의 데이터만 필터링
		if (userRole !== 'admin') {
			makeupData = makeupData.filter((s) => {
				const scheduleDate = new Date(s.makeup_date);
				return (
					scheduleDate.getFullYear() === currentYear && scheduleDate.getMonth() === currentMonth
				);
			});
		}

		return makeupData;
	} catch (error) {
		console.error('캘린더 데이터 로드 오류:', error);
		return [];
	}
}

// 선택된 날짜의 보강 일정을 로드하는 함수
export async function loadSelectedDateSchedules(fetchWithAuth, API_ENDPOINTS, date, userRole, currentUser) {
	try {
		const dateStr = toLocalDateString(date);

		// admin인 경우 모든 보강 일정을, 아닌 경우 자신의 보강 일정만 조회
		let res;
		if (userRole === 'admin') {
			res = await fetchWithAuth(API_ENDPOINTS.MAKEUPS_BY_DATE(dateStr));
		} else {
			// 학생인 경우 자신의 보강 일정만 조회
			const userId = currentUser?.id || '';
			res = await fetchWithAuth(API_ENDPOINTS.MAKEUPS_BY_USER_AND_DATE(userId, dateStr));
		}

		return res.ok ? await res.json() : [];
	} catch (error) {
		console.error('날짜별 보강 일정 로드 오류:', error);
		return [];
	}
}

// 캘린더를 생성하는 함수
export function generateCalendar(makeupData, examPeriods, currentYear, currentMonth) {
	const firstDay = new Date(currentYear, currentMonth, 1);
	const lastDay = new Date(currentYear, currentMonth + 1, 0);
	const startDate = new Date(firstDay);
	startDate.setDate(startDate.getDate() - firstDay.getDay()); // 일요일 시작

	const calendarDays = [];
	const today = new Date();

	// 이번 달의 마지막 날짜까지만 표시
	const endDate = new Date(lastDay);
	endDate.setDate(endDate.getDate() + (6 - lastDay.getDay())); // 토요일까지

	let currentDate = new Date(startDate);
	while (currentDate <= endDate) {
		const isCurrentMonth = currentDate.getMonth() === currentMonth;
		const isToday = currentDate.toDateString() === today.toDateString();

		const dateStr = toLocalDateString(currentDate);
		const daySchedules = (makeupData || []).filter((s) => s.makeup_date === dateStr);
		
		// 해당 날짜의 시험 기간 찾기
		const dayExamPeriods = filterExamPeriodsForDate(examPeriods, currentDate);

		calendarDays.push({
			date: new Date(currentDate),
			isToday,
			isCurrentMonth,
			makeupSchedules: daySchedules,
			examPeriods: dayExamPeriods,
		});

		currentDate.setDate(currentDate.getDate() + 1);
	}

	return calendarDays;
}

// 날짜를 YYYY-MM-DD 형식으로 변환하는 유틸리티 함수
export function toLocalDateString(date) {
	const year = date.getFullYear();
	const month = String(date.getMonth() + 1).padStart(2, '0');
	const day = String(date.getDate()).padStart(2, '0');
	return `${year}-${month}-${day}`;
}
