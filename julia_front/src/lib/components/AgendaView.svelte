<script>
	import { getMakeupStatusColor, getMakeupStatusName } from '$lib/utils/makeupStatus.js';

	// Props
	export let agendaItems = [];
	export let classes = [];


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

	// 시험 기간에 반 이름을 추가하는 함수
	function getExamPeriodDisplayName(examPeriod) {
		const classItem = classes.find(c => c.class_id === examPeriod.class_id);
		const className = classItem ? classItem.class_name : '알 수 없는 반';
		return `${className} - ${examPeriod.name}`;
	}

	// 시간 포맷팅 함수
	function fmtTime(t) {
		if (!t) return '';
		return t.slice(0, 5); // HH:MM
	}
</script>

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
		margin-top: 14px;
		min-height: calc(var(--vh) - 200px); /* 최소 높이 설정 */
		height: auto; /* 내용에 맞게 자동 조정 */
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

	/* 상태 표시 스타일 */
	.ag-status {
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
</style>
