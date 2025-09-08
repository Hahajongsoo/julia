<script>
	import { createEventDispatcher } from 'svelte';
	import { API_ENDPOINTS } from '$lib/config';
	import { fetchWithAuth } from '$lib/auth';
	import { getMakeupStatusColor, getMakeupStatusName, makeupStatusOptions } from '$lib/utils/makeupStatus.js';

	const dispatch = createEventDispatcher();

	// Props
	export let isVisible = false;
	export let selectedDate = null;
	export let selectedDateSchedules = [];
	export let userRole = 'student';
	export let currentUser = null;

	// 모달 관련 상태
	let isModalLoading = false;
	let schedules = selectedDateSchedules;

	// selectedDateSchedules가 변경될 때마다 schedules 업데이트
	$: schedules = selectedDateSchedules;

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

	// 반 및 학생 관련 상태
	let classes = [];
	let students = [];
	let selectedClassId = '';

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

	// 과제 생성 관련 상태
	let showAssignmentForm = false;
	let assignmentFormData = {
		content: '',
		status: 'pending'
	};
	let selectedMakeupForAssignment = null;
	let isCreatingAssignment = false;


	// 날짜 포맷팅 함수
	function formatDate(date) {
		const year = date.getFullYear();
		const month = date.getMonth() + 1;
		const day = date.getDate();
		const dayNames = ['일', '월', '화', '수', '목', '금', '토'];
		const dayName = dayNames[date.getDay()];
		return `${year}년 ${month}월 ${day}일 (${dayName})`;
	}

	// 시간 포맷팅 함수
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

	// 모달 닫기
	function closeModal() {
		dispatch('close');
		showCreateForm = false;
		showEditForm = false;
		showAssignmentForm = false;
		selectedMakeup = null;
		selectedMakeupForAssignment = null;
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
		assignmentFormData = {
			content: '',
			status: 'pending'
		};
	}

	// 반 목록 가져오기
	async function loadClasses() {
		try {
			const res = await fetchWithAuth(API_ENDPOINTS.CLASSES);
			if (res.ok) {
				classes = await res.json();
			} else {
				console.error('반 목록 로드 실패');
			}
		} catch (e) {
			console.error('반 목록 로드 오류:', e);
		}
	}

	// 선택된 반의 학생 목록 가져오기
	async function loadStudentsByClass(classId) {
		if (!classId) {
			students = [];
			return;
		}
		
		try {
			const res = await fetchWithAuth(`${API_ENDPOINTS.CLASSES}/${classId}/users`);
			if (res.ok) {
				students = await res.json();
			} else {
				console.error('학생 목록 로드 실패');
				students = [];
			}
		} catch (e) {
			console.error('학생 목록 로드 오류:', e);
			students = [];
		}
	}

	// 반 선택 변경 시 학생 목록 업데이트
	$: if (selectedClassId) {
		loadStudentsByClass(selectedClassId);
		// 반이 변경되면 학생 선택 초기화
		createFormData.user_id = '';
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
		// 반과 학생 목록 초기화
		selectedClassId = '';
		students = [];
		// 반 목록 로드
		loadClasses();
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
				dispatch('refresh');
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
		// 반과 학생 목록 초기화
		selectedClassId = '';
		students = [];
		// 반 목록 로드
		loadClasses();
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

	// 과제 생성 폼 열기
	function openAssignmentForm(makeup) {
		selectedMakeupForAssignment = makeup;
		showAssignmentForm = true;
		assignmentFormData = {
			content: `${makeup.makeup_date} 보강 일정에 대한 과제입니다.\n`,
			status: 'pending'
		};
	}

	// 과제 생성 폼 닫기
	function closeAssignmentForm() {
		showAssignmentForm = false;
		selectedMakeupForAssignment = null;
		assignmentFormData = {
			content: '',
			status: 'pending'
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
				dispatch('refresh');
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
				dispatch('refresh');
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
				dispatch('refresh');
			} else {
				const errorData = await res.json();
				alert(`상태 변경 실패: ${errorData.message || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('상태 변경 오류:', e);
			alert('상태 변경 중 오류가 발생했습니다.');
		}
	}

	// 과제 생성 요청
	async function createAssignment() {
		if (!assignmentFormData.content) {
			alert('내용을 입력해주세요.');
			return;
		}

		isCreatingAssignment = true;
		try {
			const assignmentData = {
				user_id: selectedMakeupForAssignment.user_id,
				content: assignmentFormData.content,
				status: assignmentFormData.status
			};

			const res = await fetchWithAuth(API_ENDPOINTS.ASSIGNMENTS, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(assignmentData),
			});

			if (res.ok) {
				alert('과제가 성공적으로 생성되었습니다.');
				closeAssignmentForm();
				dispatch('refresh');
			} else {
				const errorData = await res.json();
				alert(`과제 생성 실패: ${errorData.message || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('과제 생성 오류:', e);
			alert('과제 생성 중 오류가 발생했습니다.');
		} finally {
			isCreatingAssignment = false;
		}
	}

	// 키보드 이벤트 처리 (ESC로 모달 닫기)
	function handleKeydown(event) {
		if (event.key === 'Escape') {
			if (showCreateForm) {
				closeCreateForm();
			} else if (showEditForm) {
				closeEditForm();
			} else if (showAssignmentForm) {
				closeAssignmentForm();
			} else {
				closeModal();
			}
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

{#if isVisible}
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
						<label for="create-class">반 선택</label>
						<select
							id="create-class"
							bind:value={selectedClassId}
							required
							class="form-input"
						>
							<option value="">반을 선택하세요</option>
							{#each classes as classItem}
								<option value={classItem.class_id}>{classItem.class_name}</option>
							{/each}
						</select>
					</div>

					<div class="form-group">
						<label for="create-user">학생 선택</label>
						<select
							id="create-user"
							bind:value={createFormData.user_id}
							required
							class="form-input"
							disabled={!selectedClassId}
						>
							<option value="">학생을 선택하세요</option>
							{#each students as student}
								<option value={student.id}>{student.id}</option>
							{/each}
						</select>
						{#if !selectedClassId}
							<div class="help-text">먼저 반을 선택해주세요</div>
						{:else if students.length === 0}
							<div class="help-text">선택된 반에 학생이 없습니다</div>
						{/if}
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
								{#each makeupStatusOptions as option}
									<option value={option.value}>{option.label}</option>
								{/each}
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
						<label for="edit-class">반 선택</label>
						<select
							id="edit-class"
							bind:value={selectedClassId}
							required
							class="form-input"
						>
							<option value="">반을 선택하세요</option>
							{#each classes as classItem}
								<option value={classItem.class_id}>{classItem.class_name}</option>
							{/each}
						</select>
					</div>

					<div class="form-group">
						<label for="edit-user">학생 선택</label>
						<select
							id="edit-user"
							bind:value={editFormData.user_id}
							required
							class="form-input"
							disabled={!selectedClassId}
						>
							<option value="">학생을 선택하세요</option>
							{#each students as student}
								<option value={student.id}>{student.id}</option>
							{/each}
						</select>
						{#if !selectedClassId}
							<div class="help-text">먼저 반을 선택해주세요</div>
						{:else if students.length === 0}
							<div class="help-text">선택된 반에 학생이 없습니다</div>
						{/if}
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
								{#each makeupStatusOptions as option}
									<option value={option.value}>{option.label}</option>
								{/each}
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
		{:else if showAssignmentForm}
			<!-- 과제 생성 폼 -->
			<div class="create-form">
				<h3>새 과제 생성</h3>
				<form on:submit|preventDefault={createAssignment}>
					<div class="form-group">
						<label for="assignment-content">과제 내용</label>
						<textarea
							id="assignment-content"
							bind:value={assignmentFormData.content}
							placeholder="과제 내용을 입력하세요"
							required
							class="form-textarea"
							rows="6"
						/>
					</div>

					<div class="form-group">
						<label for="assignment-status">상태</label>
						<select
							id="assignment-status"
							bind:value={assignmentFormData.status}
							required
							class="form-input"
						>
							<option value="pending">대기중</option>
							<option value="completed">완료</option>
						</select>
					</div>

					<div class="form-actions">
						<button
							type="button"
							class="btn ghost"
							on:click={closeAssignmentForm}
							disabled={isCreatingAssignment}
						>
							취소
						</button>
						<button type="submit" class="btn primary" disabled={isCreatingAssignment}>
							{isCreatingAssignment ? '생성 중...' : '과제 생성'}
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
			{:else if schedules.length === 0}
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
						<span class="schedule-count">총 {schedules.length}개의 보강 일정</span>
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
					{#each schedules as s, index}
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
												{#each makeupStatusOptions as option}
													<option value={option.value}>{option.label}</option>
												{/each}
											</select>
										</div>
										<button
											class="btn-icon assignment-btn"
											on:click={() => openAssignmentForm(s)}
											title="과제 생성"
										>
											<svg width="14" height="14" viewBox="0 0 24 24" aria-hidden="true">
												<path
													d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"
													stroke="currentColor"
													fill="none"
													stroke-width="2"
												/>
												<polyline
													points="14,2 14,8 20,8"
													stroke="currentColor"
													fill="none"
													stroke-width="2"
												/>
												<line
													x1="16"
													y1="13"
													x2="8"
													y2="13"
													stroke="currentColor"
													stroke-width="2"
												/>
												<line
													x1="16"
													y1="17"
													x2="8"
													y2="17"
													stroke="currentColor"
													stroke-width="2"
												/>
												<polyline
													points="10,9 9,9 8,9"
													stroke="currentColor"
													fill="none"
													stroke-width="2"
												/>
											</svg>
										</button>
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

	.assignment-btn {
		color: #10b981;
	}

	.assignment-btn:hover {
		background: color-mix(in srgb, #10b981 10%, transparent);
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

	/* 상태 표시 스타일 */
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
		border: 1px solid var(--line);
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
		background: var(--bg);
	}

	.status-select:focus {
		outline: none;
		border-color: var(--brand);
		box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.1);
	}

	/* 도움말 텍스트 */
	.help-text {
		font-size: 12px;
		color: var(--muted);
		margin-top: 4px;
	}
</style>
