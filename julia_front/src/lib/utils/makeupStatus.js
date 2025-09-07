// 보강 상태별 색상 정의
export const makeupStatusColors = {
	pending: { bg: '#eef6ff', border: '#d6e7ff', text: '#1e40af' }, // 미정 - 예전 파란색
	present: { bg: '#dcfce7', border: '#bbf7d0', text: '#166534' }, // 출석 - 초록색
	absent: { bg: '#fef2f2', border: '#fecaca', text: '#991b1b' }, // 결석 - 빨간색
	postponed: { bg: '#fef3c7', border: '#fde68a', text: '#92400e' }, // 미룸 - 노란색
	exempt: { bg: '#f3f4f6', border: '#9ca3af', text: '#6b7280' } // 면제 - 회색 (더 진한 회색)
};

// 보강 상태별 색상 가져오기 함수
export function getMakeupStatusColor(status) {
	return makeupStatusColors[status] || makeupStatusColors.pending;
}

// 보강 상태 영어-한글 매핑
export const makeupStatusNames = {
	pending: '미정',
	present: '출석',
	absent: '결석',
	postponed: '미룸',
	exempt: '면제'
};

// 영어 상태를 한글로 변환하는 함수
export function getMakeupStatusName(status) {
	return makeupStatusNames[status] || '미정';
}

// 보강 상태 옵션 배열 (폼에서 사용)
export const makeupStatusOptions = [
	{ value: 'pending', label: '미정' },
	{ value: 'present', label: '출석' },
	{ value: 'absent', label: '결석' },
	{ value: 'postponed', label: '미룸' },
	{ value: 'exempt', label: '면제' }
];
