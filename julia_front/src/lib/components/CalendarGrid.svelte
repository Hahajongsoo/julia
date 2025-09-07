<script>
	import { createEventDispatcher } from 'svelte';
	import { getMakeupStatusColor, getMakeupStatusName } from '$lib/utils/makeupStatus.js';
	import { getExamPeriodColor, getExamPeriodDisplayName, isEnglishExam } from '$lib/utils/examPeriod.js';

	const dispatch = createEventDispatcher();

	// Props
	export let calendarDays = [];
	export let classes = [];
	export let dayNames = ['일', '월', '화', '수', '목', '금', '토'];

	// 스크롤 상태 관리
	let scrollState = new Map();

	// 이벤트 핸들러들
	function handleDateClick(date) {
		dispatch('dateClick', date);
	}

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

	function handleCellMouseEnter(event) {
		const cell = event.currentTarget;
		
		// 약간의 지연을 두고 스크롤 가능 여부 확인 (DOM 업데이트 대기)
		setTimeout(() => {
			if (!cell) return; // cell이 null인 경우 방어
			
			const { scrollHeight, clientHeight } = cell;
			
			// 스크롤 가능한 내용이 있으면 클래스 추가
			if (scrollHeight > clientHeight + 2) {
				cell.classList.add('scrollable');
				// 스크롤 상태 초기화
				scrollState.set(cell, { lastScrollTime: 0, consecutiveScrolls: 0 });
			}
		}, 10);
	}

	function handleCellMouseLeave(event) {
		const cell = event.currentTarget;
		if (!cell) return; // cell이 null인 경우 방어
		
		cell.classList.remove('scrollable');
		// 스크롤 상태 정리
		scrollState.delete(cell);
	}

	function handleTouchStart(event) {
		dispatch('touchStart', event);
	}

	function handleTouchMove(event) {
		dispatch('touchMove', event);
	}

	function handleTouchEnd(event) {
		dispatch('touchEnd', event);
	}

	// 유틸리티 함수들
	function dayClasses({ date, isCurrentMonth, isToday }) {
		const w = date.getDay(); // 0=일 ... 3=수, 6=토
		return [
			'cal-day',
			!isCurrentMonth && 'is-out',
			isToday && 'is-today',
			w === 0 && 'sun',
			w === 6 && 'sat',
			w === 3 && 'wed'
		].filter(Boolean).join(' ');
	}

	function fmtTime(timeStr) {
		if (!timeStr) return '';
		const [hours, minutes] = timeStr.split(':');
		return `${hours}:${minutes}`;
	}
</script>

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
							{@const displayName = getExamPeriodDisplayName(ep, classes)}
							{@const isEnglishDate = isEnglishExam(ep, day.date)}
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
	}

	.card {
		background: var(--card);
		border: 1px solid var(--line);
		border-radius: 16px;
		padding: 12px;
	}

	.calendar {
		display: grid;
		grid-template-rows: auto 1fr;
		gap: 8px;
		height: 100%;
		min-height: 400px;
	}

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

	/* 모바일 대응 */
	@media (max-width: 768px) {
		.weekdays {
			gap: 0;
			margin-bottom: 4px;
		}

		.wday {
			font-size: 12px;
			padding: 6px 2px;
		}

		.grid {
			gap: 0;
		}
	}

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

	.english-indicator {
		background: var(--brand);
		color: white;
		font-size: 8px;
		padding: 1px 3px;
		border-radius: 3px;
		margin-left: 4px;
		font-weight: 700;
	}

	.exam-pill.english-exam {
		border-style: dashed;
	}

	/* 모바일 대응 */
	@media (max-width: 768px) {
		.cal-day {
			padding: 4px;
			min-height: 60px;
		}

		.num {
			font-size: 12px;
		}

		.num.today {
			width: 16px;
			height: 16px;
			font-size: 10px;
		}

		.pill,
		.exam-pill {
			font-size: 10px;
			padding: 1px 4px;
		}

		.pill-time {
			font-size: 9px;
		}

		.english-indicator {
			font-size: 7px;
			padding: 1px 2px;
		}
	}
</style>
