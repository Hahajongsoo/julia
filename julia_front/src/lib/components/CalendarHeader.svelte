<script>
	import { createEventDispatcher } from 'svelte';
	import { goto } from '$app/navigation';

	const dispatch = createEventDispatcher();

	// Props
	export let currentYear;
	export let currentMonth;
	export let currentUser = null;
	export let userRole = 'student';
	export let viewMode = 'month';
	export let isLoading = false;
	export let isPushEnabled = false;
	export let isPushLoading = false;

	// 월 이름 배열
	const monthNames = [
		'1월', '2월', '3월', '4월', '5월', '6월',
		'7월', '8월', '9월', '10월', '11월', '12월',
	];

	// month input: value는 계산한 문자열을 단방향 주입하고, 변경은 on:change에서 처리
	function monthInputValue() {
		return `${currentYear}-${String(currentMonth + 1).padStart(2, '0')}`;
	}

	function onMonthInputChange(e) {
		const val = e.target.value; // yyyy-MM
		if (!val) return;
		const [y, m] = val.split('-').map(Number);
		dispatch('monthChange', { year: y, month: m - 1 });
	}

	// 이벤트 핸들러들
	function handleDashboardClick() {
		goto('/dashboard');
	}

	function handleLogout() {
		dispatch('logout');
	}

	function handleGoToday() {
		dispatch('goToday');
	}

	function handlePreviousMonth() {
		dispatch('previousMonth');
	}

	function handleNextMonth() {
		dispatch('nextMonth');
	}

	function handleViewModeChange(mode) {
		dispatch('viewModeChange', mode);
	}

	function handleEnablePush() {
		dispatch('enablePush');
	}

	function handleDisablePush() {
		dispatch('disablePush');
	}

	function handleToggleAssignmentPanel() {
		dispatch('toggleAssignmentPanel');
	}

	function handleToggleMemoPanel() {
		dispatch('toggleMemoPanel');
	}
</script>

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
					on:click={handleDashboardClick}
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
			<button class="btn ghost" on:click={handleGoToday} aria-label="오늘로 이동">오늘</button>
			
			<!-- 과제 아이콘 추가 -->
			<button 
				class="btn ghost assignment-toggle" 
				on:click={handleToggleAssignmentPanel}
				aria-label="과제 목록 열기"
				title="과제 목록"
			>
				<svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true">
					<path d="M9 12l2 2 4-4" stroke="currentColor" fill="none" stroke-width="2"/>
					<path d="M21 12c0 1.66-1.34 3-3 3H6c-1.66 0-3-1.34-3-3s1.34-3 3-3h12c1.66 0 3 1.34 3 3z" stroke="currentColor" fill="none" stroke-width="2"/>
				</svg>
				과제
			</button>
			
			<!-- 메모 아이콘 추가 -->
			{#if userRole === 'admin'}
				<button 
					class="btn ghost memo-toggle" 
					on:click={handleToggleMemoPanel}
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
					on:click={handlePreviousMonth}
					aria-label="이전 달"
					disabled={isLoading}
				>
					<svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true"
						><path d="M15 18l-6-6 6-6" stroke="currentColor" fill="none" stroke-width="2" /></svg
					>
				</button>
				<button class="btn icon" on:click={handleNextMonth} aria-label="다음 달" disabled={isLoading}>
					<svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true"
						><path d="M9 6l6 6-6 6" stroke="currentColor" fill="none" stroke-width="2" /></svg
					>
				</button>
			</div>
			<div class="seg view">
				<button
					class="btn {viewMode === 'month' ? 'primary' : 'ghost'}"
					on:click={() => handleViewModeChange('month')}
					aria-pressed={viewMode === 'month'}>월</button
				>
				<button
					class="btn {viewMode === 'agenda' ? 'primary' : 'ghost'}"
					on:click={() => handleViewModeChange('agenda')}
					aria-pressed={viewMode === 'agenda'}>목록</button
				>
			</div>
		</div>
	</div>
</header>

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

	/* 과제 아이콘 스타일 */
	.assignment-toggle {
		display: flex;
		align-items: center;
		gap: 6px;
		margin-right: 8px;
	}

	.assignment-toggle svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	.assignment-toggle:hover {
		background: var(--bg);
		color: var(--text);
		border-color: var(--brand);
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

	/* 모바일 대응 */
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

		/* 모바일에서 푸시 버튼 크기 조정 */
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

		.assignment-toggle {
			font-size: 11px;
			padding: 3px 4px;
		}
		
		.assignment-toggle svg {
			width: 14px;
			height: 14px;
		}
		
		.memo-toggle {
			font-size: 11px;
			padding: 3px 4px;
		}
		
		.memo-toggle svg {
			width: 14px;
			height: 14px;
		}
	}

	/* 구형 iOS 대체 (max() 미지원) */
	@supports not (padding: max(0px)) {
		.topbar {
			padding-top: calc(14px + constant(safe-area-inset-top));
		}
	}
</style>
