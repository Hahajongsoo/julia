<script>
	import { createEventDispatcher } from 'svelte';
	import { fetchWithAuth } from '$lib/auth';
	import { API_ENDPOINTS } from '$lib/config';

	const dispatch = createEventDispatcher();

	export let isVisible = false;
	export let userRole = 'student';
	export let currentUser = null;

	let assignments = [];
	let assignmentsLoading = false;
	let showCreateAssignmentForm = false;
	let createAssignmentData = {
		user_id: '',
		content: '',
		status: 'pending'
	};
	let isCreatingAssignment = false;
	let showEditAssignmentForm = false;
	let editAssignmentData = {
		assignment_id: null,
		user_id: '',
		content: '',
		status: 'pending'
	};
	let isEditingAssignment = false;
	let selectedAssignment = null;

	// 반별 토글 상태 관리 (배열로 변경)
	let expandedClasses = [];
	
	// 반응형 선언으로 토글 상태 추적
	$: expandedClassesReactive = expandedClasses;

	// 반 토글 함수
	function toggleClass(classId) {
		const index = expandedClasses.indexOf(classId);
		if (index > -1) {
			// 이미 펼쳐져 있으면 접기
			expandedClasses = expandedClasses.filter(id => id !== classId);
		} else {
			// 접혀져 있으면 펼치기
			expandedClasses = [...expandedClasses, classId];
		}
	}

	// 반이 펼쳐져 있는지 확인 (반응형으로 변경)
	$: isClassExpanded = (classId) => {
		return expandedClasses.includes(classId);
	};

	// 과제 패널 열기/닫기
	function togglePanel() {
		isVisible = !isVisible;
		if (isVisible) {
			loadAssignments();
		}
		dispatch('toggle', { isVisible });
	}

	function closePanel() {
		isVisible = false;
		dispatch('close');
	}

	// 과제 로드 (반별 토글로 변경)
	async function loadAssignments() {
		assignmentsLoading = true;
		try {
			let response;
			if (userRole === 'admin') {
				// admin인 경우 /classes/assignments 사용
				response = await fetchWithAuth(API_ENDPOINTS.CLASS_ASSIGNMENTS);
			} else {
				// 학생인 경우 /assignments/user/:userID 사용
				const userId = currentUser?.id || '';
				response = await fetchWithAuth(API_ENDPOINTS.ASSIGNMENTS_BY_USER(userId));
			}

			if (response.ok) {
				const data = await response.json();
				if (userRole === 'admin') {
					// admin인 경우 ClassAssignments 구조 그대로 사용 (반별 토글)
					assignments = data.classes || [];
				} else {
					// 학생인 경우 단일 배열로 변환 (백엔드에서 이미 pending만 필터링됨)
					assignments = [{
						class_id: 0,
						class_name: '내 과제',
						students: [{
							user_id: currentUser?.id || '',
							assignments: data
						}]
					}];
				}
			} else {
				console.error('Failed to load assignments');
				assignments = [];
			}
		} catch (error) {
			console.error('Error loading assignments:', error);
			assignments = [];
		} finally {
			assignmentsLoading = false;
		}
	}

	// 과제 상태 변경 (체크박스 체크 시)
	async function updateAssignmentStatus(assignment, isCompleted) {
		try {
			const updateData = {
				assignment_id: assignment.assignment_id,
				user_id: assignment.user_id,
				content: assignment.content,
				status: isCompleted ? 'completed' : 'pending'
			};

			const response = await fetchWithAuth(API_ENDPOINTS.ASSIGNMENTS, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(updateData)
			});

			if (response.ok) {
				// 로컬 상태 업데이트 - 새로운 배열 생성으로 반응성 보장
				assignments = assignments.map(classBlock => ({
					...classBlock,
					students: classBlock.students.map(student => ({
						...student,
						assignments: student.assignments
							.map(a => a.assignment_id === assignment.assignment_id 
								? { ...a, status: isCompleted ? 'completed' : 'pending' }
								: a
							)
							.filter(a => !(a.assignment_id === assignment.assignment_id && isCompleted))
					}))
				}));
				
				
				dispatch('toast', { 
					message: isCompleted ? '과제가 완료되었습니다.' : '과제 상태가 변경되었습니다.', 
					type: 'success' 
				});
			} else {
				dispatch('toast', { message: '과제 상태 변경에 실패했습니다.', type: 'error' });
			}
		} catch (error) {
			console.error('Assignment status update error:', error);
			dispatch('toast', { message: '과제 상태 변경 중 오류가 발생했습니다.', type: 'error' });
		}
	}

	// 과제 생성 폼 열기
	function openCreateAssignmentForm() {
		showCreateAssignmentForm = true;
		createAssignmentData = {
			user_id: '',
			content: '',
			status: 'pending'
		};
	}

	// 과제 생성 폼 닫기
	function closeCreateAssignmentForm() {
		showCreateAssignmentForm = false;
		createAssignmentData = {
			user_id: '',
			content: '',
			status: 'pending'
		};
	}

	// 과제 생성 요청
	async function createAssignment() {
		if (!createAssignmentData.content || !createAssignmentData.user_id) {
			dispatch('toast', { message: '학생명과 과제 내용을 입력해주세요.', type: 'error' });
			return;
		}

		isCreatingAssignment = true;
		try {
			const createData = {
				user_id: createAssignmentData.user_id,
				content: createAssignmentData.content,
				status: createAssignmentData.status
			};

			const response = await fetchWithAuth(API_ENDPOINTS.ASSIGNMENTS, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(createData)
			});

			if (response.ok) {
				dispatch('toast', { message: '과제가 성공적으로 생성되었습니다.', type: 'success' });
				closeCreateAssignmentForm();
				// 과제 목록 다시 로드
				await loadAssignments();
			} else {
				const errorData = await response.json();
				dispatch('toast', { 
					message: `과제 생성 실패: ${errorData.error || '알 수 없는 오류가 발생했습니다.'}`, 
					type: 'error' 
				});
			}
		} catch (error) {
			console.error('Assignment creation error:', error);
			dispatch('toast', { message: '과제 생성 중 오류가 발생했습니다.', type: 'error' });
		} finally {
			isCreatingAssignment = false;
		}
	}

	// 과제 수정 폼 열기
	function openEditAssignmentForm(assignment) {
		selectedAssignment = assignment;
		showEditAssignmentForm = true;
		editAssignmentData = {
			assignment_id: assignment.assignment_id,
			user_id: assignment.user_id,
			content: assignment.content,
			status: assignment.status
		};
	}

	// 과제 수정 폼 닫기
	function closeEditAssignmentForm() {
		showEditAssignmentForm = false;
		selectedAssignment = null;
		editAssignmentData = {
			assignment_id: null,
			user_id: '',
			content: '',
			status: 'pending'
		};
	}

	// 과제 수정 요청
	async function updateAssignment() {
		if (!editAssignmentData.content) {
			dispatch('toast', { message: '과제 내용을 입력해주세요.', type: 'error' });
			return;
		}

		isEditingAssignment = true;
		try {
			const updateData = {
				assignment_id: editAssignmentData.assignment_id,
				user_id: editAssignmentData.user_id,
				content: editAssignmentData.content,
				status: editAssignmentData.status
			};

			const response = await fetchWithAuth(API_ENDPOINTS.ASSIGNMENTS, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(updateData)
			});

			if (response.ok) {
				dispatch('toast', { message: '과제가 성공적으로 수정되었습니다.', type: 'success' });
				closeEditAssignmentForm();
				// 과제 목록 다시 로드
				await loadAssignments();
			} else {
				const errorData = await response.json();
				dispatch('toast', { 
					message: `과제 수정 실패: ${errorData.error || '알 수 없는 오류가 발생했습니다.'}`, 
					type: 'error' 
				});
			}
		} catch (error) {
			console.error('Assignment update error:', error);
			dispatch('toast', { message: '과제 수정 중 오류가 발생했습니다.', type: 'error' });
		} finally {
			isEditingAssignment = false;
		}
	}

	// isVisible이 true가 될 때 자동으로 과제 로드
	$: if (isVisible) {
		loadAssignments();
	}

	// 외부에서 호출할 수 있도록 함수 노출
	export { togglePanel, closePanel };
</script>

{#if isVisible}
	<div class="assignment-panel">
		<div class="assignment-header">
			<h3>📋 과제 목록</h3>
			<button class="assignment-close-btn" on:click={closePanel} aria-label="과제 패널 닫기">
				<svg width="20" height="20" viewBox="0 0 24 24">
					<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" fill="none" stroke-width="2"/>
				</svg>
			</button>
		</div>
		
		<div class="assignment-content">
			{#if showCreateAssignmentForm}
				<!-- 과제 생성 폼 -->
				<div class="assignment-create-form">
					<div class="assignment-form-header">
						<h4>새 과제 생성</h4>
						<button class="assignment-form-close" on:click={closeCreateAssignmentForm} aria-label="폼 닫기">
							<svg width="16" height="16" viewBox="0 0 24 24">
								<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" fill="none" stroke-width="2"/>
							</svg>
						</button>
					</div>
					
					<form on:submit|preventDefault={createAssignment}>
						<div class="assignment-form-group">
							<label for="assignment-user">학생명</label>
							<input
								id="assignment-user"
								type="text"
								bind:value={createAssignmentData.user_id}
								placeholder="학생 이름을 입력하세요"
								required
								class="assignment-form-input"
							/>
						</div>

						<div class="assignment-form-group">
							<label for="assignment-content">과제 내용</label>
							<textarea
								id="assignment-content"
								bind:value={createAssignmentData.content}
								placeholder="과제 내용을 입력하세요"
								required
								class="assignment-form-textarea"
								rows="3"
							/>
						</div>

						<div class="assignment-form-actions">
							<button
								type="button"
								class="btn ghost assignment-form-cancel"
								on:click={closeCreateAssignmentForm}
								disabled={isCreatingAssignment}
							>
								취소
							</button>
							<button 
								type="submit" 
								class="btn primary assignment-form-submit" 
								disabled={isCreatingAssignment}
							>
								{isCreatingAssignment ? '생성 중...' : '과제 생성'}
							</button>
						</div>
					</form>
				</div>
			{:else if showEditAssignmentForm}
				<!-- 과제 수정 폼 -->
				<div class="assignment-edit-form">
					<div class="assignment-form-header">
						<h4>과제 수정</h4>
						<button class="assignment-form-close" on:click={closeEditAssignmentForm} aria-label="폼 닫기">
							<svg width="16" height="16" viewBox="0 0 24 24">
								<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" fill="none" stroke-width="2"/>
							</svg>
						</button>
					</div>
					
					<form on:submit|preventDefault={updateAssignment}>
						<div class="assignment-form-group">
							<label for="edit-assignment-content">과제 내용</label>
							<textarea
								id="edit-assignment-content"
								bind:value={editAssignmentData.content}
								placeholder="과제 내용을 입력하세요"
								required
								class="assignment-form-textarea"
								rows="3"
							/>
						</div>

						<div class="assignment-form-actions">
							<button
								type="button"
								class="btn ghost assignment-form-cancel"
								on:click={closeEditAssignmentForm}
								disabled={isEditingAssignment}
							>
								취소
							</button>
							<button 
								type="submit" 
								class="btn primary assignment-form-submit" 
								disabled={isEditingAssignment}
							>
								{isEditingAssignment ? '수정 중...' : '과제 수정'}
							</button>
						</div>
					</form>
				</div>
			{:else}
				{#if assignmentsLoading}
					<div class="assignment-loading">
						<div class="skeleton head" />
						<div class="skeleton row" />
						<div class="skeleton row" />
					</div>
				{:else if assignments.length === 0}
					<div class="assignment-empty">
						<svg width="48" height="48" viewBox="0 0 24 24" aria-hidden="true">
							<path d="M9 12l2 2 4-4" stroke="currentColor" fill="none" stroke-width="2"/>
						</svg>
						<p>완료되지 않은 과제가 없습니다.</p>
					</div>
				{:else}
					<!-- 반별 토글 과제 목록 -->
					<div class="assignment-classes">
						{#each assignments.filter(classBlock => 
							classBlock.students.some(student => student.assignments.length > 0)
						) as classBlock}
							<div class="assignment-class-section">
								<div class="assignment-class-header" on:click={() => toggleClass(classBlock.class_id)}>
									<h4 class="assignment-class-title">{classBlock.class_name}</h4>
									<div class="assignment-class-toggle">
										<svg 
											class="assignment-toggle-icon" 
											class:expanded={isClassExpanded(classBlock.class_id)}
											width="16" 
											height="16" 
											viewBox="0 0 24 24"
										>
											<path d="M6 9l6 6 6-6" stroke="currentColor" fill="none" stroke-width="2"/>
										</svg>
									</div>
								</div>
								
								{#if isClassExpanded(classBlock.class_id)}
									<div class="assignment-class-content">
										{#each classBlock.students as student}
											{#if student.assignments.length > 0}
												<div class="assignment-student-section">
													<div class="assignment-student-header">
														<span class="assignment-student-name">{student.user_id}</span>
													</div>
													
													<div class="assignment-student-list">
														{#each student.assignments as assignment (assignment.assignment_id)}
															<div class="assignment-item">
																<label class="assignment-checkbox-label">
																	<input 
																		type="checkbox" 
																		class="assignment-checkbox"
																		checked={assignment.status === 'completed'}
																		on:change={(e) => updateAssignmentStatus(assignment, e.target.checked)}
																	/>
																	<span class="assignment-checkmark"></span>
																	<div class="assignment-content-text" on:click={() => openEditAssignmentForm(assignment)}>
																		<div class="assignment-main">
																			<span class="assignment-text">{assignment.content}</span>
																		</div>
																		<div class="assignment-date">
																			{new Date(assignment.created_at).toLocaleDateString('ko-KR', {
																				month: 'short',
																				day: 'numeric'
																			})}
																		</div>
																	</div>
																</label>
															</div>
														{/each}
													</div>
												</div>
											{/if}
										{/each}
									</div>
								{/if}
							</div>
						{/each}
					</div>
				{/if}
			{/if}
		</div>
		
		<!-- 과제 생성 버튼 (폼이 열려있지 않을 때만 표시) -->
		{#if !showCreateAssignmentForm && !showEditAssignmentForm && userRole === 'admin'}
			<div class="assignment-footer">
				<button 
					class="btn primary assignment-create-btn" 
					on:click={openCreateAssignmentForm}
					aria-label="새 과제 생성"
				>
					<svg width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
						<path d="M12 5v14M5 12h14" stroke="currentColor" fill="none" stroke-width="2"/>
					</svg>
					새 과제 생성
				</button>
			</div>
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
		display: flex;
		align-items: center;
		gap: 8px;
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

	.btn svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
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

	/* 과제 패널 스타일 */
	.assignment-panel {
		position: fixed;
		top: 0;
		left: 0;
		width: 400px;
		height: 100vh;
		background: var(--card);
		border-right: 1px solid var(--line);
		z-index: 1000;
		display: flex;
		flex-direction: column;
		box-shadow: 4px 0 20px rgba(0, 0, 0, 0.1);
		transform: translateX(-100%);
		animation: slideInLeft 0.3s ease-out forwards;
	}

	@keyframes slideInLeft {
		to {
			transform: translateX(0);
		}
	}

	.assignment-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20px;
		border-bottom: 1px solid var(--line);
		background: var(--bg);
	}

	.assignment-header h3 {
		margin: 0;
		font-size: 18px;
		font-weight: 700;
		color: var(--text);
	}

	.assignment-close-btn {
		background: none;
		border: none;
		color: var(--muted);
		cursor: pointer;
		padding: 8px;
		border-radius: 8px;
		transition: all 0.2s ease;
	}

	.assignment-close-btn:hover {
		background: var(--line);
		color: var(--text);
	}

	.assignment-content {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 20px;
		overflow-y: auto;
	}

	.assignment-loading {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.assignment-empty {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		flex: 1;
		gap: 16px;
		color: var(--muted);
		text-align: center;
	}

	.assignment-empty svg {
		opacity: 0.5;
	}

	.assignment-empty p {
		font-size: 16px;
		margin: 0;
	}

	.assignment-list {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.assignment-item {
		border: 1px solid var(--line);
		border-radius: 12px;
		background: color-mix(in srgb, var(--card) 92%, var(--bg) 8%);
		transition: all 0.2s ease;
	}

	.assignment-item:hover {
		border-color: var(--brand);
		transform: translateX(2px);
	}

	.assignment-checkbox-label {
		display: flex;
		align-items: flex-start;
		gap: 12px;
		padding: 16px;
		cursor: pointer;
		width: 100%;
	}

	.assignment-checkbox {
		display: none;
	}

	.assignment-checkmark {
		width: 20px;
		height: 20px;
		border: 2px solid var(--line);
		border-radius: 4px;
		background: var(--card);
		position: relative;
		flex-shrink: 0;
		transition: all 0.2s ease;
	}

	.assignment-checkbox:checked + .assignment-checkmark {
		background: var(--brand);
		border-color: var(--brand);
	}

	.assignment-checkbox:checked + .assignment-checkmark::after {
		content: '';
		position: absolute;
		left: 6px;
		top: 2px;
		width: 6px;
		height: 10px;
		border: solid white;
		border-width: 0 2px 2px 0;
		transform: rotate(45deg);
	}

	.assignment-content-text {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 8px;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.assignment-content-text:hover {
		background: color-mix(in srgb, var(--brand) 5%, transparent);
		border-radius: 6px;
		padding: 4px;
		margin: -4px;
	}

	.assignment-main {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.assignment-text {
		font-size: 14px;
		font-weight: 500;
		color: var(--text);
		line-height: 1.4;
	}

	.assignment-meta {
		display: flex;
		gap: 8px;
		align-items: center;
	}

	.assignment-student {
		font-size: 12px;
		font-weight: 600;
		color: var(--brand);
		background: color-mix(in srgb, var(--brand) 10%, transparent);
		padding: 2px 6px;
		border-radius: 4px;
	}

	.assignment-class {
		font-size: 12px;
		color: var(--muted);
		background: var(--bg);
		padding: 2px 6px;
		border-radius: 4px;
	}

	.assignment-date {
		font-size: 12px;
		color: var(--muted);
		align-self: flex-end;
	}

	/* 과제 생성/수정 폼 스타일 */
	.assignment-create-form,
	.assignment-edit-form {
		background: color-mix(in srgb, var(--bg) 30%, transparent);
		border: 1px solid var(--line);
		border-radius: 12px;
		padding: 16px;
		margin-bottom: 16px;
	}

	.assignment-form-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 16px;
	}

	.assignment-form-header h4 {
		margin: 0;
		font-size: 16px;
		font-weight: 700;
		color: var(--text);
	}

	.assignment-form-close {
		background: none;
		border: none;
		color: var(--muted);
		cursor: pointer;
		padding: 4px;
		border-radius: 4px;
		transition: all 0.2s ease;
	}

	.assignment-form-close:hover {
		background: var(--line);
		color: var(--text);
	}

	.assignment-form-group {
		display: flex;
		flex-direction: column;
		gap: 6px;
		margin-bottom: 16px;
	}

	.assignment-form-group label {
		font-weight: 600;
		font-size: 13px;
		color: var(--text);
	}

	.assignment-form-input,
	.assignment-form-textarea {
		border: 1px solid var(--line);
		background: var(--card);
		color: var(--text);
		padding: 10px;
		border-radius: 8px;
		font-size: 14px;
		transition: all 0.2s ease;
		font-family: inherit;
	}

	.assignment-form-input:focus,
	.assignment-form-textarea:focus {
		outline: none;
		border-color: var(--brand);
		box-shadow: 0 0 0 2px color-mix(in srgb, var(--brand) 20%, transparent);
	}

	.assignment-form-textarea {
		resize: vertical;
		min-height: 60px;
	}

	.assignment-form-actions {
		display: flex;
		gap: 8px;
		justify-content: flex-end;
	}

	.assignment-form-cancel,
	.assignment-form-submit {
		padding: 8px 16px;
		font-size: 13px;
		font-weight: 600;
		border-radius: 6px;
	}

	/* 과제 생성 버튼 스타일 */
	.assignment-footer {
		padding: 16px 20px;
		border-top: 1px solid var(--line);
		background: var(--bg);
	}

	.assignment-create-btn {
		width: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		padding: 12px 16px;
		font-size: 14px;
		font-weight: 600;
		border-radius: 8px;
		transition: all 0.2s ease;
	}

	.assignment-create-btn:hover {
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
	}

	.assignment-create-btn svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	/* 반별 토글 스타일 */
	.assignment-classes {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.assignment-class-section {
		border: 1px solid var(--line);
		border-radius: 12px;
		background: color-mix(in srgb, var(--card) 95%, var(--bg) 5%);
		overflow: hidden;
	}

	.assignment-class-header {
		background: var(--bg);
		padding: 12px 16px;
		border-bottom: 1px solid var(--line);
		display: flex;
		justify-content: space-between;
		align-items: center;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.assignment-class-header:hover {
		background: color-mix(in srgb, var(--bg) 80%, var(--brand) 20%);
	}

	.assignment-class-title {
		margin: 0;
		font-size: 16px;
		font-weight: 700;
		color: var(--text);
	}

	.assignment-class-toggle {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 24px;
		height: 24px;
		border-radius: 4px;
		transition: all 0.2s ease;
	}

	.assignment-class-toggle:hover {
		background: color-mix(in srgb, var(--brand) 20%, transparent);
	}

	.assignment-toggle-icon {
		transition: transform 0.2s ease;
		color: var(--muted);
	}

	.assignment-toggle-icon.expanded {
		transform: rotate(180deg);
		color: var(--brand);
	}

	.assignment-class-content {
		animation: slideDown 0.2s ease-out;
	}

	@keyframes slideDown {
		from {
			opacity: 0;
			max-height: 0;
		}
		to {
			opacity: 1;
			max-height: 1000px;
		}
	}

	.assignment-student-section {
		border-bottom: 1px solid var(--line);
	}

	.assignment-student-section:last-child {
		border-bottom: none;
	}

	.assignment-student-header {
		background: color-mix(in srgb, var(--bg) 50%, transparent);
		padding: 8px 16px;
		border-bottom: 1px solid var(--line);
	}

	.assignment-student-name {
		font-size: 14px;
		font-weight: 600;
		color: var(--brand);
	}

	.assignment-student-list {
		padding: 8px;
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	/* 모바일 대응 */
	@media (max-width: 768px) {
		.assignment-panel {
			width: 100%;
			left: 0;
			right: 0;
		}

		.assignment-form-actions {
			flex-direction: column;
		}

		.assignment-form-cancel,
		.assignment-form-submit {
			width: 100%;
		}
	}
</style>
