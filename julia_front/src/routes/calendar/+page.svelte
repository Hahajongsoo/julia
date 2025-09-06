<script>
	import { onMount } from 'svelte';
	import { isAuthenticated, fetchWithAuth } from '$lib/auth';
	import { goto } from '$app/navigation';
	import { API_ENDPOINTS } from '$lib/config';
	import { enablePush, disablePush, getPushSubscription, getUserSubscriptions } from '$lib/push';
	// SvelteKit props (required by framework, but not used in this component)

	let currentDate = new Date();
	let currentMonth = currentDate.getMonth();
	let currentYear = currentDate.getFullYear();
	let calendarDays = [];
	let isLoading = false;
	let viewMode = 'month'; // 'month' | 'agenda'

	// 모달 관련 상태
	let showModal = false;
	let selectedDate = null;
	let selectedDateSchedules = [];
	let isModalLoading = false;

	// 보강 생성 관련 상태
	let showCreateForm = false;
	let createFormData = {
		user_id: '',
		makeup_date: '',
		start_time: '',
		reason: '',
		status: 'pending',
	};
	let isCreating = false;

	// 보강 수정 관련 상태
	let showEditForm = false;
	let editFormData = {
		user_id: '',
		makeup_date: '',
		start_time: '',
		reason: '',
		status: 'pending',
	};
	let selectedMakeup = null;
	let isEditing = false;

	// 사용자 정보 관련 상태
	let currentUser = null;
	let userRole = 'student'; // 기본값

	// 시험 기간 관련 상태
	let examPeriods = [];
	let examPeriodsLoading = false;

	// 메모 관련 상태
	let showMemoPanel = false;
	let memoContent = '';
	let memoLastSaved = null;
	let autoSaveTimeout = null;

	// 보강 상태별 색상 정의
	const makeupStatusColors = {
		pending: { bg: '#eef6ff', border: '#d6e7ff', text: '#1e40af' }, // 미정 - 예전 파란색
		present: { bg: '#dcfce7', border: '#bbf7d0', text: '#166534' }, // 출석 - 초록색
		absent: { bg: '#fef2f2', border: '#fecaca', text: '#991b1b' }, // 결석 - 빨간색
		postponed: { bg: '#fef3c7', border: '#fde68a', text: '#92400e' }, // 미룸 - 노란색
		exempt: { bg: '#f3f4f6', border: '#9ca3af', text: '#6b7280' } // 면제 - 회색 (더 진한 회색)
	};

	// 보강 상태별 색상 가져오기 함수
	function getMakeupStatusColor(status) {
		return makeupStatusColors[status] || makeupStatusColors.pending;
	}

	// 보강 상태 영어-한글 매핑
	const makeupStatusNames = {
		pending: '미정',
		present: '출석',
		absent: '결석',
		postponed: '미룸',
		exempt: '면제'
	};

	// 영어 상태를 한글로 변환하는 함수
	function getMakeupStatusName(status) {
		return makeupStatusNames[status] || '미정';
	}

	// 캘린더 셀 마우스 휠 이벤트 처리
	function handleCellWheel(event) {
		const cell = event.currentTarget;
		
		// 스크롤 가능한 클래스가 있는지 확인
		if (!cell.classList.contains('scrollable')) {
			// 스크롤 가능하지 않으면 기본 동작 허용 (페이지 스크롤)
			return;
		}
		
		const { scrollTop, scrollHeight, clientHeight } = cell;
		const currentTime = Date.now();
		const state = scrollState.get(cell) || { lastScrollTime: 0, consecutiveScrolls: 0 };
		
		// 스크롤 가능한 상태인지 확인 (여유 공간을 더 크게 설정)
		const scrollBuffer = 15; // 스크롤 끝에서 15px 여유 공간
		const canScrollUp = scrollTop > scrollBuffer;
		const canScrollDown = scrollTop < scrollHeight - clientHeight - scrollBuffer;
		
		// 연속 스크롤 방지 (500ms 내에 3번 이상 스크롤하면 잠시 대기)
		if (currentTime - state.lastScrollTime < 500) {
			state.consecutiveScrolls++;
		} else {
			state.consecutiveScrolls = 0;
		}
		
		// 위로 스크롤 시도
		if (event.deltaY < 0) {
			if (canScrollUp) {
				// 셀 내부에서 위로 스크롤 가능
				event.preventDefault();
				cell.scrollTop -= 20;
				state.lastScrollTime = currentTime;
			} else if (scrollTop <= scrollBuffer && state.consecutiveScrolls < 3) {
				// 맨 위 근처에서는 연속 스크롤 제한
				event.preventDefault();
				state.lastScrollTime = currentTime;
			}
		}
		// 아래로 스크롤 시도
		else if (event.deltaY > 0) {
			if (canScrollDown) {
				// 셀 내부에서 아래로 스크롤 가능
				event.preventDefault();
				cell.scrollTop += 20;
				state.lastScrollTime = currentTime;
			} else if (scrollTop >= scrollHeight - clientHeight - scrollBuffer && state.consecutiveScrolls < 3) {
				// 맨 아래 근처에서는 연속 스크롤 제한
				event.preventDefault();
				state.lastScrollTime = currentTime;
			}
		}
		
		// 상태 업데이트
		scrollState.set(cell, state);
	}

	// 스크롤 상태 추적을 위한 변수
	let scrollState = new Map();

	// 셀에 마우스 오버 시 스크롤 가능 상태 확인
	function handleCellMouseEnter(event) {
		const cell = event.currentTarget;
		
		// 약간의 지연을 두고 스크롤 가능 여부 확인 (DOM 업데이트 대기)
		setTimeout(() => {
			const { scrollHeight, clientHeight } = cell;
			
			// 스크롤 가능한 내용이 있으면 클래스 추가
			if (scrollHeight > clientHeight + 2) {
				cell.classList.add('scrollable');
				// 스크롤 상태 초기화
				scrollState.set(cell, { lastScrollTime: 0, consecutiveScrolls: 0 });
			}
		}, 10);
	}

	// 셀에서 마우스 아웃 시 스크롤 가능 상태 클래스 제거
	function handleCellMouseLeave(event) {
		const cell = event.currentTarget;
		cell.classList.remove('scrollable');
		// 스크롤 상태 정리
		scrollState.delete(cell);
	}

	// 시험 기간 색상 팔레트 (보강 색상과 구분되는 색상들)
	const examPeriodColors = [
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
	function getExamPeriodColor(examPeriodId) {
		const colorIndex = examPeriodId % examPeriodColors.length;
		return examPeriodColors[colorIndex];
	}

	// 반 정보를 가져오는 함수
	let classes = [];
	async function loadClasses() {
		try {
			const response = await fetchWithAuth(`${API_ENDPOINTS.API_BASE_URL}/classes`);
			if (response.ok) {
				classes = await response.json();
			} else {
				console.error('Failed to load classes');
				classes = [];
			}
		} catch (error) {
			console.error('Error loading classes:', error);
			classes = [];
		}
	}

	// 시험 기간에 반 이름을 추가하는 함수
	function getExamPeriodDisplayName(examPeriod) {
		const classItem = classes.find(c => c.class_id === examPeriod.class_id);
		const className = classItem ? classItem.class_name : '알 수 없는 반';
		return `${className} - ${examPeriod.name}`;
	}

	// 푸시 구독 관련 상태
	let pushSubscription = null;
	let isPushEnabled = false;
	let isPushLoading = false;

	// 스와이프 관련 상태
	let touchStartX = 0;
	let touchStartY = 0;
	let touchEndX = 0;
	let touchEndY = 0;
	let isSwiping = false;

	const monthNames = [
		'1월',
		'2월',
		'3월',
		'4월',
		'5월',
		'6월',
		'7월',
		'8월',
		'9월',
		'10월',
		'11월',
		'12월',
	];
	const dayNames = ['일', '월', '화', '수', '목', '금', '토'];

	onMount(async () => {
		// --- Dynamic viewport height variable (mobile zoom / iOS UI changes safe) ---
		const setVh = () => {
			const h = window.visualViewport ? window.visualViewport.height : window.innerHeight;
			document.documentElement.style.setProperty('--vh', `${h}px`);
		};
		setVh();
		window.addEventListener('resize', setVh);
		window.visualViewport && window.visualViewport.addEventListener('resize', setVh);

		// if (!isAuthenticated()) {
		// 	goto('/');
		// 	return;
		// }
		await loadCurrentUser();
		await loadClasses();
		await loadCalendarData();
		await checkPushSubscription();
	});

	// 현재 사용자 정보 로드
	async function loadCurrentUser() {
		try {
			const res = await fetchWithAuth(API_ENDPOINTS.ME);
			if (res.ok) {
				currentUser = await res.json();
				userRole = currentUser.role || 'student';
				console.log('사용자 정보 로드 성공:', currentUser);
				console.log('현재 사용자 role:', userRole);
			} else {
				console.error('사용자 정보 로드 실패');
				userRole = 'student';
			}
		} catch (e) {
			console.error('사용자 정보 로드 오류:', e);
			userRole = 'student';
		}
	}

	// 시험 기간 로드
	async function loadExamPeriods() {
		examPeriodsLoading = true;
		try {
			const response = await fetchWithAuth(API_ENDPOINTS.EXAM_PERIODS);
			if (response.ok) {
				examPeriods = await response.json();
			} else {
				console.error('Failed to load exam periods');
				examPeriods = [];
			}
		} catch (error) {
			console.error('Error loading exam periods:', error);
			examPeriods = [];
		} finally {
			examPeriodsLoading = false;
		}
	}

	// 푸시 구독 상태 확인
	async function checkPushSubscription() {
		try {
			// 브라우저의 구독 상태 확인
			pushSubscription = await getPushSubscription();
			const browserHasSubscription = !!pushSubscription;
			
			// 서버의 구독 상태 확인 (사용자 ID가 있을 때만)
			let serverHasSubscription = false;
			if (currentUser?.id) {
				const serverSubscriptions = await getUserSubscriptions(currentUser.id);
				serverHasSubscription = serverSubscriptions.length > 0;
			}
			
			// 브라우저와 서버 모두에 구독이 있어야 활성화된 것으로 간주
			isPushEnabled = browserHasSubscription && serverHasSubscription;
			
			console.log('푸시 구독 상태:', {
				browser: browserHasSubscription,
				server: serverHasSubscription,
				final: isPushEnabled
			});
		} catch (e) {
			console.error('푸시 구독 상태 확인 오류:', e);
			isPushEnabled = false;
		}
	}

	// 푸시 구독 켜기
	async function handleEnablePush() {
		if (!currentUser?.id) {
			alert('로그인이 필요합니다.');
			return;
		}

		isPushLoading = true;
		try {
			await enablePush(currentUser.id);
			await checkPushSubscription();
			alert('푸시 알림이 활성화되었습니다.');
		} catch (e) {
			console.error('푸시 활성화 오류:', e);
			alert('푸시 알림 활성화에 실패했습니다: ' + e.message);
		} finally {
			isPushLoading = false;
		}
	}

	// 푸시 구독 끄기
	async function handleDisablePush() {
		if (!currentUser?.id) {
			alert('로그인이 필요합니다.');
			return;
		}

		isPushLoading = true;
		try {
			await disablePush(currentUser.id);
			await checkPushSubscription();
			alert('푸시 알림이 비활성화되었습니다.');
		} catch (e) {
			console.error('푸시 비활성화 오류:', e);
			alert('푸시 알림 비활성화에 실패했습니다: ' + e.message);
		} finally {
			isPushLoading = false;
		}
	}

	async function loadCalendarData() {
		isLoading = true;
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

			// 시험 기간 데이터도 함께 로드
			await loadExamPeriods();
			generateCalendar(makeupData);
		} catch (e) {
			console.error('캘린더 데이터 로드 오류:', e);
			generateCalendar([]);
		} finally {
			isLoading = false;
		}
	}

	function generateCalendar(makeupData) {
		const firstDay = new Date(currentYear, currentMonth, 1);
		const lastDay = new Date(currentYear, currentMonth + 1, 0);
		const startDate = new Date(firstDay);
		startDate.setDate(startDate.getDate() - firstDay.getDay()); // 일요일 시작

		calendarDays = [];
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
			const dayExamPeriods = (examPeriods || []).filter((ep) => {
				// 현재 날짜를 YYYY-MM-DD 형식으로 변환
				const currentDateStr = toLocalDateString(currentDate);
				
				// 시험 기간의 시작일과 종료일을 YYYY-MM-DD 형식으로 비교
				// 문자열 비교로 시간대 문제를 완전히 회피
				const isInPeriod = currentDateStr >= ep.start_date && currentDateStr <= ep.end_date;
				
				// 영어 시험 날짜와도 비교 (영어 시험 날짜가 있는 경우)
				const isEnglishDate = ep.english_date && currentDateStr === ep.english_date;
				
				return isInPeriod || isEnglishDate;
			});

			calendarDays.push({
				date: new Date(currentDate),
				isToday,
				isCurrentMonth,
				makeupSchedules: daySchedules,
				examPeriods: dayExamPeriods,
			});

			currentDate.setDate(currentDate.getDate() + 1);
		}
	}

	// 날짜 클릭 핸들러
	async function handleDateClick(date) {
		selectedDate = date;
		showModal = true;
		showCreateForm = false; // 모달 열 때 생성 폼 초기화
		await loadSelectedDateSchedules(date);
	}

	// 선택된 날짜의 보강 일정 로드
	async function loadSelectedDateSchedules(date) {
		isModalLoading = true;
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

			selectedDateSchedules = res.ok ? await res.json() : [];
		} catch (e) {
			console.error('날짜별 보강 일정 로드 오류:', e);
			selectedDateSchedules = [];
		} finally {
			isModalLoading = false;
		}
	}

	// 보강 생성 폼 열기
	function openCreateForm() {
		showCreateForm = true;
		// 선택된 날짜를 기본값으로 설정
		createFormData = {
			user_id: '',
			makeup_date: toLocalDateString(selectedDate),
			start_time: '',
			reason: '',
			status: 'pending',
		};
	}

	// 보강 생성 폼 닫기
	function closeCreateForm() {
		showCreateForm = false;
		createFormData = {
			user_id: '',
			makeup_date: '',
			start_time: '',
			reason: '',
			status: 'pending',
		};
	}

	// 보강 생성 요청
	async function createMakeup() {
		if (!createFormData.start_time || !createFormData.user_id) {
			alert('시간과 학생명을 입력해주세요.');
			return;
		}

		isCreating = true;
		try {
			const res = await fetchWithAuth(API_ENDPOINTS.MAKEUPS, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(createFormData),
			});

			if (res.ok) {
				alert('보강 일정이 성공적으로 생성되었습니다.');
				closeCreateForm();
				// 선택된 날짜의 보강 일정 다시 로드
				await loadSelectedDateSchedules(selectedDate);
				// 캘린더 데이터도 다시 로드
				await loadCalendarData();
			} else {
				const errorData = await res.json();
				alert(`보강 생성 실패: ${errorData.message || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('보강 생성 오류:', e);
			alert('보강 생성 중 오류가 발생했습니다.');
		} finally {
			isCreating = false;
		}
	}

	// 보강 수정 폼 열기
	function openEditForm(makeup) {
		selectedMakeup = makeup;
		showEditForm = true;
		editFormData = {
			user_id: makeup.user_id,
			makeup_date: makeup.makeup_date,
			start_time: makeup.start_time,
			reason: makeup.reason || '',
			status: makeup.status || 'pending',
		};
	}

	// 보강 수정 폼 닫기
	function closeEditForm() {
		showEditForm = false;
		selectedMakeup = null;
		editFormData = {
			user_id: '',
			makeup_date: '',
			start_time: '',
			reason: '',
			status: 'pending',
		};
	}

	// 보강 수정 요청
	async function updateMakeup() {
		if (!editFormData.start_time || !editFormData.user_id) {
			alert('시간과 학생명을 입력해주세요.');
			return;
		}

		isEditing = true;
		try {
			const res = await fetchWithAuth(
				API_ENDPOINTS.UPDATE_MAKEUP(
					selectedMakeup.user_id,
					selectedMakeup.makeup_date,
					selectedMakeup.start_time
				),
				{
					method: 'PUT',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify(editFormData),
				}
			);

			if (res.ok) {
				alert('보강 일정이 성공적으로 수정되었습니다.');
				closeEditForm();
				// 선택된 날짜의 보강 일정 다시 로드
				await loadSelectedDateSchedules(selectedDate);
				// 캘린더 데이터도 다시 로드
				await loadCalendarData();
			} else {
				const errorData = await res.json();
				alert(`보강 수정 실패: ${errorData.message || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('보강 수정 오류:', e);
			alert('보강 수정 중 오류가 발생했습니다.');
		} finally {
			isEditing = false;
		}
	}

	// 보강 삭제 요청
	async function deleteMakeup(makeup) {
		if (!confirm('정말로 이 보강 일정을 삭제하시겠습니까?')) {
			return;
		}

		try {
			const res = await fetchWithAuth(
				API_ENDPOINTS.DELETE_MAKEUP(makeup.user_id, makeup.makeup_date, makeup.start_time),
				{
					method: 'DELETE',
				}
			);

			if (res.ok) {
				alert('보강 일정이 성공적으로 삭제되었습니다.');
				// 선택된 날짜의 보강 일정 다시 로드
				await loadSelectedDateSchedules(selectedDate);
				// 캘린더 데이터도 다시 로드
				await loadCalendarData();
			} else {
				const errorData = await res.json();
				alert(`보강 삭제 실패: ${errorData.message || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('보강 삭제 오류:', e);
			alert('보강 삭제 중 오류가 발생했습니다.');
		}
	}

	// 보강 상태 변경 요청
	async function updateMakeupStatus(makeup, newStatus) {
		try {
			const updateData = {
				user_id: makeup.user_id,
				makeup_date: makeup.makeup_date,
				start_time: makeup.start_time,
				reason: makeup.reason,
				status: newStatus
			};

			const res = await fetchWithAuth(
				API_ENDPOINTS.UPDATE_MAKEUP(makeup.user_id, makeup.makeup_date, makeup.start_time),
				{
					method: 'PUT',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify(updateData),
				}
			);

			if (res.ok) {
				// 선택된 날짜의 보강 일정 다시 로드
				await loadSelectedDateSchedules(selectedDate);
				// 캘린더 데이터도 다시 로드
				await loadCalendarData();
			} else {
				const errorData = await res.json();
				alert(`상태 변경 실패: ${errorData.message || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('상태 변경 오류:', e);
			alert('상태 변경 중 오류가 발생했습니다.');
		}
	}

	// 모달 닫기
	function closeModal() {
		showModal = false;
		selectedDate = null;
		selectedDateSchedules = [];
		showCreateForm = false;
		showEditForm = false;
		selectedMakeup = null;
		createFormData = {
			user_id: '',
			makeup_date: '',
			start_time: '',
			reason: '',
			status: 'pending',
		};
		editFormData = {
			user_id: '',
			makeup_date: '',
			start_time: '',
			reason: '',
			status: 'pending',
		};
	}

	// 키보드 이벤트 처리 (ESC로 모달 닫기)
	function handleKeydown(event) {
		if (event.key === 'Escape' && showModal) {
			if (showCreateForm) {
				closeCreateForm();
			} else if (showEditForm) {
				closeEditForm();
			} else {
				closeModal();
			}
		}
	}

	// 터치 시작 이벤트
	function handleTouchStart(event) {
		touchStartX = event.touches[0].clientX;
		touchStartY = event.touches[0].clientY;
		isSwiping = false;
	}

	// 터치 이동 이벤트
	function handleTouchMove(event) {
		if (!touchStartX || !touchStartY) return;

		touchEndX = event.touches[0].clientX;
		touchEndY = event.touches[0].clientY;

		const deltaX = touchStartX - touchEndX;
		const deltaY = touchStartY - touchEndY;

		// 수평 스와이프가 수직 스와이프보다 클 때만 스와이프로 인식
		if (Math.abs(deltaX) > Math.abs(deltaY) && Math.abs(deltaX) > 50) {
			isSwiping = true;
		}
	}

	// 터치 종료 이벤트
	function handleTouchEnd(event) {
		if (!isSwiping || !touchStartX || !touchEndX) {
			touchStartX = 0;
			touchStartY = 0;
			touchEndX = 0;
			touchEndY = 0;
			return;
		}

		const deltaX = touchStartX - touchEndX;
		const minSwipeDistance = 50;

		if (Math.abs(deltaX) > minSwipeDistance) {
			if (deltaX > 0) {
				// 왼쪽으로 스와이프 - 다음 달
				nextMonth();
			} else {
				// 오른쪽으로 스와이프 - 이전 달
				previousMonth();
			}
		}

		// 터치 상태 초기화
		touchStartX = 0;
		touchStartY = 0;
		touchEndX = 0;
		touchEndY = 0;
		isSwiping = false;
	}

	// 로그아웃 핸들러
	async function handleLogout() {
		try {
			const res = await fetchWithAuth(API_ENDPOINTS.LOGOUT, {
				method: 'POST',
			});

			if (res.ok) {
				// 로컬 스토리지에서 인증 정보 제거
				localStorage.removeItem('authToken');
				// 로그인 페이지로 이동
				goto('/');
			} else {
				console.error('로그아웃 실패');
				// 강제로 로그인 페이지로 이동
				localStorage.removeItem('authToken');
				goto('/');
			}
		} catch (error) {
			console.error('로그아웃 오류:', error);
			// 오류가 발생해도 로그인 페이지로 이동
			localStorage.removeItem('authToken');
			goto('/');
		}
	}

	// 메모 관련 함수들
	function toggleMemoPanel() {
		showMemoPanel = !showMemoPanel;
		if (showMemoPanel) {
			loadMemo();
		}
	}

	function closeMemoPanel() {
		showMemoPanel = false;
		// 자동 저장 취소
		if (autoSaveTimeout) {
			clearTimeout(autoSaveTimeout);
		}
	}

	// 메모 로드
	async function loadMemo() {
		try {
			const response = await fetchWithAuth(API_ENDPOINTS.ADMIN_MEMO);
			if (response.ok) {
				const data = await response.json();
				memoContent = data.content || '';
				memoLastSaved = data.updated_at || null;
			} else {
				memoContent = '';
				memoLastSaved = null;
			}
		} catch (error) {
			console.error('메모 로드 오류:', error);
			memoContent = '';
			memoLastSaved = null;
		}
	}

	// 메모 내용 변경 핸들러 (자동 저장)
	function handleMemoChange() {
		// 기존 자동 저장 타이머 취소
		if (autoSaveTimeout) {
			clearTimeout(autoSaveTimeout);
		}
		
		// 3초 후 자동 저장
		autoSaveTimeout = setTimeout(() => {
			if (memoContent.trim()) {
				saveMemo(true); // 자동 저장 플래그
			}
		}, 3000);
	}

	// 메모 저장
	async function saveMemo(isAutoSave = false) {
		try {
			const response = await fetchWithAuth(API_ENDPOINTS.ADMIN_MEMO, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ content: memoContent })
			});
			
			if (response.ok) {
				memoLastSaved = new Date().toISOString();
				if (!isAutoSave) {
					// 수동 저장일 때만 알림
					showToast('메모가 저장되었습니다.', 'success');
				}
			}
		} catch (error) {
			console.error('메모 저장 오류:', error);
			if (!isAutoSave) {
				showToast('메모 저장에 실패했습니다.', 'error');
			}
		}
	}

	// 메모 전체 삭제
	async function clearMemo() {
		if (!confirm('정말로 메모를 모두 삭제하시겠습니까?')) return;
		
		try {
			const response = await fetchWithAuth(API_ENDPOINTS.ADMIN_MEMO, {
				method: 'DELETE'
			});
			
			if (response.ok) {
				memoContent = '';
				memoLastSaved = null;
				showToast('메모가 삭제되었습니다.', 'success');
			}
		} catch (error) {
			console.error('메모 삭제 오류:', error);
			showToast('메모 삭제에 실패했습니다.', 'error');
		}
	}

	// 마지막 저장 시간 포맷팅
	function formatLastSaved(dateString) {
		const date = new Date(dateString);
		const now = new Date();
		const diff = now - date;
		
		if (diff < 60000) { // 1분 미만
			return '방금 전';
		} else if (diff < 3600000) { // 1시간 미만
			return `${Math.floor(diff / 60000)}분 전`;
		} else if (diff < 86400000) { // 1일 미만
			return `${Math.floor(diff / 3600000)}시간 전`;
		} else {
			return date.toLocaleDateString('ko-KR', {
				month: 'short',
				day: 'numeric',
				hour: '2-digit',
				minute: '2-digit'
			});
		}
	}

	// 토스트 알림 함수
	function showToast(message, type = 'info') {
		// 간단한 토스트 알림 구현
		const toast = document.createElement('div');
		toast.className = `toast toast-${type}`;
		toast.textContent = message;
		toast.style.cssText = `
			position: fixed;
			top: 20px;
			right: 20px;
			background: ${type === 'success' ? '#10b981' : type === 'error' ? '#ef4444' : '#3b82f6'};
			color: white;
			padding: 12px 20px;
			border-radius: 8px;
			z-index: 10000;
			font-size: 14px;
			font-weight: 500;
			box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
			animation: slideInRight 0.3s ease-out;
		`;
		
		document.body.appendChild(toast);
		
		setTimeout(() => {
			toast.style.animation = 'slideOutRight 0.3s ease-in forwards';
			setTimeout(() => {
				document.body.removeChild(toast);
			}, 300);
		}, 3000);
	}

	function previousMonth() {
		currentMonth--;
		if (currentMonth < 0) {
			currentMonth = 11;
			currentYear--;
		}
		loadCalendarData();
	}

	function nextMonth() {
		currentMonth++;
		if (currentMonth > 11) {
			currentMonth = 0;
			currentYear++;
		}
		loadCalendarData();
	}

	function goToday() {
		const now = new Date();
		currentYear = now.getFullYear();
		currentMonth = now.getMonth();
		loadCalendarData();
	}

	// month input: value는 계산한 문자열을 단방향 주입하고, 변경은 on:change에서 처리
	function monthInputValue() {
		return `${currentYear}-${String(currentMonth + 1).padStart(2, '0')}`;
	}
	function onMonthInputChange(e) {
		const val = e.target.value; // yyyy-MM
		if (!val) return;
		const [y, m] = val.split('-').map(Number);
		currentYear = y;
		currentMonth = m - 1;
		loadCalendarData();
	}

	function dayClasses({ date, isCurrentMonth, isToday }) {
		const w = date.getDay(); // 0=일 ... 3=수, 6=토
		return [
			'cal-day',
			!isCurrentMonth && 'is-out',
			isToday && 'is-today',
			w === 0 && 'sun',
			w === 6 && 'sat',
			w === 3 && 'wed', // 스타일링 포인트용
		]
			.filter(Boolean)
			.join(' ');
	}

	function fmtTime(t) {
		if (!t) return '';
		return t.slice(0, 5); // HH:MM
	}

	// 로컬 시간대를 고려한 날짜 문자열 변환 함수
	function toLocalDateString(date) {
		const year = date.getFullYear();
		const month = String(date.getMonth() + 1).padStart(2, '0');
		const day = String(date.getDate()).padStart(2, '0');
		return `${year}-${month}-${day}`;
	}

	// 날짜 포맷팅 함수
	function formatDate(date) {
		const year = date.getFullYear();
		const month = date.getMonth() + 1;
		const day = date.getDate();
		const dayName = dayNames[date.getDay()];
		return `${year}년 ${month}월 ${day}일 (${dayName})`;
	}

	// Agenda: 이번 달의 일정만 날짜별로 뭉치기
	$: agendaItems = calendarDays
		.filter((d) => d.isCurrentMonth && (d.makeupSchedules?.length || d.examPeriods?.length))
		.map((d) => ({
			key: toLocalDateString(d.date),
			label: `${d.date.getMonth() + 1}/${d.date.getDate()} (${dayNames[d.date.getDay()]})`,
			makeupItems: d.makeupSchedules || [],
			examPeriodItems: d.examPeriods || [],
		}));
</script>

<svelte:head>
	<title>Julia - 월별 캘린더</title>
	<meta name="description" content="Julia 시스템 월별 보강 캘린더" />
	<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover" />
</svelte:head>

<svelte:window on:keydown={handleKeydown} />

<div class="wrap">
	<header class="topbar">
		<div class="brand">
			<h1 aria-label="Julia 보강 캘린더">
				<span class="logo-dot" aria-hidden="true" /> Julia 보강 캘린더
			</h1>

			<div class="month-line">
				<strong class="month-text">{currentYear}년 {monthNames[currentMonth]}</strong>
				<label class="month-input-wrapper" for="month-input">
					<input
						id="month-input"
						class="month-input"
						type="month"
						value={monthInputValue()}
						on:change={onMonthInputChange}
						aria-label="월 선택"
					/>
					<svg class="calendar-icon" width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
						<path
							d="M8 2v4M16 2v4M3 10h18M5 4h14a2 2 0 012 2v14a2 2 0 01-2 2H5a2 2 0 01-2-2V6a2 2 0 012-2z"
							stroke="currentColor"
							fill="none"
							stroke-width="2"
						/>
					</svg>
				</label>
			</div>
			{#if currentUser}
				<div class="user-info">
					<span class="user-role {userRole}">{userRole === 'admin' ? '관리자' : '학생'}</span>
					<span class="user-name">{currentUser.id || '사용자'}</span>
				</div>
			{/if}
		</div>

		<div class="actions">
			<div class="actions-top">
				{#if userRole === 'admin'}
					<button
						class="btn ghost dashboard"
						on:click={() => goto('/dashboard')}
						aria-label="관리자 페이지"
					>
						<svg
							class="dashboard-icon"
							width="16"
							height="16"
							viewBox="0 0 24 24"
							aria-hidden="true"
						>
							<path
								d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2M12 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8z"
								stroke="currentColor"
								fill="none"
								stroke-width="2"
							/>
						</svg>
						관리 페이지
					</button>
				{/if}
				{#if userRole !== 'admin'}
					{#if isPushEnabled}
						<button 
							class="btn ghost push-toggle" 
							on:click={handleDisablePush}
							disabled={isPushLoading}
							aria-label="푸시 알림 끄기"
						>
							<svg width="14" height="14" viewBox="0 0 24 24" aria-hidden="true">
								<path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" stroke="currentColor" fill="none" stroke-width="2"/>
								<path d="M13.73 21a2 2 0 0 1-3.46 0" stroke="currentColor" fill="none" stroke-width="2"/>
							</svg>
							{isPushLoading ? '끄는 중...' : '푸시 끄기'}
						</button>
					{:else}
						<button 
							class="btn ghost push-toggle" 
							on:click={handleEnablePush}
							disabled={isPushLoading}
							aria-label="푸시 알림 켜기"
						>
							<svg width="14" height="14" viewBox="0 0 24 24" aria-hidden="true">
								<path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" stroke="currentColor" fill="none" stroke-width="2"/>
								<path d="M13.73 21a2 2 0 0 1-3.46 0" stroke="currentColor" fill="none" stroke-width="2"/>
							</svg>
							{isPushLoading ? '켜는 중...' : '푸시 켜기'}
						</button>
					{/if}
				{/if}
				<button class="btn ghost logout" on:click={handleLogout} aria-label="로그아웃">
					<svg class="logout-icon" width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
						<path
							d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4M16 17l5-5-5-5M21 12H9"
							stroke="currentColor"
							fill="none"
							stroke-width="2"
						/>
					</svg>
					로그아웃
				</button>
			</div>
			<div class="actions-bottom">
				<button class="btn ghost" on:click={goToday} aria-label="오늘로 이동">오늘</button>
				
				<!-- 메모 아이콘 추가 -->
				{#if userRole === 'admin'}
					<button 
						class="btn ghost memo-toggle" 
						on:click={toggleMemoPanel}
						aria-label="메모장 열기"
						title="메모장"
					>
						<svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true">
							<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" fill="none" stroke-width="2"/>
							<polyline points="14,2 14,8 20,8" stroke="currentColor" fill="none" stroke-width="2"/>
							<line x1="16" y1="13" x2="8" y2="13" stroke="currentColor" stroke-width="2"/>
							<line x1="16" y1="17" x2="8" y2="17" stroke="currentColor" stroke-width="2"/>
							<polyline points="10,9 9,9 8,9" stroke="currentColor" fill="none" stroke-width="2"/>
						</svg>
						메모
					</button>
				{/if}
				
				<div class="seg">
					<button
						class="btn icon"
						on:click={previousMonth}
						aria-label="이전 달"
						disabled={isLoading}
					>
						<svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true"
							><path d="M15 18l-6-6 6-6" stroke="currentColor" fill="none" stroke-width="2" /></svg
						>
					</button>
					<button class="btn icon" on:click={nextMonth} aria-label="다음 달" disabled={isLoading}>
						<svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true"
							><path d="M9 6l6 6-6 6" stroke="currentColor" fill="none" stroke-width="2" /></svg
						>
					</button>
				</div>
				<div class="seg view">
					<button
						class="btn {viewMode === 'month' ? 'primary' : 'ghost'}"
						on:click={() => (viewMode = 'month')}
						aria-pressed={viewMode === 'month'}>월</button
					>
					<button
						class="btn {viewMode === 'agenda' ? 'primary' : 'ghost'}"
						on:click={() => (viewMode = 'agenda')}
						aria-pressed={viewMode === 'agenda'}>목록</button
					>
				</div>
			</div>
		</div>
	</header>

	{#if isLoading}
		<div class="card">
			<div class="skeleton head" />
			<div class="skeleton row" />
			<div class="skeleton row" />
			<div class="skeleton row" />
		</div>
	{:else if viewMode === 'month'}
		<section
			class="card calendar"
			aria-label="월별 캘린더"
			on:touchstart={handleTouchStart}
			on:touchmove={handleTouchMove}
			on:touchend={handleTouchEnd}
		>
			<div class="weekdays">
				{#each dayNames as d, i}
					<div class="wday {i === 0 ? 'sun' : ''} {i === 6 ? 'sat' : ''} {i === 3 ? 'wed' : ''}">
						{d}
					</div>
				{/each}
			</div>

			<div class="grid">
				{#each calendarDays as day}
					<div
						class={dayClasses(day)}
						tabindex="0"
						aria-label={`${day.date.getMonth() + 1}월 ${day.date.getDate()}일`}
						on:click={() => handleDateClick(day.date)}
						on:keydown={(e) => e.key === 'Enter' && handleDateClick(day.date)}
						on:wheel={handleCellWheel}
						on:mouseenter={handleCellMouseEnter}
						on:mouseleave={handleCellMouseLeave}
						role="button"
					>
						<div class="num-wrap">
							<span class="num {day.isToday ? 'today' : ''}">{day.date.getDate()}</span>
						</div>

						{#if day.examPeriods.length > 0}
							<div class="exam-periods">
								{#each day.examPeriods as ep}
									{@const color = getExamPeriodColor(ep.exam_period_id)}
									{@const displayName = getExamPeriodDisplayName(ep)}
									{@const currentDateStr = toLocalDateString(day.date)}
									{@const isEnglishDate = ep.english_date && currentDateStr === ep.english_date}
									<div 
										class="exam-pill {isEnglishDate ? 'english-exam' : ''}" 
										title={displayName + ' (' + ep.start_date + ' ~ ' + ep.end_date + ')' + (isEnglishDate ? ' - 영어시험' : '')}
										style="background-color: {color.bg}; border-color: {color.border};"
									>
										<span class="exam-pill-name" style="color: {color.text};">
											{isEnglishDate ? '📚 ' + displayName : displayName}
										</span>
										{#if isEnglishDate}
											<span class="english-indicator">EN</span>
										{/if}
									</div>
								{/each}
							</div>
						{/if}
						
						{#if day.makeupSchedules.length > 0}
							<div class="list">
								{#each day.makeupSchedules as s}
									{@const statusColor = getMakeupStatusColor(s.status)}
									<div 
										class="pill" 
										title={(s.user_id || '학생') + ' ' + fmtTime(s.start_time) + ' (' + getMakeupStatusName(s.status) + ')'}
										style="background-color: {statusColor.bg}; border-color: {statusColor.border};"
									>
										<span class="pill-name" style="color: {statusColor.text};">{s.user_id || '학생'}</span>
										<span class="pill-time" style="color: {statusColor.text};">{fmtTime(s.start_time)}</span>
									</div>
								{/each}
							</div>
						{/if}
					</div>
				{/each}
			</div>
		</section>
	{:else}
		<section class="card agenda" aria-label="목록형 보강 일정">
			{#if agendaItems.length === 0}
				<p class="empty">이번 달 등록된 보강 일정과 시험 기간이 없습니다.</p>
			{:else}
				{#each agendaItems as g}
					<div class="ag-day">
						<div class="ag-head">{g.label}</div>
						
						<!-- 시험 기간 섹션 -->
						{#if g.examPeriodItems.length > 0}
							<div class="ag-section">
								<div class="ag-section-title">시험 기간</div>
								<ul class="ag-list">
									{#each g.examPeriodItems as ep}
										{@const color = getExamPeriodColor(ep.exam_period_id)}
										{@const displayName = getExamPeriodDisplayName(ep)}
										{@const currentDateStr = g.key}
										{@const isEnglishDate = ep.english_date && currentDateStr === ep.english_date}
										<li 
											class="ag-item exam-period-item {isEnglishDate ? 'english-exam' : ''}"
											style="border-left-color: {color.text}; background: color-mix(in srgb, {color.bg} 20%, var(--card) 80%);"
										>
											<div class="ag-main">
												<strong class="ag-name" style="color: {color.text};">
													{isEnglishDate ? '📚 ' + displayName : displayName}
													{#if isEnglishDate}
														<span class="english-badge">EN</span>
													{/if}
												</strong>
												<span class="ag-time">
													{ep.start_date} ~ {ep.end_date}
													{#if isEnglishDate}
														<span class="english-label">영어시험</span>
													{/if}
												</span>
											</div>
											{#if ep.description}
												<div class="ag-reason">{ep.description}</div>
											{/if}
										</li>
									{/each}
								</ul>
							</div>
						{/if}
						
						<!-- 보강 일정 섹션 -->
						{#if g.makeupItems.length > 0}
							<div class="ag-section">
								<div class="ag-section-title">보강 일정</div>
								<ul class="ag-list">
									{#each g.makeupItems as s}
										{@const statusColor = getMakeupStatusColor(s.status)}
										<li class="ag-item" style="border-left-color: {statusColor.border};">
											<div class="ag-main">
												<strong class="ag-name" style="color: {statusColor.text};">{s.user_id || '학생'}</strong>
												<span class="ag-time" style="color: {statusColor.text};">{fmtTime(s.start_time)}</span>
												<span class="ag-status" style="color: {statusColor.text}; background-color: {statusColor.bg}; border-color: {statusColor.border};">
													{getMakeupStatusName(s.status)}
												</span>
											</div>
											{#if s.reason}
												<div class="ag-reason">{s.reason}</div>
											{/if}
										</li>
									{/each}
								</ul>
							</div>
						{/if}
					</div>
				{/each}
			{/if}
		</section>
	{/if}

	{#if showModal}
		<div
			class="modal-overlay"
			on:click={closeModal}
			on:keydown={(e) => e.key === 'Escape' && closeModal()}
			role="button"
			tabindex="0"
		/>
		<div class="modal-content" on:click|stopPropagation>
			<div class="modal-header">
				<h2>{formatDate(selectedDate)}</h2>
				<button class="btn-close" on:click={closeModal} aria-label="닫기">
					<svg width="20" height="20" viewBox="0 0 24 24" aria-hidden="true">
						<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" fill="none" stroke-width="2" />
					</svg>
				</button>
			</div>

			{#if showCreateForm}
				<!-- 보강 생성 폼 -->
				<div class="create-form">
					<h3>새 보강 일정 생성</h3>
					<form on:submit|preventDefault={createMakeup}>
						<div class="form-group">
							<label for="create-date">날짜</label>
							<input
								id="create-date"
								type="date"
								bind:value={createFormData.makeup_date}
								required
								class="form-input"
							/>
						</div>

						<div class="form-group">
							<label for="create-time">시간</label>
							<input
								id="create-time"
								type="time"
								bind:value={createFormData.start_time}
								required
								class="form-input"
							/>
						</div>

						<div class="form-group">
							<label for="create-user">학생명</label>
							<input
								id="create-user"
								type="text"
								bind:value={createFormData.user_id}
								placeholder="학생 이름을 입력하세요"
								required
								class="form-input"
							/>
						</div>

						<div class="form-group">
							<label for="create-reason">사유 (선택사항)</label>
							<textarea
								id="create-reason"
								bind:value={createFormData.reason}
								placeholder="보강 사유를 입력하세요"
								class="form-textarea"
								rows="3"
							/>
						</div>

						<div class="form-group">
							<label for="create-status">상태</label>
							<select
								id="create-status"
								bind:value={createFormData.status}
								required
								class="form-input"
							>
								<option value="pending">미정</option>
								<option value="present">출석</option>
								<option value="absent">결석</option>
								<option value="postponed">미룸</option>
								<option value="exempt">면제</option>
							</select>
						</div>

						<div class="form-actions">
							<button
								type="button"
								class="btn ghost"
								on:click={closeCreateForm}
								disabled={isCreating}
							>
								취소
							</button>
							<button type="submit" class="btn primary" disabled={isCreating}>
								{isCreating ? '생성 중...' : '보강 생성'}
							</button>
						</div>
					</form>
				</div>
			{:else if showEditForm}
				<!-- 보강 수정 폼 -->
				<div class="create-form">
					<h3>보강 일정 수정</h3>
					<form on:submit|preventDefault={updateMakeup}>
						<div class="form-group">
							<label for="edit-date">날짜</label>
							<input
								id="edit-date"
								type="date"
								bind:value={editFormData.makeup_date}
								required
								class="form-input"
							/>
						</div>

						<div class="form-group">
							<label for="edit-time">시간</label>
							<input
								id="edit-time"
								type="time"
								bind:value={editFormData.start_time}
								required
								class="form-input"
							/>
						</div>

						<div class="form-group">
							<label for="edit-user">학생명</label>
							<input
								id="edit-user"
								type="text"
								bind:value={editFormData.user_id}
								placeholder="학생 이름을 입력하세요"
								required
								class="form-input"
							/>
						</div>

						<div class="form-group">
							<label for="edit-reason">사유 (선택사항)</label>
							<textarea
								id="edit-reason"
								bind:value={editFormData.reason}
								placeholder="보강 사유를 입력하세요"
								class="form-textarea"
								rows="3"
							/>
						</div>

						<div class="form-group">
							<label for="edit-status">상태</label>
							<select
								id="edit-status"
								bind:value={editFormData.status}
								required
								class="form-input"
							>
								<option value="pending">미정</option>
								<option value="present">출석</option>
								<option value="absent">결석</option>
								<option value="postponed">미룸</option>
								<option value="exempt">면제</option>
							</select>
						</div>

						<div class="form-actions">
							<button type="button" class="btn ghost" on:click={closeEditForm} disabled={isEditing}>
								취소
							</button>
							<button type="submit" class="btn primary" disabled={isEditing}>
								{isEditing ? '수정 중...' : '보강 수정'}
							</button>
						</div>
					</form>
				</div>
			{:else}
				<!-- 기존 보강 일정 목록 -->
				{#if isModalLoading}
					<div class="modal-loading">
						<div class="skeleton head" />
						<div class="skeleton row" />
						<div class="skeleton row" />
					</div>
				{:else if selectedDateSchedules.length === 0}
					<!-- Debug: userRole = {userRole} -->
					{#if userRole === 'admin'}
						<div
							class="modal-empty"
							on:click={openCreateForm}
							on:keydown={(e) => e.key === 'Enter' && openCreateForm()}
							role="button"
							tabindex="0"
						>
							<svg width="48" height="48" viewBox="0 0 24 24" aria-hidden="true">
								<path d="M8 12h8M12 8v8" stroke="currentColor" fill="none" stroke-width="2" />
							</svg>
							<p>선택된 날짜에 보강 일정이 없습니다.</p>
							<p class="click-hint">클릭하여 새 보강 일정을 생성하세요</p>
						</div>
					{:else}
						<div class="modal-empty">
							<svg width="48" height="48" viewBox="0 0 24 24" aria-hidden="true">
								<path d="M8 12h8M12 8v8" stroke="currentColor" fill="none" stroke-width="2" />
							</svg>
							<p>선택된 날짜에 보강 일정이 없습니다.</p>
						</div>
					{/if}
				{:else}
					<div class="modal-header-actions">
						<div class="schedule-summary">
							<span class="schedule-count">총 {selectedDateSchedules.length}개의 보강 일정</span>
						</div>
						{#if userRole === 'admin'}
							<button class="btn primary create-btn compact" on:click={openCreateForm}>
								<svg width="14" height="14" viewBox="0 0 24 24" aria-hidden="true">
									<path d="M12 5v14M5 12h14" stroke="currentColor" fill="none" stroke-width="2" />
								</svg>
								추가
							</button>
						{/if}
					</div>
					<ul class="modal-list">
						{#each selectedDateSchedules as s, index}
							{@const statusColor = getMakeupStatusColor(s.status)}
							<li class="modal-item" style="border-left-color: {statusColor.border};">
								<div class="modal-main">
									<div class="modal-info">
										<strong class="modal-name" style="color: {statusColor.text};">{s.user_id || '학생'}</strong>
										<span class="modal-index">#{index + 1}</span>
										<span class="modal-status" style="color: {statusColor.text}; background-color: {statusColor.bg}; border-color: {statusColor.border};">
											{getMakeupStatusName(s.status)}
										</span>
									</div>
									<div class="modal-actions-item">
										<span class="modal-time" style="color: {statusColor.text};">{fmtTime(s.start_time)}</span>
										{#if userRole === 'admin'}
											<div class="status-change-container">
												<select
													class="status-select"
													value={s.status || 'pending'}
													on:change={(e) => updateMakeupStatus(s, e.target.value)}
													title="상태 변경"
												>
													<option value="pending">미정</option>
													<option value="present">출석</option>
													<option value="absent">결석</option>
													<option value="postponed">미룸</option>
													<option value="exempt">면제</option>
												</select>
											</div>
											<div class="item-actions">
												<button
													class="btn-icon edit-btn"
													on:click={() => openEditForm(s)}
													title="수정"
												>
													<svg width="14" height="14" viewBox="0 0 24 24" aria-hidden="true">
														<path
															d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"
															stroke="currentColor"
															fill="none"
															stroke-width="2"
														/>
														<path
															d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"
															stroke="currentColor"
															fill="none"
															stroke-width="2"
														/>
													</svg>
												</button>
												<button
													class="btn-icon delete-btn"
													on:click={() => deleteMakeup(s)}
													title="삭제"
												>
													<svg width="14" height="14" viewBox="0 0 24 24" aria-hidden="true">
														<path
															d="M3 6h18M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0v14M10 11v6M14 11v6"
															stroke="currentColor"
															fill="none"
															stroke-width="2"
														/>
													</svg>
												</button>
											</div>
										{/if}
									</div>
								</div>
								{#if s.reason}
									<div class="modal-reason">{s.reason}</div>
								{/if}
							</li>
						{/each}
					</ul>
				{/if}
			{/if}
		</div>
	{/if}

	<!-- 메모장 패널 -->
	{#if showMemoPanel}
		<div class="memo-panel">
			<div class="memo-header">
				<h3>📝 관리자 메모장</h3>
				<button class="memo-close-btn" on:click={closeMemoPanel} aria-label="메모장 닫기">
					<svg width="20" height="20" viewBox="0 0 24 24">
						<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" fill="none" stroke-width="2"/>
					</svg>
				</button>
			</div>
			
			<div class="memo-content">
				<!-- 메모 편집 영역 -->
				<div class="memo-editor">
					<textarea 
						bind:value={memoContent}
						placeholder="여기에 메모를 자유롭게 작성하세요..."
						class="memo-textarea"
						on:input={handleMemoChange}
					></textarea>
				</div>
				
				<!-- 메모 액션 버튼들 -->
				<div class="memo-actions">
					<button 
						class="btn primary memo-save-btn" 
						on:click={() => saveMemo(false)}
						disabled={!memoContent.trim()}
					>
						<svg width="16" height="16" viewBox="0 0 24 24">
							<path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" stroke="currentColor" fill="none" stroke-width="2"/>
							<polyline points="17,21 17,13 7,13 7,21" stroke="currentColor" fill="none" stroke-width="2"/>
							<polyline points="7,3 7,8 15,8" stroke="currentColor" fill="none" stroke-width="2"/>
						</svg>
						저장
					</button>
					
					<button 
						class="btn ghost memo-clear-btn" 
						on:click={clearMemo}
						disabled={!memoContent.trim()}
					>
						<svg width="16" height="16" viewBox="0 0 24 24">
							<path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0v14M10 11v6M14 11v6" stroke="currentColor" fill="none" stroke-width="2"/>
						</svg>
						전체 삭제
					</button>
				</div>
				
				<!-- 메모 정보 -->
				<div class="memo-info">
					<div class="memo-stats">
						<span class="memo-char-count">{memoContent.length}자</span>
						{#if memoLastSaved}
							<span class="memo-last-saved">마지막 저장: {formatLastSaved(memoLastSaved)}</span>
						{/if}
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	:root {
		--bg: #f6f7fb;
		--card: #ffffff;
		--text: #1f2937;
		--muted: #6b7280;
		--line: #e5e7eb;
		--brand: #3b82f6;
		--brand-600: #2563eb;
		--sun: #ef4444;
		--sat: #2563eb;
		--today-bg: #fff4d6;
		--today-br: #f59e0b;
		--pill-bg: #eef6ff;
		--pill-br: #d6e7ff;
		--badge: #111827;
		/* Safe area + dynamic viewport height */
		--safe-top: env(safe-area-inset-top, 0px);
		--vh: 100vh; /* onMount에서 실제 값으로 덮어씀 */
	}



	* {
		box-sizing: border-box;
	}
	.wrap {
		max-width: 1100px;
		margin: 20px auto;
		/* 상단이 잘리는 문제 방지: 기본 패딩 vs safe-area vs 비례 여백 중 가장 큰 값 */
		padding-top: max(clamp(12px, 2vw, 24px), var(--safe-top), 1.5vh);
		padding-left: clamp(12px, 2vw, 24px);
		padding-right: clamp(12px, 2vw, 24px);
		padding-bottom: clamp(12px, 2vw, 24px);
		min-height: calc(var(--vh) - 40px); /* 최소 높이 유지 */
		height: auto; /* 내용에 맞게 자동 조정 */
		background: var(--bg);
		color: var(--text);
		border-radius: 20px;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
	}

	/* 모바일에서 좌우 여백 증가 */
	@media (max-width: 768px) {
		.wrap {
			margin: 6px 6px;
			border-radius: 16px;
		}
	}

	/* 상단바: 모바일에서 sticky 제거(가림 방지), 태블릿 이상에서 sticky 적용 */
	.topbar {
		display: flex;
		flex-wrap: wrap;
		gap: 1px 12px;
		justify-content: space-between;
		align-items: center;
		background: var(--card);
		border: 1px solid var(--line);
		border-radius: 12px;
		/* 상단 잘림 방지: 안전영역/줌에 비례한 패딩 */
		padding-top: max(14px, var(--safe-top), 1vh);
		padding-bottom: 14px;
		padding-left: clamp(12px, 2vw, 18px);
		padding-right: clamp(12px, 2vw, 18px);
		position: static; /* 모바일 기본: sticky 해제 */
		z-index: auto;
	}
	@media (min-width: 769px) {
		.topbar {
			position: sticky;
			top: 0;
			z-index: 10;
			border-radius: 16px;
		}
	}

	.brand h1 {
		display: flex;
		align-items: center;
		gap: 10px;
		font-size: clamp(18px, 2.2vw, 22px);
		margin: 0 0 6px 0;
		font-weight: 700;
	}
	.logo-dot {
		width: 10px;
		height: 10px;
		border-radius: 999px;
		background: var(--brand);
		box-shadow: 0 0 0 4px color-mix(in srgb, var(--brand) 20%, transparent);
	}
	.month-line {
		display: flex;
		align-items: center;
		gap: 10px;
		flex-wrap: wrap;
	}
	.month-text {
		font-size: clamp(14px, 1.6vw, 16px);
	}

	.user-info {
		display: flex;
		align-items: center;
		gap: 8px;
		margin-top: 4px;
		font-size: 12px;
	}

	.user-role {
		padding: 2px 6px;
		border-radius: 4px;
		font-weight: 600;
		font-size: 11px;
	}

	.user-role.admin {
		background: var(--brand);
		color: white;
	}

	.user-role.student {
		background: var(--muted);
		color: white;
	}

	.user-name {
		color: var(--muted);
		font-weight: 500;
	}
	.month-input-wrapper {
		position: relative;
		display: inline-block;
	}
	.month-input {
		appearance: none;
		border: 1px solid var(--line);
		background: var(--card);
		color: var(--text);
		padding: 6px 10px;
		border-radius: 10px;
		font-size: 14px;
	}
	.calendar-icon {
		position: absolute;
		right: 8px;
		top: 50%;
		transform: translateY(-50%);
		pointer-events: none;
		opacity: 0.6;
		display: none; /* 데스크탑에서는 숨김 */
	}

	.actions {
		display: flex;
		flex-direction: column;
		gap: 8px;
		align-items: flex-end;
	}
	.actions-top {
		display: flex;
		justify-content: flex-end;
	}
	.actions-bottom {
		display: flex;
		align-items: center;
		gap: 8px;
		flex-wrap: wrap;
	}
	.seg {
		display: flex;
		gap: 6px;
		background: color-mix(in srgb, var(--brand) 8%, transparent);
		padding: 4px;
		border-radius: 10px;
	}
	.seg.view {
		margin-left: 4px;
	}

	.btn {
		border: 1px solid var(--line);
		background: var(--card);
		color: var(--text);
		border-radius: 10px;
		padding: 8px 12px;
		font-weight: 600;
		cursor: pointer;
		transition: 0.15s ease;
	}
	.btn:hover {
		transform: translateY(-1px);
	}
	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
		transform: none;
	}

	.btn.ghost {
		background: transparent;
	}
	.btn.primary {
		background: var(--brand);
		border-color: var(--brand);
		color: white;
	}
	.btn.icon {
		width: 38px;
		height: 38px;
		display: grid;
		place-items: center;
	}
	.btn.icon svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}
	.btn.dashboard {
		background: transparent;
		color: var(--muted);
		border-color: var(--line);
		align-items: center;
		gap: 6px;
		font-size: 14px;
		margin-right: 12px;
	}
	.btn.dashboard:hover {
		background: var(--bg);
		color: var(--text);
		border-color: var(--brand);
	}
	.btn.dashboard svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	.btn.logout {
		background: transparent;
		color: var(--muted);
		border-color: var(--line);
		align-items: center;
		gap: 6px;
		font-size: 14px;
		margin-left: auto;
	}
	.btn.logout:hover {
		background: var(--bg);
		color: var(--text);
		border-color: var(--brand);
	}
	.btn.logout svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	.card {
		background: var(--card);
		border: 1px solid var(--line);
		border-radius: 16px;
		padding: 12px;
		margin-top: 14px;
		min-height: calc(var(--vh) - 200px); /* 최소 높이 설정 */
		height: auto; /* 내용에 맞게 자동 조정 */
	}

	/* 스와이프 시 시각적 피드백 */
	.calendar {
		user-select: none;
		touch-action: pan-y;
	}

	/* Skeleton */
	.skeleton {
		border-radius: 10px;
		background: linear-gradient(
			90deg,
			rgba(0, 0, 0, 0.06),
			rgba(0, 0, 0, 0.12),
			rgba(0, 0, 0, 0.06)
		);
		background-size: 200% 100%;
		animation: shine 1.2s linear infinite;
	}
	@keyframes shine {
		to {
			background-position: -200% 0;
		}
	}
	.skeleton.head {
		height: 24px;
		margin-bottom: 12px;
	}
	.skeleton.row {
		height: 80px;
		margin: 8px 0;
	}

	/* 캘린더 헤더(요일) & 그리드 */
	/* 모바일(기본): 수요일, 토요일 확장 */
	.weekdays {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr 1.5fr 1fr 1fr 1.5fr; /* 수요일, 토요일 확장 */
		position: static; /* 모바일 sticky 해제 */
		background: var(--card);
		z-index: 3;
		border-bottom: 1px solid var(--line);
	}
	.wday {
		text-align: center;
		padding: 8px 4px;
		font-weight: 700;
		font-size: clamp(11px, 1.4vw, 12px);
		color: var(--muted);
	}
	.wday.sun {
		color: var(--sun);
	}
	.wday.sat {
		color: var(--sat);
	}
	.wday.wed {
		color: var(--brand-600);
	}

	.grid {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr 1.5fr 1fr 1fr 1.5fr; /* 수요일, 토요일 확장 */
		grid-auto-rows: minmax(80px, 1fr); /* 최소 높이 증가 */
		gap: 3px;
		padding-top: 8px;
		min-height: calc(var(--vh) - 300px); /* 최소 높이 설정 */
		height: auto; /* 내용에 맞게 자동 조정 */
	}

	/* 태블릿 이상: 수요일(4번째), 토요일(7번째) 열 넓힘 */
	@media (min-width: 769px) {
		.weekdays {
			position: sticky;
			top: calc(64px + var(--safe-top)); /* 상단바 높이 + 안전영역 */
			grid-template-columns: 1fr 1fr 1fr 1.35fr 1fr 1fr 1.35fr;
		}
		.grid {
			grid-template-columns: 1fr 1fr 1fr 1.35fr 1fr 1fr 1.35fr;
			grid-auto-rows: minmax(85px, 1fr); /* 최소 높이 증가 */
			gap: 4px;
		}
	}

	/* 데스크톱: 조금 더 넓게 */
	@media (min-width: 980px) {
		.weekdays {
			grid-template-columns: 1fr 1fr 1fr 1.5fr 1fr 1fr 1.5fr;
		}
		.grid {
			grid-template-columns: 1fr 1fr 1fr 1.5fr 1fr 1fr 1.5fr;
			grid-auto-rows: minmax(95px, 1fr);
		}
	}

	/* 개별 Day 셀 */
	.cal-day {
		border: 1px solid var(--line);
		border-radius: 12px;
		padding: 6px;
		display: flex;
		flex-direction: column;
		min-height: 96px;
		height: auto;
		max-height: 120px; /* 최대 높이 제한 */
		outline: none;
		background: color-mix(in srgb, var(--card) 92%, var(--bg) 8%);
		cursor: pointer;
		transition: all 0.2s ease;
		overflow: hidden; /* 기본적으로는 숨김 */
		position: relative;
	}

	/* 스크롤 가능한 셀에 마우스 오버 시 */
	.cal-day.scrollable {
		overflow-y: auto; /* 세로 스크롤 활성화 */
		overflow-x: hidden; /* 가로 스크롤은 숨김 */
		/* 스크롤바 숨기기 */
		scrollbar-width: none; /* Firefox */
		-ms-overflow-style: none; /* IE and Edge */
	}

	/* 스크롤 가능한 셀 내부의 리스트와 시험 기간 컨테이너 */
	.cal-day.scrollable .list,
	.cal-day.scrollable .exam-periods {
		overflow: visible; /* 스크롤 가능한 상태에서는 내부 오버플로우 허용 */
		max-height: none; /* 최대 높이 제한 해제 */
	}

	/* 웹킷 기반 브라우저에서 스크롤바 숨기기 */
	.cal-day.scrollable::-webkit-scrollbar {
		display: none;
	}

	/* 스크롤 가능한 셀에 시각적 힌트 */
	.cal-day.scrollable::after {
		content: '';
		position: absolute;
		top: 2px;
		right: 2px;
		width: 4px;
		height: 4px;
		background: var(--brand);
		border-radius: 50%;
		opacity: 0.6;
		pointer-events: none;
	}
	.cal-day:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
		border-color: var(--brand);
	}
	.cal-day:focus {
		box-shadow: 0 0 0 2px var(--brand-600) inset;
	}
	.cal-day.is-out {
		opacity: 0.6;
	}
	.cal-day.is-today {
		background: var(--today-bg);
		border-color: var(--today-br);
	}
	.cal-day.sun .num {
		color: var(--sun);
	}
	.cal-day.sat .num {
		color: var(--sat);
	}
	.cal-day.wed .num {
		color: var(--brand-600);
	}
	/* 월화목금 글자색을 회색으로 */
	.cal-day:not(.sun):not(.sat):not(.wed) .num {
		color: var(--muted);
	}
	/* 나머지 요일(월,화,목,금,일) 셀 배경만 회색 처리 - 이번 달 셀만 */
	.cal-day:not(.sat):not(.wed):not(.is-out) {
		background: color-mix(in srgb, var(--card) 95%, var(--muted) 5%);
		opacity: 0.8;
	}

	.num-wrap {
		display: flex;
		align-items: center;
		justify-content: flex-start;
		flex-shrink: 0; /* 날짜 번호 영역은 고정 */
		margin-bottom: 4px;
		min-height: 20px; /* 최소 높이 보장 */
	}
	.num {
		font-weight: 800;
		font-size: clamp(13px, 1.6vw, 15px);
		line-height: 1;
		display: inline-block;
		padding: 4px 6px;
		border-radius: 6px;
		transition: all 0.2s ease;
	}
	.num.today {
		background: var(--today-br);
		color: white;
		box-shadow: 0 2px 4px rgba(245, 158, 11, 0.3);
	}

	.list {
		margin-top: 8px;
		display: flex;
		flex-direction: column;
		gap: 4px;
		flex: 1; /* 남은 공간을 모두 차지 */
		min-height: 0; /* flex 아이템이 축소될 수 있도록 */
		position: relative;
		max-height: calc(100% - 40px); /* 날짜 번호와 여백을 제외한 최대 높이 */
	}
	.pill {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 6px;
		padding: 4px 6px;
		border-radius: 8px;
		border: 1px solid var(--pill-br);
		background: var(--pill-bg);
		font-size: clamp(10px, 1.4vw, 11px);
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
		min-height: 20px;
		max-height: 24px;
		flex-shrink: 0;
	}
	.pill-name {
		font-weight: 700;
		overflow: hidden;
		text-overflow: ellipsis;
		flex-shrink: 0; /* 이름은 항상 표시 */
		min-width: 0; /* flex item이 축소될 수 있도록 */
	}
	.pill-time {
		color: var(--muted);
		font-variant-numeric: tabular-nums;
		flex: 1; /* 시간은 남은 공간 사용 */
		font-size: clamp(9px, 1.2vw, 10px);
		text-align: right;
	}
	.pill.more {
		background: var(--brand);
		border-color: var(--brand-600);
		color: white;
	}
	.pill.more .pill-name {
		color: white;
	}
	.pill.more .pill-time {
		color: rgba(255, 255, 255, 0.8);
	}

	/* 시험 기간 스타일 */
	.exam-periods {
		margin-top: 1.5px;
		display: flex;
		flex-direction: column;
		gap: 1.5px;
		flex: 1; /* 남은 공간을 모두 차지 */
		min-height: 0; /* flex 아이템이 축소될 수 있도록 */
		position: relative;
		max-height: calc(100% - 40px);
	}

	.exam-pill {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 3px;
		padding: 1.5px 3.5px;
		border-radius: 3.5px;
		border: 1px solid;
		font-size: clamp(6.5px, 0.9vw, 7.5px);
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
		min-height: 13px;
		max-height: 15px;
		flex-shrink: 0;
	}

	.exam-pill-name {
		font-weight: 550;
		overflow: hidden;
		text-overflow: ellipsis;
		flex: 1;
		text-align: center;
		font-size: clamp(5.5px, 0.8vw, 6.5px);
	}

	.exam-pill.more {
		background: #6b7280;
		border-color: #6b7280;
		color: white;
	}

	.exam-pill.more .exam-pill-name {
		color: white;
	}

	/* Agenda List */
	.agenda {
		padding: 10px;
	}
	.empty {
		color: var(--muted);
		padding: 12px 6px;
	}
	.ag-day {
		border-top: 1px dashed var(--line);
		padding: 10px 0;
	}
	.ag-day:first-child {
		border-top: none;
	}
	.ag-head {
		font-weight: 800;
		margin-bottom: 8px;
	}
	.ag-list {
		list-style: none;
		padding: 0;
		margin: 0;
		display: grid;
		gap: 8px;
	}
	.ag-item {
		border: 1px solid var(--line);
		border-radius: 12px;
		padding: 10px;
		background: color-mix(in srgb, var(--card) 92%, var(--bg) 8%);
	}
	.ag-main {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 10px;
	}
	.ag-name {
		font-weight: 800;
	}
	.ag-time {
		color: var(--muted);
		font-variant-numeric: tabular-nums;
	}
	.ag-reason {
		margin-top: 6px;
		color: var(--muted);
		font-size: 13px;
	}

	.ag-section {
		margin-bottom: 16px;
	}

	.ag-section:last-child {
		margin-bottom: 0;
	}

	.ag-section-title {
		font-size: 12px;
		font-weight: 700;
		color: var(--muted);
		margin-bottom: 8px;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.exam-period-item {
		border-left: 4px solid;
	}

	/* 영어 시험 날짜 강조 스타일 - 원래 색상 유지 */
	.exam-pill.english-exam {
		position: relative;
		border-width: 2px !important;
		box-shadow: 0 0 8px rgba(0, 0, 0, 0.2), 0 0 16px rgba(255, 255, 255, 0.3);
		animation: englishGlow 2s ease-in-out infinite;
		transform: scale(1.02);
	}

	.exam-pill.english-exam .exam-pill-name {
		font-weight: 700;
		text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
	}

	.english-indicator {
		position: absolute;
		top: -2px;
		right: -2px;
		background: #ff6b6b;
		color: white;
		font-size: 6px;
		font-weight: 900;
		padding: 1px 3px;
		border-radius: 3px;
		line-height: 1;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
		z-index: 1;
	}

	.ag-item.english-exam {
		position: relative;
		border-left-width: 6px !important;
		box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1), 0 0 20px rgba(255, 255, 255, 0.2);
		transform: translateX(2px);
	}

	.english-badge {
		background: #ff6b6b;
		color: white;
		font-size: 8px;
		font-weight: 900;
		padding: 1px 4px;
		border-radius: 3px;
		margin-left: 6px;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
	}

	.english-label {
		background: rgba(255, 107, 107, 0.1);
		color: #ff6b6b;
		font-size: 10px;
		font-weight: 700;
		padding: 2px 6px;
		border-radius: 4px;
		margin-left: 8px;
		border: 1px solid rgba(255, 107, 107, 0.3);
	}

	@keyframes englishGlow {
		0%, 100% {
			box-shadow: 0 0 8px rgba(0, 0, 0, 0.2), 0 0 16px rgba(255, 255, 255, 0.3);
		}
		50% {
			box-shadow: 0 0 12px rgba(0, 0, 0, 0.3), 0 0 24px rgba(255, 255, 255, 0.5);
		}
	}

	/* Modal */
	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.5);
		z-index: 9999;
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 20px;
		backdrop-filter: blur(4px);
		width: 100vw;
		height: 100vh;
		min-height: 100vh;
		/* 전역 main 스타일 무시 */
		align-items: center !important;
		justify-content: center !important;
	}
	.modal-content {
		background: var(--card);
		border: 1px solid var(--line);
		border-radius: 16px;
		padding: 24px;
		width: 100%;
		max-width: 500px;
		max-height: 80vh;
		overflow-y: auto;
		z-index: 10000;
		position: relative;
		display: flex;
		flex-direction: column;
		gap: 16px;
		box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
		animation: modalSlideIn 0.3s ease-out;
		margin: auto;
		align-self: center;
		/* 강제 중앙 정렬 */
		transform: translate(-50%, -50%);
		position: absolute;
		top: 50%;
		left: 50%;
	}
	@keyframes modalSlideIn {
		from {
			opacity: 0;
			transform: translate(-50%, -50%) scale(0.95);
		}
		to {
			opacity: 1;
			transform: translate(-50%, -50%) scale(1);
		}
	}
	.modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 20px;
		padding-bottom: 16px;
		border-bottom: 2px solid var(--line);
	}
	.modal-content h2 {
		font-size: clamp(18px, 2.2vw, 22px);
		font-weight: 700;
		margin: 0;
		color: var(--text);
	}
	.btn-close {
		background: none;
		border: none;
		color: var(--muted);
		cursor: pointer;
		padding: 8px;
		border-radius: 8px;
		transition: all 0.2s ease;
	}
	.btn-close:hover {
		background: var(--bg);
		color: var(--text);
	}
	.modal-list {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 12px;
		flex: 1;
	}
	.modal-item {
		border: 1px solid var(--line);
		border-radius: 12px;
		padding: 16px;
		background: color-mix(in srgb, var(--card) 92%, var(--bg) 8%);
		transition: all 0.2s ease;
	}
	.modal-item:hover {
		border-color: var(--brand);
		transform: translateX(4px);
	}
	.modal-loading {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}
	.modal-empty {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
		padding: 40px 20px;
		color: var(--muted);
		text-align: center;
	}
	.modal-empty svg {
		opacity: 0.5;
	}
	.modal-empty p {
		font-size: 16px;
		margin: 0;
	}

	.modal-main {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 10px;
		margin-bottom: 8px;
	}

	.modal-actions-item {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.item-actions {
		display: flex;
		gap: 4px;
		opacity: 0;
		transition: opacity 0.2s ease;
	}

	.modal-item:hover .item-actions {
		opacity: 1;
	}

	.btn-icon {
		background: none;
		border: none;
		padding: 6px;
		border-radius: 6px;
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.btn-icon:hover {
		transform: scale(1.1);
	}

	.edit-btn {
		color: var(--brand);
	}

	.edit-btn:hover {
		background: color-mix(in srgb, var(--brand) 10%, transparent);
	}

	.delete-btn {
		color: var(--sun);
	}

	.delete-btn:hover {
		background: color-mix(in srgb, var(--sun) 10%, transparent);
	}

	.btn-icon svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}
	.modal-info {
		display: flex;
		align-items: center;
		gap: 8px;
	}
	.modal-name {
		font-weight: 800;
		font-size: 16px;
		color: var(--text);
	}
	.modal-index {
		font-size: 12px;
		color: var(--muted);
		font-weight: 600;
		background: var(--bg);
		padding: 2px 6px;
		border-radius: 4px;
	}
	.modal-time {
		color: var(--brand);
		font-variant-numeric: tabular-nums;
		font-weight: 600;
		font-size: 14px;
		background: var(--pill-bg);
		padding: 4px 8px;
		border-radius: 6px;
	}
	.modal-reason {
		margin-top: 8px;
		color: var(--muted);
		font-size: 14px;
		line-height: 1.4;
		padding: 8px 12px;
		background: color-mix(in srgb, var(--bg) 30%, transparent);
		border-radius: 8px;
	}

	/* 보강 생성 폼 스타일 */
	.create-form {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.create-form h3 {
		font-size: 18px;
		font-weight: 700;
		margin: 0;
		color: var(--text);
		text-align: center;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.form-group label {
		font-weight: 600;
		font-size: 14px;
		color: var(--text);
	}

	.form-input,
	.form-textarea {
		border: 1px solid var(--line);
		background: var(--card);
		color: var(--text);
		padding: 12px;
		border-radius: 10px;
		font-size: 14px;
		transition: all 0.2s ease;
	}

	.form-input:focus,
	.form-textarea:focus {
		outline: none;
		border-color: var(--brand);
		box-shadow: 0 0 0 2px color-mix(in srgb, var(--brand) 20%, transparent);
	}

	.form-textarea {
		resize: vertical;
		min-height: 80px;
		font-family: inherit;
	}

	.form-actions {
		display: flex;
		gap: 12px;
		justify-content: flex-end;
		margin-top: 8px;
	}

	.create-btn {
		display: flex;
		align-items: center;
		gap: 8px;
		margin-top: 16px;
		padding: 12px 20px;
		font-size: 14px;
		font-weight: 600;
	}

	.create-btn svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	.modal-header-actions {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 16px;
		padding: 12px 16px;
		background: color-mix(in srgb, var(--bg) 30%, transparent);
		border-radius: 12px;
		border: 1px solid var(--line);
	}

	.schedule-summary {
		display: flex;
		align-items: center;
	}

	.schedule-count {
		font-size: 14px;
		color: var(--text);
		font-weight: 600;
	}

	.create-btn.compact {
		padding: 8px 16px;
		font-size: 13px;
		font-weight: 600;
		border-radius: 8px;
		display: flex;
		align-items: center;
		gap: 6px;
		transition: all 0.2s ease;
	}

	.create-btn.compact:hover {
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
	}

	.create-btn.compact svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	.modal-empty {
		transition: all 0.2s ease;
	}

	.modal-empty[role='button'] {
		cursor: pointer;
	}

	.modal-empty[role='button']:hover {
		transform: scale(1.02);
		background: color-mix(in srgb, var(--bg) 20%, transparent);
		border-radius: 12px;
	}

	.click-hint {
		font-size: 14px;
		color: var(--brand);
		font-weight: 600;
		margin-top: 8px;
		opacity: 0.8;
	}

	/* 구형 iOS 대체 (max() 미지원) */
	@supports not (padding: max(0px)) {
		.wrap {
			padding-top: calc(clamp(12px, 2vw, 24px) + constant(safe-area-inset-top));
		}
		.topbar {
			padding-top: calc(14px + constant(safe-area-inset-top));
		}
	}
	@media (max-width: 768px) {
		.month-line {
			flex-wrap: nowrap; /* 월 텍스트/인풋도 한 줄 */
			gap: 6px;
		}
		.month-text {
			font-size: 15px;
			margin-left: 16px;
		}
		.month-input {
			margin-left: 18px;
			margin-bottom: 4px;
			width: 33px; /* 모바일에서 달력 아이콘만 표시 */
			height: 26px;
			padding: 0;
			border-radius: 8px;
			font-size: 0; /* 텍스트 숨김 */
			background: transparent;
			border: 1px solid var(--line);
		}
		.calendar-icon {
			display: block; /* 모바일에서만 표시 */
			opacity: 1;
			width: 18px;
			height: 18px;
		}

		.actions {
			margin-left: auto; /* brand와 최대한 간격 벌림 */
			display: flex;
			justify-content: flex-end;
			align-items: center;
			gap: 6px; /* 버튼 간격 줄이기 */
			flex: 0 0 auto; /* 크기 줄도록 */
		}

		/* 기본 버튼 크기 축소 */
		.btn {
			padding: 2px 3px;
			font-size: 10px;
			border-radius: 6px;
		}

		.btn.icon {
			width: 24px;
			height: 26px;
		}

		/* 세그먼트 그룹(월/목록, 화살표 등)도 더 작게 */
		.seg {
			gap: 3px;
			padding: 3px;
			border-radius: 6px;
		}
		.seg .btn {
			font-size: 11px;
			padding: 4px 5px;
		}

		/* 대시보드 버튼 */
		.btn.dashboard {
			font-size: 11px;
			padding: 3px 4px;
			border-radius: 5px;
			opacity: 0.9;
		}
		.btn.dashboard .dashboard-icon {
			width: 12px;
			height: 12px;
		}

		/* 로그아웃 버튼은 특별히 더 작게 */
		.btn.logout {
			margin-left: 8px;
			font-size: 11px;
			padding: 3px 4px;
			border-radius: 5px;
			opacity: 0.9;
		}
		.btn.logout .logout-icon {
			width: 12px;
			height: 12px;
		}

		/* 푸시 토글 버튼 스타일 */
		.push-toggle {
			display: flex;
			align-items: center;
			gap: 6px;
			font-size: 14px;
			margin-right: 12px;
		}

		.push-toggle svg {
			stroke: currentColor;
			fill: none;
			stroke-width: 2;
		}

		.push-toggle:hover {
			background: var(--bg);
			color: var(--text);
			border-color: var(--brand);
		}

		.push-toggle:disabled {
			opacity: 0.6;
			cursor: not-allowed;
		}

		/* 모바일에서 푸시 버튼 크기 조정 */
		@media (max-width: 768px) {
			.push-toggle {
				font-size: 11px;
				padding: 3px 4px;
				border-radius: 5px;
				margin-right: 8px;
			}

			.push-toggle svg {
				width: 12px;
				height: 12px;
			}
		}
	}

	/* 메모 아이콘 스타일 */
	.memo-toggle {
		display: flex;
		align-items: center;
		gap: 6px;
		margin-right: 8px;
	}

	.memo-toggle svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	.memo-toggle:hover {
		background: var(--bg);
		color: var(--text);
		border-color: var(--brand);
	}

	/* 메모 패널 스타일 */
	.memo-panel {
		position: fixed;
		top: 0;
		right: 0;
		width: 400px;
		height: 100vh;
		background: var(--card);
		border-left: 1px solid var(--line);
		z-index: 1000;
		display: flex;
		flex-direction: column;
		box-shadow: -4px 0 20px rgba(0, 0, 0, 0.1);
		transform: translateX(100%);
		animation: slideInRight 0.3s ease-out forwards;
	}

	@keyframes slideInRight {
		to {
			transform: translateX(0);
		}
	}

	.memo-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20px;
		border-bottom: 1px solid var(--line);
		background: var(--bg);
	}

	.memo-header h3 {
		margin: 0;
		font-size: 18px;
		font-weight: 700;
		color: var(--text);
	}

	.memo-close-btn {
		background: none;
		border: none;
		color: var(--muted);
		cursor: pointer;
		padding: 8px;
		border-radius: 8px;
		transition: all 0.2s ease;
	}

	.memo-close-btn:hover {
		background: var(--line);
		color: var(--text);
	}

	.memo-content {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 20px;
	}

	.memo-editor {
		flex: 1;
		display: flex;
		flex-direction: column;
		margin-bottom: 20px;
	}

	.memo-textarea {
		flex: 1;
		border: 1px solid var(--line);
		border-radius: 12px;
		padding: 16px;
		font-family: inherit;
		font-size: 14px;
		line-height: 1.6;
		resize: none;
		background: var(--card);
		color: var(--text);
		transition: all 0.2s ease;
	}

	.memo-textarea:focus {
		outline: none;
		border-color: var(--brand);
		box-shadow: 0 0 0 2px color-mix(in srgb, var(--brand) 20%, transparent);
	}

	.memo-textarea::placeholder {
		color: var(--muted);
		opacity: 0.7;
	}

	.memo-actions {
		display: flex;
		gap: 12px;
		margin-bottom: 16px;
	}

	.memo-save-btn,
	.memo-clear-btn {
		display: flex;
		align-items: center;
		gap: 8px;
		flex: 1;
		justify-content: center;
	}

	.memo-save-btn svg,
	.memo-clear-btn svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	.memo-clear-btn {
		color: var(--sun);
		border-color: var(--sun);
	}

	.memo-clear-btn:hover {
		background: var(--sun);
		color: white;
	}

	.memo-info {
		border-top: 1px solid var(--line);
		padding-top: 16px;
	}

	.memo-stats {
		display: flex;
		justify-content: space-between;
		align-items: center;
		font-size: 12px;
		color: var(--muted);
	}

	.memo-char-count {
		font-weight: 600;
	}

	.memo-last-saved {
		opacity: 0.8;
	}

	/* 토스트 애니메이션 */
	@keyframes slideInRight {
		from {
			transform: translateX(100%);
			opacity: 0;
		}
		to {
			transform: translateX(0);
			opacity: 1;
		}
	}

	@keyframes slideOutRight {
		from {
			transform: translateX(0);
			opacity: 1;
		}
		to {
			transform: translateX(100%);
			opacity: 0;
		}
	}

	/* 모바일 대응 */
	@media (max-width: 768px) {
		.memo-panel {
			width: 100%;
			left: 0;
			right: 0;
		}
		
		.memo-toggle {
			font-size: 11px;
			padding: 3px 4px;
		}
		
		.memo-toggle svg {
			width: 14px;
			height: 14px;
		}
		
		.memo-actions {
			flex-direction: column;
		}
		
		.memo-save-btn,
		.memo-clear-btn {
			flex: none;
		}
	}

	/* 상태 표시 스타일 */
	.ag-status,
	.modal-status {
		display: inline-block;
		padding: 2px 6px;
		border-radius: 4px;
		font-size: 11px;
		font-weight: 500;
		border: 1px solid;
		margin-left: 8px;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	/* 상태 변경 컨테이너 */
	.status-change-container {
		margin-right: 8px;
	}

	/* 상태 선택 드롭다운 */
	.status-select {
		padding: 4px 8px;
		border: 1px solid var(--border);
		border-radius: 6px;
		background: var(--bg);
		color: var(--text);
		font-size: 12px;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
		min-width: 70px;
	}

	.status-select:hover {
		border-color: var(--brand);
		background: var(--hover);
	}

	.status-select:focus {
		outline: none;
		border-color: var(--brand);
		box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.1);
	}
</style>
