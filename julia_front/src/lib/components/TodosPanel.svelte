<script>
	import { createEventDispatcher } from 'svelte';
	import { fetchWithAuth } from '$lib/auth';
	import { API_ENDPOINTS } from '$lib/config';

	const dispatch = createEventDispatcher();

	export let isVisible = false;
	export let userRole = 'student';
	export let currentUser = null;

	let todos = [];
	let isLoading = false;
	let newTodoTitle = '';
	let newTodoDescription = '';
	let showAddForm = false;
	
	// 관리자용 상태
	let classes = [];
	let expandedClasses = [-1]; // 기본적으로 "할 일" 섹션을 열어둠
	let classTodos = {}; // 반별 할 일 저장
	let newTodoClassId = null; // 새 할 일 추가 시 선택할 반 ID (null = 선택 안함)

	// Todos 패널 열기/닫기
	function togglePanel() {
		isVisible = !isVisible;
		if (isVisible) {
			loadTodos();
		}
		dispatch('toggle', { isVisible });
	}

	// 반 토글 함수 (과제와 동일한 방식)
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

	// currentUser가 변경될 때도 할 일 로드
	$: if (isVisible && currentUser?.id) {
		if (userRole === 'admin') {
			loadClasses();
		}
		loadTodos();
	}

	function closePanel() {
		isVisible = false;
		dispatch('close');
	}

	// 반 목록 로드
	async function loadClasses() {
		try {
			const response = await fetchWithAuth(API_ENDPOINTS.CLASSES);
			if (response.ok) {
				classes = await response.json();
			} else {
				console.error('반 목록 로드 실패');
			}
		} catch (error) {
			console.error('반 목록 로드 오류:', error);
		}
	}

	// Todos 로드
	async function loadTodos() {
		if (!currentUser?.id) {
			todos = [];
			return;
		}

		isLoading = true;
		
		try {
			const response = await fetchWithAuth(API_ENDPOINTS.TODOS_BY_USER(currentUser.id));
			if (response.ok) {
				const classTodosData = await response.json();
				
				// 반별로 분류 및 null 값 필터링
				classTodos = {};
				const allTodos = [];
				
				for (const classData of classTodosData) {
					const validTodos = (classData.todos || []).filter(todo => todo !== null);
					classTodos[classData.class_id] = validTodos;
					allTodos.push(...validTodos);
				}
				
				// 관리자는 모든 섹션을 기본적으로 열어둠
				if (userRole === 'admin') {
					expandedClasses = Object.keys(classTodos).map(id => parseInt(id));
				}
				
				todos = allTodos;
			} else {
				todos = [];
			}
		} catch (error) {
			console.error('할 일 로드 오류:', error);
			todos = [];
		} finally {
			isLoading = false;
		}
	}

	// 새 Todo 추가
	async function addTodo() {
		if (!newTodoTitle.trim()) return;

		try {
			const todoData = {
				title: newTodoTitle.trim(),
				description: newTodoDescription.trim(),
				user_id: currentUser?.id || '',
				completed: false
			};

			// 관리자는 class_id를 payload에 포함
			if (userRole === 'admin' && newTodoClassId !== null) {
				todoData.class_id = newTodoClassId;
			}

			const response = await fetchWithAuth(API_ENDPOINTS.TODOS, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(todoData)
			});

			if (response.ok) {
				const newTodo = await response.json();
				
				// todos 배열에 추가
				todos = [...todos, newTodo];
				
				// classTodos도 업데이트
				const classId = newTodo.class_id || -1;
				if (!classTodos[classId]) {
					classTodos[classId] = [];
				}
				classTodos[classId] = [...classTodos[classId], newTodo];
				classTodos = { ...classTodos };
				
				newTodoTitle = '';
				newTodoDescription = '';
				newTodoClassId = null; // 관리자용 반 선택 초기화
				showAddForm = false;
				dispatch('toast', { message: '할 일이 추가되었습니다.', type: 'success' });
			} else {
				dispatch('toast', { message: '할 일 추가에 실패했습니다.', type: 'error' });
			}
		} catch (error) {
			console.error('Todo 추가 오류:', error);
			dispatch('toast', { message: '할 일 추가에 실패했습니다.', type: 'error' });
		}
	}

	// Todo 완료 상태 토글
	async function toggleTodo(todo) {
		try {
			const updatedTodo = {
				...todo,
				completed: !todo.completed
			};

			const response = await fetchWithAuth(`${API_ENDPOINTS.TODOS}/${todo.id}`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(updatedTodo)
			});

			if (response.ok) {
				// todos 배열 업데이트
				todos = todos.map(t => t.id === todo.id ? { ...t, completed: !t.completed } : t);
				
				// classTodos도 업데이트 (관리자용)
				for (const classId in classTodos) {
					classTodos[classId] = classTodos[classId].map(t => 
						t.id === todo.id ? { ...t, completed: !t.completed } : t
					);
				}
				classTodos = { ...classTodos };
			} else {
				dispatch('toast', { message: '상태 변경에 실패했습니다.', type: 'error' });
			}
		} catch (error) {
			console.error('Todo 상태 변경 오류:', error);
			dispatch('toast', { message: '상태 변경에 실패했습니다.', type: 'error' });
		}
	}

	// Todo 삭제
	async function deleteTodo(todo) {
		if (!confirm('정말로 이 할 일을 삭제하시겠습니까?')) return;

		try {
			const response = await fetchWithAuth(`${API_ENDPOINTS.TODOS}/${todo.id}`, {
				method: 'DELETE'
			});

			if (response.ok) {
				// todos 배열에서 제거
				todos = todos.filter(t => t.id !== todo.id);
				
				// classTodos에서도 제거
				for (const classId in classTodos) {
					classTodos[classId] = classTodos[classId].filter(t => t.id !== todo.id);
				}
				classTodos = { ...classTodos };
				
				dispatch('toast', { message: '할 일이 삭제되었습니다.', type: 'success' });
			} else {
				dispatch('toast', { message: '할 일 삭제에 실패했습니다.', type: 'error' });
			}
		} catch (error) {
			console.error('Todo 삭제 오류:', error);
			dispatch('toast', { message: '할 일 삭제에 실패했습니다.', type: 'error' });
		}
	}

	// 완료된 Todo 개수 계산
	$: completedCount = todos.filter(t => t.completed).length;
	$: totalCount = todos.length;
	$: progressPercentage = totalCount > 0 ? (completedCount / totalCount) * 100 : 0;

	// 외부에서 호출할 수 있도록 함수 노출
	export { togglePanel, closePanel };
</script>

{#if isVisible}
	<div class="todos-panel">
		<div class="todos-header">
			<h3>📝 할 일 목록</h3>
			<button class="todos-close-btn" on:click={closePanel} aria-label="할 일 목록 닫기">
				<svg width="20" height="20" viewBox="0 0 24 24">
					<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" fill="none" stroke-width="2"/>
				</svg>
			</button>
		</div>
		
		<div class="todos-content">
			<!-- 통계 -->
			<div class="todos-stats">
				<span class="stats-text">완료: {completedCount}/{totalCount}</span>
				{#if totalCount > 0}
					<div class="progress-bar">
						<div class="progress-fill" style="width: {progressPercentage}%"></div>
					</div>
				{/if}
			</div>
				<!-- 새 할 일 추가 폼 -->
				<div class="add-todo-section">
				{#if showAddForm}
					<div class="add-todo-form">
						{#if userRole === 'admin'}
							<div class="form-group">
								<label for="new-todo-class">반 선택 (선택사항)</label>
								<select
									id="new-todo-class"
									bind:value={newTodoClassId}
									class="todo-select"
								>
									<option value={null}>선택 안함</option>
									{#each classes as classItem}
										<option value={classItem.class_id}>{classItem.class_name}</option>
									{/each}
								</select>
							</div>
						{/if}
						<input
							bind:value={newTodoTitle}
							placeholder="할 일 제목을 입력하세요"
							class="todo-input"
							on:keydown={(e) => e.key === 'Enter' && addTodo()}
						/>
						<textarea
							bind:value={newTodoDescription}
							placeholder="설명 (선택사항)"
							class="todo-textarea"
							rows="2"
						></textarea>
						<div class="form-actions">
							<button class="btn primary" on:click={addTodo} disabled={!newTodoTitle.trim()}>
								추가
							</button>
							<button class="btn ghost" on:click={() => { showAddForm = false; newTodoTitle = ''; newTodoDescription = ''; newTodoClassId = null; }}>
								취소
							</button>
						</div>
					</div>
				{:else}
					<button class="btn primary add-todo-btn" on:click={() => showAddForm = true}>
						<svg width="16" height="16" viewBox="0 0 24 24">
							<path d="M12 5v14M5 12h14" stroke="currentColor" fill="none" stroke-width="2"/>
						</svg>
						새 할 일 추가
					</button>
				{/if}
			</div>

			<!-- 할 일 목록 -->
			<div class="todos-list">
				{#if isLoading}
					<div class="loading">할 일을 불러오는 중...</div>
				{:else if userRole === 'admin'}
					<!-- 관리자용: 반별 토글 방식 -->
					<!-- 일반 할 일 섹션 -->
					{#if classTodos[-1] && classTodos[-1].length > 0}
						<div class="todo-class-section">
							<div class="todo-class-header" on:click={() => toggleClass(-1)}>
								<h4 class="todo-class-title">할 일</h4>
								<div class="todo-class-toggle">
									<svg 
										class="todo-toggle-icon" 
										class:expanded={isClassExpanded(-1)}
										width="16" 
										height="16" 
										viewBox="0 0 24 24"
									>
										<path d="M6 9l6 6 6-6" stroke="currentColor" fill="none" stroke-width="2"/>
									</svg>
								</div>
							</div>
							
							{#if isClassExpanded(-1)}
								<div class="todo-class-content">
									{#each classTodos[-1] as todo}
										<div class="todo-item" class:completed={todo.completed}>
											<div class="todo-main">
												<button 
													class="todo-checkbox" 
													on:click={() => toggleTodo(todo)}
													aria-label={todo.completed ? '완료 취소' : '완료 표시'}
												>
													{#if todo.completed}
														<svg width="16" height="16" viewBox="0 0 24 24">
															<path d="M9 12l2 2 4-4" stroke="currentColor" fill="none" stroke-width="2"/>
														</svg>
													{/if}
												</button>
												<div class="todo-content">
													<h4 class="todo-title" class:completed={todo.completed}>{todo.title}</h4>
													{#if todo.description}
														<p class="todo-description">{todo.description}</p>
													{/if}
													<div class="todo-meta">
														<span class="todo-date">{new Date(todo.created_at).toLocaleString('ko-KR', {
															year: 'numeric',
															month: '2-digit',
															day: '2-digit',
															hour: '2-digit',
															minute: '2-digit'
														})}</span>
													</div>
												</div>
												<button 
													class="todo-delete-btn" 
													on:click={() => deleteTodo(todo)}
													aria-label="삭제"
												>
													<svg width="16" height="16" viewBox="0 0 24 24">
														<path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0v14M10 11v6M14 11v6" stroke="currentColor" fill="none" stroke-width="2"/>
													</svg>
												</button>
											</div>
										</div>
									{/each}
								</div>
							{/if}
						</div>
					{/if}

					<!-- 반별 할 일 섹션 -->
					{#each classes as classItem}
						{#if classTodos[classItem.class_id] && classTodos[classItem.class_id].length > 0}
							<div class="todo-class-section">
								<div class="todo-class-header" on:click={() => toggleClass(classItem.class_id)}>
									<h4 class="todo-class-title">{classItem.class_name}</h4>
									<div class="todo-class-toggle">
										<svg 
											class="todo-toggle-icon" 
											class:expanded={isClassExpanded(classItem.class_id)}
											width="16" 
											height="16" 
											viewBox="0 0 24 24"
										>
											<path d="M6 9l6 6 6-6" stroke="currentColor" fill="none" stroke-width="2"/>
										</svg>
									</div>
								</div>
								
								{#if isClassExpanded(classItem.class_id)}
									<div class="todo-class-content">
										{#each classTodos[classItem.class_id] as todo}
											<div class="todo-item" class:completed={todo.completed}>
												<div class="todo-main">
													<button 
														class="todo-checkbox" 
														on:click={() => toggleTodo(todo)}
														aria-label={todo.completed ? '완료 취소' : '완료 표시'}
													>
														{#if todo.completed}
															<svg width="16" height="16" viewBox="0 0 24 24">
																<path d="M9 12l2 2 4-4" stroke="currentColor" fill="none" stroke-width="2"/>
															</svg>
														{/if}
													</button>
													<div class="todo-content">
														<h4 class="todo-title" class:completed={todo.completed}>{todo.title}</h4>
														{#if todo.description}
															<p class="todo-description">{todo.description}</p>
														{/if}
														<div class="todo-meta">
															<span class="todo-date">{new Date(todo.created_at).toLocaleString('ko-KR', {
																year: 'numeric',
																month: '2-digit',
																day: '2-digit',
																hour: '2-digit',
																minute: '2-digit'
															})}</span>
														</div>
													</div>
													<button 
														class="todo-delete-btn" 
														on:click={() => deleteTodo(todo)}
														aria-label="삭제"
													>
														<svg width="16" height="16" viewBox="0 0 24 24">
															<path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0v14M10 11v6M14 11v6" stroke="currentColor" fill="none" stroke-width="2"/>
														</svg>
													</button>
												</div>
											</div>
										{/each}
									</div>
								{/if}
							</div>
						{/if}
					{/each}

					<!-- 할 일이 없는 경우 -->
					{#if todos.length === 0}
						<div class="empty-state">
							<svg width="48" height="48" viewBox="0 0 24 24">
								<path d="M9 12l2 2 4-4" stroke="currentColor" fill="none" stroke-width="2"/>
								<path d="M21 12c0 1.66-1.34 3-3 3H6c-1.66 0-3-1.34-3-3s1.34-3 3-3h12c1.66 0 3 1.34 3 3z" stroke="currentColor" fill="none" stroke-width="2"/>
							</svg>
							<p>할 일이 없습니다</p>
							<p class="empty-hint">새 할 일을 추가해보세요!</p>
						</div>
					{/if}
				{:else}
					<!-- 학생용: 일반 목록 -->
					{#if todos.length === 0}
						<div class="empty-state">
							<svg width="48" height="48" viewBox="0 0 24 24">
								<path d="M9 12l2 2 4-4" stroke="currentColor" fill="none" stroke-width="2"/>
								<path d="M21 12c0 1.66-1.34 3-3 3H6c-1.66 0-3-1.34-3-3s1.34-3 3-3h12c1.66 0 3 1.34 3 3z" stroke="currentColor" fill="none" stroke-width="2"/>
							</svg>
							<p>할 일이 없습니다</p>
							<p class="empty-hint">새 할 일을 추가해보세요!</p>
						</div>
					{:else}
						{#each todos as todo}
							<div class="todo-item" class:completed={todo.completed}>
								<div class="todo-main">
									<button 
										class="todo-checkbox" 
										on:click={() => toggleTodo(todo)}
										aria-label={todo.completed ? '완료 취소' : '완료 표시'}
									>
										{#if todo.completed}
											<svg width="16" height="16" viewBox="0 0 24 24">
												<path d="M9 12l2 2 4-4" stroke="currentColor" fill="none" stroke-width="2"/>
											</svg>
										{/if}
									</button>
									<div class="todo-content">
										<h4 class="todo-title" class:completed={todo.completed}>{todo.title}</h4>
										{#if todo.description}
											<p class="todo-description">{todo.description}</p>
										{/if}
										<div class="todo-meta">
											<span class="todo-date">{new Date(todo.created_at).toLocaleString('ko-KR', {
												year: 'numeric',
												month: '2-digit',
												day: '2-digit',
												hour: '2-digit',
												minute: '2-digit'
											})}</span>
										</div>
									</div>
									<button 
										class="todo-delete-btn" 
										on:click={() => deleteTodo(todo)}
										aria-label="삭제"
									>
										<svg width="16" height="16" viewBox="0 0 24 24">
											<path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0v14M10 11v6M14 11v6" stroke="currentColor" fill="none" stroke-width="2"/>
										</svg>
									</button>
								</div>
							</div>
						{/each}
					{/if}
				{/if}
			</div>
		</div>
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
		--sun: #ef4444;
		--success: #10b981;
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

	/* Todos 패널 스타일 */
	.todos-panel {
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

	.todos-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20px;
		border-bottom: 1px solid var(--line);
		background: var(--bg);
	}

	.todos-header h3 {
		margin: 0;
		font-size: 18px;
		font-weight: 700;
		color: var(--text);
	}

	.todos-close-btn {
		background: none;
		border: none;
		color: var(--muted);
		cursor: pointer;
		padding: 8px;
		border-radius: 8px;
		transition: all 0.2s ease;
	}

	.todos-close-btn:hover {
		background: var(--line);
		color: var(--text);
	}

	.todos-content {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 20px;
		overflow-y: auto;
	}

	/* 반별 토글 섹션 스타일 (과제와 동일) */
	.todo-class-section {
		margin-bottom: 16px;
		border: 1px solid var(--line);
		border-radius: 12px;
		overflow: hidden;
		background: var(--card);
	}

	.todo-class-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 16px 20px;
		background: var(--bg);
		cursor: pointer;
		transition: background-color 0.2s ease;
		border-bottom: 1px solid var(--line);
	}

	.todo-class-header:hover {
		background: color-mix(in srgb, var(--brand) 5%, var(--bg));
	}

	.todo-class-title {
		font-size: 16px;
		font-weight: 600;
		color: var(--text);
		margin: 0;
	}

	.todo-class-toggle {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		border-radius: 6px;
		background: transparent;
		transition: all 0.2s ease;
	}

	.todo-class-toggle:hover {
		background: color-mix(in srgb, var(--brand) 10%, transparent);
	}

	.todo-toggle-icon {
		width: 16px;
		height: 16px;
		color: var(--muted);
		transition: transform 0.2s ease;
	}

	.todo-toggle-icon.expanded {
		transform: rotate(180deg);
		color: var(--brand);
	}

	.todo-class-content {
		padding: 0;
		background: var(--card);
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 8px;
		margin-bottom: 12px;
	}

	.form-group label {
		font-weight: 600;
		font-size: 14px;
		color: var(--text);
	}

	.todo-select {
		border: 1px solid var(--line);
		border-radius: 8px;
		padding: 12px;
		font-family: inherit;
		font-size: 14px;
		background: var(--card);
		color: var(--text);
		transition: all 0.2s ease;
	}

	.todo-select:focus {
		outline: none;
		border-color: var(--brand);
		box-shadow: 0 0 0 2px color-mix(in srgb, var(--brand) 20%, transparent);
	}

	.todos-stats {
		display: flex;
		align-items: center;
		gap: 12px;
		margin-bottom: 20px;
		padding: 12px;
		background: var(--bg);
		border-radius: 10px;
	}

	.stats-text {
		font-size: 14px;
		font-weight: 600;
		color: var(--text);
	}

	.progress-bar {
		flex: 1;
		height: 6px;
		background: var(--line);
		border-radius: 3px;
		overflow: hidden;
	}

	.progress-fill {
		height: 100%;
		background: var(--success);
		transition: width 0.3s ease;
	}

	.add-todo-section {
		margin-bottom: 20px;
	}

	.add-todo-btn {
		width: 100%;
		justify-content: center;
	}

	.add-todo-form {
		display: flex;
		flex-direction: column;
		gap: 12px;
		padding: 16px;
		background: var(--bg);
		border-radius: 10px;
		border: 1px solid var(--line);
	}

	.todo-input,
	.todo-textarea {
		border: 1px solid var(--line);
		border-radius: 8px;
		padding: 12px;
		font-family: inherit;
		font-size: 14px;
		background: var(--card);
		color: var(--text);
		transition: all 0.2s ease;
	}

	.todo-input:focus,
	.todo-textarea:focus {
		outline: none;
		border-color: var(--brand);
		box-shadow: 0 0 0 2px color-mix(in srgb, var(--brand) 20%, transparent);
	}

	.todo-textarea {
		resize: vertical;
		min-height: 60px;
	}

	.form-actions {
		display: flex;
		gap: 8px;
	}

	.todos-list {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.loading {
		text-align: center;
		padding: 40px;
		color: var(--muted);
		font-style: italic;
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 40px 20px;
		text-align: center;
		color: var(--muted);
	}

	.empty-state svg {
		opacity: 0.5;
		margin-bottom: 16px;
	}

	.empty-state p {
		margin: 0;
		font-size: 16px;
	}

	.empty-hint {
		font-size: 14px;
		margin-top: 8px;
		opacity: 0.8;
	}

	.todo-item {
		border: 1px solid var(--line);
		border-radius: 10px;
		padding: 16px;
		background: var(--card);
		transition: all 0.2s ease;
	}

	.todo-item:hover {
		border-color: var(--brand);
		transform: translateX(4px);
	}

	.todo-item.completed {
		opacity: 0.7;
		background: color-mix(in srgb, var(--success) 5%, var(--card));
	}

	.todo-main {
		display: flex;
		align-items: flex-start;
		gap: 12px;
	}

	.todo-checkbox {
		width: 24px;
		height: 24px;
		border: 2px solid var(--line);
		border-radius: 6px;
		background: var(--card);
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.2s ease;
		flex-shrink: 0;
		margin-top: 2px;
	}

	.todo-checkbox:hover {
		border-color: var(--brand);
	}

	.todo-item.completed .todo-checkbox {
		background: var(--success);
		border-color: var(--success);
		color: white;
	}

	.todo-content {
		flex: 1;
		min-width: 0;
	}

	.todo-title {
		margin: 0 0 8px 0;
		font-size: 16px;
		font-weight: 600;
		color: var(--text);
		line-height: 1.4;
	}

	.todo-title.completed {
		text-decoration: line-through;
		color: var(--muted);
	}

	.todo-description {
		margin: 0 0 12px 0;
		font-size: 14px;
		color: var(--muted);
		line-height: 1.5;
	}

	.todo-meta {
		display: flex;
		justify-content: flex-end;
		align-items: center;
		font-size: 12px;
		color: var(--muted);
	}

	.todo-delete-btn {
		background: none;
		border: none;
		color: var(--muted);
		cursor: pointer;
		padding: 4px;
		border-radius: 4px;
		transition: all 0.2s ease;
		flex-shrink: 0;
	}

	.todo-delete-btn:hover {
		background: var(--sun);
		color: white;
	}

	.todo-delete-btn svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	/* 모바일 대응 */
	@media (max-width: 768px) {
		.todos-panel {
			width: 100%;
			left: 0;
			right: 0;
		}
		
		.form-actions {
			flex-direction: column;
		}
		
		.todo-meta {
			justify-content: flex-start;
		}
	}
</style>
