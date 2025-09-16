<script>
	import { onMount } from 'svelte';
	import { isAuthenticated, fetchWithAuth } from '$lib/auth';
	import { goto } from '$app/navigation';
	import { API_ENDPOINTS } from '$lib/config';
	import { enablePush, disablePush, getPushSubscription, getUserSubscriptions } from '$lib/push';
	import Toast from '$lib/components/Toast.svelte';
	import TodosPanel from '$lib/components/TodosPanel.svelte';
	import AssignmentPanel from '$lib/components/AssignmentPanel.svelte';
	import MakeupModal from '$lib/components/MakeupModal.svelte';
	import AgendaView from '$lib/components/AgendaView.svelte';
	import CalendarHeader from '$lib/components/CalendarHeader.svelte';
	import CalendarGrid from '$lib/components/CalendarGrid.svelte';
	import { getMakeupStatusColor, getMakeupStatusName } from '$lib/utils/makeupStatus.js';
	import { 
		getExamPeriodColor, 
		getExamPeriodDisplayName, 
		loadExamPeriods, 
		filterExamPeriodsForDate,
		isEnglishExam 
	} from '$lib/utils/examPeriod.js';
	import { 
		loadClasses, 
		loadExamPeriodsData, 
		loadCalendarData, 
		loadSelectedDateSchedules, 
		generateCalendar,
		toLocalDateString 
	} from '$lib/utils/calendarData.js';
	import { TouchState, createTouchHandlers } from '$lib/utils/touchHandler.js';
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

	// 사용자 정보 관련 상태
	let currentUser = null;
	let userRole = 'student'; // 기본값

	// 시험 기간 관련 상태
	let examPeriods = [];
	let examPeriodsLoading = false;

	// 할 일 관련 상태
	let showTodosPanel = false;

	// 과제 관련 상태
	let showAssignmentPanel = false;

	// 토스트 관련 상태
	let toastComponent;




	// 반 정보를 가져오는 함수
	let classes = [];
	async function loadClassesData() {
		classes = await loadClasses(fetchWithAuth, API_ENDPOINTS);
	}


	// 푸시 구독 관련 상태
	let pushSubscription = null;
	let isPushEnabled = false;
	let isPushLoading = false;

	// 터치 상태 관리
	let touchState = new TouchState();

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

		// 로그인 상태 확인
		const authenticated = await isAuthenticated();
		if (!authenticated) {
			await goto('/');
			return;
		}
		await loadCurrentUser();
		await loadClassesData();
		await loadCalendarDataWrapper();
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
	async function loadExamPeriodsDataWrapper() {
		examPeriodsLoading = true;
		try {
			examPeriods = await loadExamPeriodsData(fetchWithAuth, API_ENDPOINTS);
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

	async function loadCalendarDataWrapper() {
		isLoading = true;
		try {
			const makeupData = await loadCalendarData(fetchWithAuth, API_ENDPOINTS, currentYear, currentMonth, userRole, currentUser);

			// 시험 기간 데이터도 함께 로드
			await loadExamPeriodsDataWrapper();
			calendarDays = generateCalendar(makeupData, examPeriods, currentYear, currentMonth);
		} catch (e) {
			console.error('캘린더 데이터 로드 오류:', e);
			calendarDays = generateCalendar([], examPeriods, currentYear, currentMonth);
		} finally {
			isLoading = false;
		}
	}


	// 날짜 클릭 핸들러
	async function handleDateClick(date) {
		selectedDate = date;
		showModal = true;
		await loadSelectedDateSchedulesWrapper(date);
	}

	// 선택된 날짜의 보강 일정 로드
	async function loadSelectedDateSchedulesWrapper(date) {
		selectedDateSchedules = await loadSelectedDateSchedules(fetchWithAuth, API_ENDPOINTS, date, userRole, currentUser);
	}


	// 모달 닫기
	function closeModal() {
		showModal = false;
		selectedDate = null;
		selectedDateSchedules = [];
	}

	// 키보드 이벤트 처리 (ESC로 모달 닫기)
	function handleKeydown(event) {
		if (event.key === 'Escape') {
			if (showModal) {
					closeModal();
			} else if (showTodosPanel) {
				showTodosPanel = false;
			} else if (showAssignmentPanel) {
				showAssignmentPanel = false;
			}
		}
	}

	// 터치 시작 이벤트
	// 터치 이벤트 핸들러 생성
	const touchHandlers = createTouchHandlers(touchState, {
		onNextMonth: nextMonth,
		onPreviousMonth: previousMonth
	});

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

	// 할 일 관련 함수들
	function toggleTodosPanel() {
		showTodosPanel = !showTodosPanel;
	}

	// 토스트 알림 함수
	function showToast(message, type = 'info') {
		if (toastComponent) {
			toastComponent.showToast(message, type);
		}
	}

	// 과제 관련 함수들
	function toggleAssignmentPanel() {
		showAssignmentPanel = !showAssignmentPanel;
	}

	function previousMonth() {
		currentMonth--;
		if (currentMonth < 0) {
			currentMonth = 11;
			currentYear--;
		}
		loadCalendarDataWrapper();
	}

	function nextMonth() {
		currentMonth++;
		if (currentMonth > 11) {
			currentMonth = 0;
			currentYear++;
		}
		loadCalendarDataWrapper();
	}

	function goToday() {
		const now = new Date();
		currentYear = now.getFullYear();
		currentMonth = now.getMonth();
		loadCalendarDataWrapper();
	}



	function fmtTime(t) {
		if (!t) return '';
		return t.slice(0, 5); // HH:MM
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
	<CalendarHeader 
		{currentYear}
		{currentMonth}
		{currentUser}
		{userRole}
		{viewMode}
		{isLoading}
		{isPushEnabled}
		{isPushLoading}
		on:monthChange={(e) => {
			currentYear = e.detail.year;
			currentMonth = e.detail.month;
			loadCalendarDataWrapper();
		}}
		on:logout={handleLogout}
		on:goToday={goToday}
		on:previousMonth={previousMonth}
		on:nextMonth={nextMonth}
		on:viewModeChange={(e) => viewMode = e.detail}
		on:enablePush={handleEnablePush}
		on:disablePush={handleDisablePush}
		on:toggleAssignmentPanel={toggleAssignmentPanel}
		on:toggleTodosPanel={toggleTodosPanel}
	/>

	{#if isLoading}
		<div class="card">
			<div class="skeleton head" />
			<div class="skeleton row" />
			<div class="skeleton row" />
			<div class="skeleton row" />
		</div>
	{:else if viewMode === 'month'}
		<CalendarGrid 
			{calendarDays}
			{classes}
			{dayNames}
			on:dateClick={(e) => handleDateClick(e.detail)}
			on:touchStart={touchHandlers.onTouchStart}
			on:touchMove={touchHandlers.onTouchMove}
			on:touchEnd={touchHandlers.onTouchEnd}
		/>
	{:else}
		<AgendaView {agendaItems} {classes} />
						{/if}
						
	<!-- 보강 일정 모달 -->
	<MakeupModal 
		bind:isVisible={showModal}
		{selectedDate}
		{selectedDateSchedules}
		{userRole}
		{currentUser}
		on:close={closeModal}
		on:refresh={async () => {
			await loadSelectedDateSchedulesWrapper(selectedDate);
			await loadCalendarDataWrapper();
		}}
	/>



	<!-- 토스트 컴포넌트 -->
	<Toast bind:this={toastComponent} />

	<!-- 할 일 패널 컴포넌트 -->
	<TodosPanel 
		bind:isVisible={showTodosPanel}
		{userRole}
		{currentUser}
		on:toast={(e) => showToast(e.detail.message, e.detail.type)}
	/>

	<!-- 과제 패널 컴포넌트 -->
	<AssignmentPanel 
		bind:isVisible={showAssignmentPanel} 
		{userRole}
		{currentUser}
		on:toast={(e) => showToast(e.detail.message, e.detail.type)}
	/>
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




	/* 구형 iOS 대체 (max() 미지원) */
	@supports not (padding: max(0px)) {
		.wrap {
			padding-top: calc(clamp(12px, 2vw, 24px) + constant(safe-area-inset-top));
		}
	}

</style>
