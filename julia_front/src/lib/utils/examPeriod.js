// 시험 기간 관리 유틸리티 모듈

// 시험 기간 색상 팔레트 (보강 색상과 구분되는 색상들)
export const examPeriodColors = [
	{ bg: '#fef3c7', border: '#fde68a', text: '#92400e' }, // 노란색
	{ bg: '#fce7f3', border: '#fbcfe8', text: '#be185d' }, // 핑크색
	{ bg: '#e0e7ff', border: '#c7d2fe', text: '#3730a3' }, // 보라색
	{ bg: '#d1fae5', border: '#a7f3d0', text: '#065f46' }, // 초록색
	{ bg: '#fef2f2', border: '#fecaca', text: '#991b1b' }, // 빨간색
	{ bg: '#f0fdf4', border: '#bbf7d0', text: '#166534' }, // 연두색
	{ bg: '#fefce8', border: '#fef08a', text: '#a16207' }, // 황금색
	{ bg: '#f0f9ff', border: '#bae6fd', text: '#0c4a6e' }, // 하늘색
	{ bg: '#fdf4ff', border: '#f3e8ff', text: '#7c2d12' }, // 라벤더색
	{ bg: '#f0fdfa', border: '#ccfbf1', text: '#134e4a' }, // 청록색
];

// 시험 기간 ID를 기반으로 색상 할당하는 함수
export function getExamPeriodColor(examPeriodId) {
	const colorIndex = examPeriodId % examPeriodColors.length;
	return examPeriodColors[colorIndex];
}

// 시험 기간에 반 이름을 추가하는 함수
export function getExamPeriodDisplayName(examPeriod, classes = []) {
	const classItem = classes.find(c => c.class_id === examPeriod.class_id);
	const className = classItem ? classItem.class_name : '알 수 없는 반';
	return `${className} - ${examPeriod.name}`;
}

// 시험 기간 데이터를 로드하는 함수
export async function loadExamPeriods(fetchWithAuth, API_ENDPOINTS) {
	try {
		const response = await fetchWithAuth(API_ENDPOINTS.EXAM_PERIODS);
		if (response.ok) {
			return await response.json();
		} else {
			console.error('Failed to load exam periods');
			return [];
		}
	} catch (error) {
		console.error('Error loading exam periods:', error);
		return [];
	}
}

// 특정 날짜의 시험 기간을 필터링하는 함수
export function filterExamPeriodsForDate(examPeriods, targetDate) {
	if (!examPeriods || !targetDate) return [];
	
	const targetDateStr = targetDate.toISOString().split('T')[0]; // YYYY-MM-DD 형식
	
	return examPeriods.filter((ep) => {
		// 시험 기간의 시작일과 종료일을 YYYY-MM-DD 형식으로 비교
		const startDate = ep.start_date;
		const endDate = ep.end_date;
		
		// 날짜 문자열 비교 (YYYY-MM-DD 형식)
		return targetDateStr >= startDate && targetDateStr <= endDate;
	});
}

// 시험 기간이 영어 시험인지 확인하는 함수
export function isEnglishExam(examPeriod, targetDate) {
	if (!examPeriod.english_date || !targetDate) return false;
	
	const targetDateStr = targetDate.toISOString().split('T')[0];
	return targetDateStr === examPeriod.english_date;
}
